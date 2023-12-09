package main

type Topic uint8

const (
	TOPIC_UNSPECIFIED Topic = iota
	TOPIC_READ_PHENOMENA
	TOPIC_PG_ISOLATION_LEVEL
)

func (t Topic) Valid() bool {
	switch t {
	case TOPIC_READ_PHENOMENA, TOPIC_PG_ISOLATION_LEVEL:
		return true
	default:
		return false
	}
}
