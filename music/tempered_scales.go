package music

import (
	"fmt"
	"math"
	"slices"
)

type TemperedScale struct {
	system      string
	description string
	algorithm   computeTemperedIntervalsFn
}

// NewBachWohltemperierteKlavierScale Bach's Well Temperament as per
// Bach's Extraordinary Temperament: Our Rosetta Stone - Bradley Lehman
// Early Music Volume 33, Number 1, February 2005, pp. 3-23 (Article)
// Reference: https://academic.oup.com/em/article-abstract/33/1/3/535436?redirectedFrom=fulltext
func NewBachWohltemperierteKlavierScale() TemperedScale {
	return TemperedScale{
		system:      "Bach's Well-Tempered Tuning",
		description: "Derived from Lehman's decoding of Bach's Well-Tempered tuning, using sixth-comma, twelfth-comma, and pure fifths.",
		algorithm:   computeBachScale,
	}
}

func NewQuarterCommaMeantoneScale() TemperedScale {
	return TemperedScale{
		system:      "Quarter-Comma Meantone",
		description: "Meantone temperament achieved by narrowing of fifths by 0.25 of a syntonic comma (81/80).",
		algorithm:   computeQuarterCommaMeantoneScale,
	}
}

func NewExtendedQuarterCommaMeantoneScale() TemperedScale {
	return TemperedScale{
		system:      "Extended Quarter-Comma Meantone",
		description: "Meantone temperament achieved by narrowing of fifths by 0.25 of a syntonic comma (81/80).",
		algorithm:   computeExtendedQuarterCommaMeantoneScale,
	}
}

func (s TemperedScale) System() string {
	return s.system
}

func (s TemperedScale) Description() string {
	return s.description
}

func (s TemperedScale) Intervals() []TemperedInterval {
	return s.algorithm()
}

type computeTemperedIntervalsFn func() []TemperedInterval

func computeQuarterCommaMeantoneScale() []TemperedInterval {
	return computeMeantoneScale(0.25, false)
}

func computeMeantoneScale(fractionOfSyntonicCommaToTemperFifthsBy float64, extendScale bool) []TemperedInterval {
	temperedFifth := PerfectFifth.ToFloat() * math.Pow(SyntonicComma.ToFloat(), -fractionOfSyntonicCommaToTemperFifthsBy)

	var ratiosOfNotesToFundamental = []TemperedInterval{2.0}
	var fifthsFromTonic = 6
	if extendScale {
		fifthsFromTonic = 9
	}

	for i := -fifthsFromTonic; i <= fifthsFromTonic; i++ {
		ratiosOfNotesToFundamental = append(ratiosOfNotesToFundamental, TemperedInterval(toDecimalPlaces(octaveReduceFloat(math.Pow(temperedFifth, float64(i))), 3)))
	}

	slices.Sort(ratiosOfNotesToFundamental)
	return ratiosOfNotesToFundamental
}

func computeExtendedQuarterCommaMeantoneScale() []TemperedInterval {
	return computeMeantoneScale(0.25, true)
}

func octaveReduceFloat(ratio float64) float64 {
	for ratio >= 2.0 || ratio < 1.0 {
		if ratio >= 2.0 {
			ratio /= 2.0
		}
		if ratio < 1.0 {
			ratio *= 2.0
		}
	}
	return ratio
}

func computeBachScale() []TemperedInterval {
	// Narrowing of the fifths as outlined by Lehman
	temperingFractions := []float64{
		0.0,         // Pure fifth
		-1.0 / 12.0, // Twelfth-comma narrowed
		-1.0 / 6.0,  // Sixth-comma narrowed
		0.0,         // Pure fifth
		-1.0 / 6.0,  // Sixth-comma narrowed
		-1.0 / 12.0, // Twelfth-comma narrowed
		0.0,         // Pure fifth
		-1.0 / 6.0,  // Sixth-comma narrowed
		-1.0 / 12.0, // Twelfth-comma narrowed
		0.0,         // Pure fifth
		-1.0 / 6.0,  // Sixth-comma narrowed
		-1.0 / 12.0, // Twelfth-comma narrowed
	}

	// Calculate tempered fifths
	temperedFifths := make([]float64, 12)
	for i := 0; i < 12; i++ {
		temperedFifths[i] = 3.0 / 2.0 * math.Pow(SyntonicComma.ToFloat(), temperingFractions[i])
	}

	// Derive ratios by walking the circle of fifths
	ratios := make([]float64, 12)
	ratios[0] = 1.0 // Start with the tonic
	for i := 1; i < 12; i++ {
		ratios[i] = ratios[i-1] * temperedFifths[(i-1)%12]
	}

	// Reduce ratios to within the octave [1.0, 2.0)
	for i := range ratios {
		ratios[i] = octaveReduceFloat(ratios[i])
	}

	ratios = append(ratios, 2.0)

	slices.Sort(ratios) // Sort the ratios in ascending order
	var intervals []TemperedInterval
	for _, ratio := range ratios {
		intervals = append(intervals, TemperedInterval(toDecimalPlaces(ratio, 3)))
	}
	return intervals
}

func toDecimalPlaces(v float64, p int) float64 {
	exponent := math.Pow10(p)
	return math.Round(v*exponent) / exponent
}

func NewEqualTemperamentScale(divisionsOfOctave uint) TemperedScale {
	return TemperedScale{
		system:      "Equal Temperament",
		description: fmt.Sprintf("%d-tone equal temperament.", divisionsOfOctave),
		algorithm: func() []TemperedInterval {
			intervals := []TemperedInterval{1.0}
			for i := 1; i <= int(divisionsOfOctave); i++ {
				intervals = append(intervals, TemperedInterval(math.Exp2(float64(i)/float64(divisionsOfOctave))))
			}
			return intervals
		},
	}
}
