package events

// ReadEvent is an event generated from BPF based on the kernel.
type ReadEvent struct {
	Pid      uint32
	Username string
	Line     string
}

// DocumentEvent is an event after it has been processed and stored in the database,
// which is where the ID and Timestamp come from.
// ID is an auto-incrementing integer.
type DocumentEvent struct {
	ID        int64
	Timestamp int64
	Username  string
	Command   string
}

// Document represents a record in typesense format.
type Document struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Username  string `json:"username"`
	Command   string `json:"command"`
}
