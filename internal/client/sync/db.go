package sync

import (
	"encoding/binary"
	"fmt"
	"sort"

	"report/internal/client/crypto"
	"report/pkg/patch"

	"github.com/dgraph-io/badger/v4"
)

type LocalDB struct {
	db     *badger.DB
	aesKey []byte
}

func OpenLocalDB(path string, aesKey []byte) (*LocalDB, error) {
	opts := badger.DefaultOptions(path).
		WithSyncWrites(true).
		WithEncryptionKey(aesKey).
		WithIndexCacheSize(64 << 20).
		WithNumVersionsToKeep(1).
		WithLogger(nil)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("open badgerdb at %s: %w", path, err)
	}
	return &LocalDB{db: db, aesKey: aesKey}, nil
}

func (d *LocalDB) Close() error {
	return d.db.Close()
}

// ── Outbox ────────────────────────────────────────────────────────────────────

func (d *LocalDB) Set(key string, value []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func (d *LocalDB) Delete(key string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (d *LocalDB) CountOutboxEntries() int {
	count := 0
	d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte("report:")
		suffix := []byte(":outbox:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			if containsBytes(it.Item().Key(), suffix) {
				count++
			}
		}
		return nil
	})
	return count
}

// ListOutboxPatches returns all pending patches sorted by ClientSeq.
// Ordering matters — we must publish in the same order they were created.
func (d *LocalDB) ListOutboxPatches() ([]patch.Patch, error) {
	var patches []patch.Patch

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("report:")
		outboxMark := []byte(":outbox:")

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			if !containsBytes(item.Key(), outboxMark) {
				continue
			}
			return item.Value(func(val []byte) error {
				plain, err := crypto.Decrypt(d.aesKey, val)
				if err != nil {
					return fmt.Errorf("decrypt outbox entry: %w", err)
				}
				var p patch.Patch
				if err := p.DecodeMsgpack(plain); err != nil {
					return fmt.Errorf("decode outbox entry: %w", err)
				}
				patches = append(patches, p)
				return nil
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort by ClientSeq to maintain causal ordering
	sort.Slice(patches, func(i, j int) bool {
		return patches[i].ClientSeq < patches[j].ClientSeq
	})
	return patches, nil
}

// ── Section data cache ────────────────────────────────────────────────────────

func (d *LocalDB) ApplyIncomingPatch(p patch.Patch) error {
	key := fmt.Sprintf("report:%s:section:%s", p.ReportID, p.SectionKey)
	// Read current section data, apply JSON patch op, write back
	// Uses RFC 6902 apply logic
	return d.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		var current []byte
		if err == badger.ErrKeyNotFound {
			current = []byte("{}")
		} else if err != nil {
			return err
		} else {
			item.Value(func(v []byte) error {
				plain, _ := crypto.Decrypt(d.aesKey, v)
				current = plain
				return nil
			})
		}
		updated, err := applyJSONPatchOp(current, p)
		if err != nil {
			return err
		}
		encrypted, err := crypto.Encrypt(d.aesKey, updated)
		if err != nil {
			return err
		}
		return txn.Set([]byte(key), encrypted)
	})
}

// ── Sequence tracking ─────────────────────────────────────────────────────────

func (d *LocalDB) GetLastSeq() uint64 {
	var seq uint64
	d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("app:last_nats_seq"))
		if err != nil {
			return nil
		}
		return item.Value(func(val []byte) error {
			if len(val) == 8 {
				seq = binary.BigEndian.Uint64(val)
			}
			return nil
		})
	})
	return seq
}

func (d *LocalDB) SaveLastSeq(seq uint64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, seq)
	d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("app:last_nats_seq"), buf)
	})
}

// ── Generic KV (for JWT, license, settings) ───────────────────────────────────

func (d *LocalDB) Get(key string) string {
	var result string
	d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return nil
		}
		return item.Value(func(val []byte) error {
			plain, err := crypto.Decrypt(d.aesKey, val)
			if err != nil {
				return err
			}
			result = string(plain)
			return nil
		})
	})
	return result
}

func (d *LocalDB) GetBytes(key string) []byte {
	var result []byte
	d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return nil
		}
		return item.Value(func(val []byte) error {
			plain, err := crypto.Decrypt(d.aesKey, val)
			if err != nil {
				return err
			}
			result = plain
			return nil
		})
	})
	return result
}
