package report

import (
	"context"
	"crypto/sha256"
	"encoding/binary"

	"github.com/jackc/pgx"
)

func sectionLocKey(reportID, sectionKey string) int64 {
	h := sha256.Sum256([]byte(reportID + "_" + sectionKey))
	return int64(binary.BigEndian.Uint64(h[:8]))
}

func AcquireSectionLock(ctx context.Context, conn *pgx.Conn, reportID, sectionKey string) (bool, error) {
	var acquired bool
	err := conn.QueryRow(ctx,
		"SELECT pg_try_advisory_xact_lock($1)",
		sectionLocKey(reportID, sectionKey)).Scan(&acquired)
	return acquired, err
}
