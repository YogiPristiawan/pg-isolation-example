package read_phenomena

type ReadPhenomena uint8

const (
	READ_PHENOMENA_UNSPECIFIED ReadPhenomena = iota

	READ_PHENOMENA_NON_REPEATABLE_READ
	READ_PHENOMENA_PHANTOM_READ
	READ_PHENOMENA_SERIALIZATION_ANOMALY
)

func (r ReadPhenomena) Valid() bool {
	switch r {
	case READ_PHENOMENA_NON_REPEATABLE_READ, READ_PHENOMENA_PHANTOM_READ, READ_PHENOMENA_SERIALIZATION_ANOMALY:
		return true
	default:
		return false
	}
}

type Product struct {
	ID       int64
	Name     string
	Quantity int64
}
