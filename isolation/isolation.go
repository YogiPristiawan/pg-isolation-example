package isolation

type Isolation uint8

const (
	ISOLATION_UNSPECIFIED Isolation = iota

	ISOLATION_READ_COMMITTED
	ISOLATION_REPEATABLE_READ
	ISOLATION_SERIALIZABLE
)

func (i Isolation) Valid() bool {
	switch i {
	case ISOLATION_READ_COMMITTED, ISOLATION_SERIALIZABLE, ISOLATION_REPEATABLE_READ:
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
