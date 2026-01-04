package music

type JustScale struct {
	system      string
	description string
	algorithm   computeIntervalsFn
}

type Symmetry string

const (
	Asymmetric Symmetry = "asymmetric"
	Symmetric1 Symmetry = "symmetric1"
	Symmetric2 Symmetry = "symmetric2"
)

func (s Symmetry) String() string {
	return string(s)
}

func NewPythagoreanScale() JustScale {
	return JustScale{
		system:      "Pythagorean",
		description: "3-limit Pythagorean ratios.",
		algorithm:   computePythagoreanIntervals,
	}
}

func New5LimitPythagoreanScale() JustScale {
	return JustScale{
		system:      "5-limit Pythagorean",
		description: "5-limit just intonation pure ratios chromatic scale derived from applying syntonic comma to Pythagorean ratios.",
		algorithm:   compute5LimitPythagoreanIntervals,
	}
}

func New5LimitJustIntonationChromaticScale(symmetry Symmetry) JustScale {
	return JustScale{
		system:      "5-limit Just Intonation",
		description: "5-limit just intonation pure ratios derived from third- and fifth-partial ratios.",
		algorithm: func() []Interval {
			return computeJustScale(buildMultiplierTablesFrom(multipliers(3), multipliers(5), multipliers(9)), fiveLimitScaleFilter(symmetry))
		},
	}
}

func NewSazScale() JustScale {
	// as per https://en.wikipedia.org/wiki/Ba%C4%9Flama and the cura that I have
	return JustScale{
		system:      "Saz",
		description: "Turkish Saz tuning ratios.",
		algorithm:   computeSazScale,
	}
}

// NewBachWohltemperierteKlavierScale Bach's Well Temperament as per
// Bach's Extraordinary Temperament: Our Rosetta Stone - Bradley Lehman
// Early Music Volume 33, Number 1, February 2005, pp. 3-23 (Article)
// Reference: https://academic.oup.com/em/article-abstract/33/1/3/535436?redirectedFrom=fulltext
func NewBachWohltemperierteKlavierScale() JustScale {
	return JustScale{
		system:      "Bach's Well-Tempered Tuning",
		description: "Derived from Lehman's decoding of Bach's Well-Tempered tuning, using sixth-comma, twelfth-comma, and pure fifths.",
		algorithm:   computeBachScale,
	}
}

func (s JustScale) System() string {
	return s.system
}

func (s JustScale) Description() string {
	return s.description
}

func (s JustScale) Intervals() []Interval {
	return s.algorithm()
}

type computeIntervalsFn func() []Interval
