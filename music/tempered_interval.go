package music

import (
	"fmt"
	"math"
)

type TemperedInterval float64

func (i TemperedInterval) toCents() float64 {
	return math.Log10(i.toFloat()) / math.Log10(2) * 1200
}

func (i TemperedInterval) toFloat() float64 {
	return float64(i)
}

func (i TemperedInterval) String() string {
	return fmt.Sprintf("%f", i.toFloat())
}
