package patch

import (
	"time"

	"github.com/vmihailenco/msgpack/v5"
	//"github.com/vmihailenco/msgpack/v5"
)

// Patch is the atomic unit of synchronisation
// NEVER sync whole documents. Always sync Patches.
// ClientSeq is a Lamport Logical clock - NEVER use wall time for ordering.
type Patch struct {
	ID       string `json:"id"           msgpack:"id"`
	ReportID string `json:"reportID"     msgpack:"reportID"`

	// OrgID is the organisation slug
	// Used to route NATS messages: reports.{OrgID}.{ReportID}.patch
	OrgID      string `json:"orgID"        msgpack:"orgID"`
	SectionKey string `json:"sectionKey"   msgpack:"sectionKey"`
	FieldPath  string `json:"FieldPath"    msgpack:"FieldPath"`

	// Op is a JSON Patch operation: "add", "replace", "remove"
	Op         string `json:"op"           msgpack:"op"`
	Value      any    `json:"value,omitempty"   msgpack:"value,omitempty"`
	AuthorID   string `json:"authorID"     msgpack:"authorID"`
	AuthorRole string `json:"authorRole"   msgpack:"authorRole"`

	// ClientSeq is a Lamport logical clock value
	// Increasing counter that establishes causal ordering regardless of client clock
	ClientSeq uint64 `json:"clientSeq"    msgpack:"clientSeq"`

	// ClientTS is the wall-clock time when the edit was made
	// Used only for display in the UI like: "edited 5 minutes ago"
	ClientTS time.Time `json:"clientTime"  msgpack:"clientTime"`
}

// Note: we are using MessagePack(vimihailenco/msgpack) for NATS payload 30%-40% smaller than JSON

// EncodeMsgpack serialises the patch to MessagePack binary
// This is what gets stored in BadgerDB and published to NATS
func (p *Patch) EncodeMsgpack() ([]byte, error) {
	return msgpack.Marshal(p)
}

// DecodeMsgpack deserialises a MessagePack binary back into a Patch
func (p *Patch) DecodeMsgpack(data []byte) error {
	return msgpack.Unmarshal(data, p)
}
