package repository

type logEntry struct {
	KeyLength   int8
	Key         string
	ValueLength int64
	Type        logEntryType
	Value       string
}

// The entire length of the entry's bytes.
func (le logEntry) Length() int {
	return 2 + int(le.KeyLength) + int(le.ValueLength)
}

type logEntryType int8

const (
	LET_Deleted logEntryType = iota
	LET_Live
)

func (t logEntryType) String() string {
	switch t {
	case LET_Deleted:
		return "Deleted"
	case LET_Live:
		return "Live"
	default:
		return "Unknown"
	}
}