package music

type Interval interface {
	JustInterval | TemperedInterval
}

type Scale interface {
	String() string
	System() string
	Description() string
	Intervals() []any
}
