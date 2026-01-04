package music

import (
	"reflect"
	"testing"
)

func TestInterval_String(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "String representation of perfect fifth",
			fields: fields{
				numerator:   3,
				denominator: 2,
			},
			want: "3:2",
		},
		{
			name: "String representation of septimal minor seventh",
			fields: fields{
				numerator:   7,
				denominator: 6,
			},
			want: "7:6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_add(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	type args struct {
		other JustInterval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   JustInterval
	}{
		{
			name: "Adding a Synthonic Comma to a lesser major second produces a greater major second",
			fields: fields{
				numerator:   10,
				denominator: 9,
			},
			args: args{
				other: JustInterval{
					numerator:   81,
					denominator: 80,
				},
			},
			want: JustInterval{
				numerator:   9,
				denominator: 8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.add(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_greaterThan(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	type args struct {
		other JustInterval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Greater than test for greater major second and lesser major second",
			fields: fields{
				numerator:   9,
				denominator: 8,
			},
			args: args{
				other: JustInterval{
					numerator:   10,
					denominator: 9,
				},
			},
			want: true,
		},
		{
			name: "Greater than test for lesser major second and greater major second",
			fields: fields{
				numerator:   10,
				denominator: 9,
			},
			args: args{
				other: JustInterval{
					numerator:   9,
					denominator: 8,
				},
			},
			want: false,
		},
		{
			name: "Greater than test for equal intervals",
			fields: fields{
				numerator:   9,
				denominator: 8,
			},
			args: args{
				other: JustInterval{
					numerator:   9,
					denominator: 8,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.greaterThan(tt.args.other); got != tt.want {
				t.Errorf("greaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isDiminishedFifth(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Diminished fifth test for diminished fifth interval",
			fields: fields{
				numerator:   64,
				denominator: 45,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.isDiminishedFifth(); got != tt.want {
				t.Errorf("isDiminishedFifth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isEqualTo(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	type args struct {
		other JustInterval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Equality test for equal intervals",
			fields: fields{
				numerator:   3,
				denominator: 2,
			},
			args: args{
				other: JustInterval{
					numerator:   3,
					denominator: 2,
				},
			},
			want: true,
		},
		{
			name: "Equality test for unequal intervals",
			fields: fields{
				numerator:   3,
				denominator: 2,
			},
			args: args{
				other: JustInterval{
					numerator:   4,
					denominator: 3,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.isEqualTo(tt.args.other); got != tt.want {
				t.Errorf("isEqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isGreaterMajorSecond(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.IsGreaterMajorSecond(); got != tt.want {
				t.Errorf("isGreaterMajorSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isGreaterMinorSeventh(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.IsGreaterMinorSeventh(); got != tt.want {
				t.Errorf("isGreaterMinorSeventh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isLesserMajorSecond(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.IsLesserMajorSecond(); got != tt.want {
				t.Errorf("isLesserMajorSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isLesserMinorSeventh(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.IsLesserMinorSeventh(); got != tt.want {
				t.Errorf("isLesserMinorSeventh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isOctave(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.isOctave(); got != tt.want {
				t.Errorf("isOctave() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isPerfect(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Perfect interval test for perfect unison",
			fields: fields{
				numerator:   1,
				denominator: 1,
			},
			want: true,
		},
		{
			name: "Perfect interval test for perfect forth",
			fields: fields{
				numerator:   4,
				denominator: 3,
			},
			want: true,
		},
		{
			name: "Perfect interval test for perfect fifth",
			fields: fields{
				numerator:   3,
				denominator: 2,
			},
			want: true,
		},
		{
			name: "Perfect interval test for perfect octave",
			fields: fields{
				numerator:   2,
				denominator: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.isPerfect(); got != tt.want {
				t.Errorf("isPerfect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isPerfectFifth(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Perfect fifth test for perfect fifth interval",
			fields: fields{
				numerator:   3,
				denominator: 2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.isPerfectFifth(); got != tt.want {
				t.Errorf("isPerfectFifth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isPerfectFourth(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.isPerfectFourth(); got != tt.want {
				t.Errorf("isPerfectFourth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isUnison(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.isUnison(); got != tt.want {
				t.Errorf("isUnison() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_lessThan(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	type args struct {
		other JustInterval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Return true if interval is less than another interval",
			fields: fields{
				numerator:   10,
				denominator: 9,
			},
			args: args{
				other: JustInterval{
					numerator:   9,
					denominator: 8,
				},
			},
			want: true,
		},
		{
			name: "Return false if interval is greater than another interval",
			fields: fields{
				numerator:   9,
				denominator: 8,
			},
			args: args{
				other: JustInterval{
					numerator:   10,
					denominator: 9,
				},
			},
			want: false,
		},
		{
			name: "Return false if interval is equal to another interval",
			fields: fields{
				numerator:   10,
				denominator: 9,
			},
			args: args{
				other: JustInterval{
					numerator:   10,
					denominator: 9,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.lessThan(tt.args.other); got != tt.want {
				t.Errorf("lessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_name(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.Name(); got != tt.want {
				t.Errorf("name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_octaveReduce(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   JustInterval
	}{
		{
			name: "Octave reduce a major ninth to a major second",
			fields: fields{
				numerator:   9,
				denominator: 4,
			},
			want: JustInterval{
				numerator:   9,
				denominator: 8,
			},
		},
		{
			name: "Octave reduce a reciprocal prefect fourth to a perfect fourth",
			fields: fields{
				numerator:   2,
				denominator: 3,
			},
			want: JustInterval{
				numerator:   4,
				denominator: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.octaveReduce(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("octaveReduce() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestInterval_reciprocal(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   JustInterval
	}{
		{
			name: "Reciprocal of an  perfect fifth is perfect fourth",
			fields: fields{
				numerator:   81,
				denominator: 80,
			},
			want: JustInterval{
				numerator:   80,
				denominator: 81,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.reciprocal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reciprocal() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestInterval_simplify(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   JustInterval
	}{
		{
			name: "Simplify a major tenth to a major third",
			fields: fields{
				numerator:   15,
				denominator: 8,
			},
			want: JustInterval{
				numerator:   5,
				denominator: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.simplify(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simplify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_sortWith(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	type args struct {
		j JustInterval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.sortWith(tt.args.j); got != tt.want {
				t.Errorf("sortWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_subtract(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	type args struct {
		other JustInterval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   JustInterval
	}{
		{
			name: "Subtracting a lesser major second from a greater major second produces a syntonic comma",
			fields: fields{
				numerator:   9,
				denominator: 8,
			},
			args: args{
				other: JustInterval{
					numerator:   10,
					denominator: 9,
				},
			},
			want: JustInterval{
				numerator:   81,
				denominator: 80,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.Subtract(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subtract() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestInterval_toCents(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.toCents(); got != tt.want {
				t.Errorf("toCents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_toFloat(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.ToFloat(); got != tt.want {
				t.Errorf("toFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_toPowerOf(t *testing.T) {
	type fields struct {
		numerator   uint
		denominator uint
	}
	type args struct {
		p int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   JustInterval
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := JustInterval{
				numerator:   tt.fields.numerator,
				denominator: tt.fields.denominator,
			}
			if got := i.ToPowerOf(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toPowerOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createMultiplierTableOf(t *testing.T) {
	type args struct {
		multiplierListA [][]uint
		multiplierListB [][]uint
	}
	tests := []struct {
		name string
		args args
		want [][]uint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createMultiplierTableOf(tt.args.multiplierListA, tt.args.multiplierListB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createMultiplierTableOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromIntArray(t *testing.T) {
	type args struct {
		i []uint
	}
	tests := []struct {
		name string
		args args
		want JustInterval
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromIntArray(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromIntArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intervalsFromIntegers(t *testing.T) {
	type args struct {
		integers [][]uint
	}
	tests := []struct {
		name string
		args args
		want []JustInterval
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intervalsFromIntegers(tt.args.integers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intervalsFromIntegers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_justIntervalsFromMultipliers(t *testing.T) {
	type args struct {
		multiplierList [][]uint
		filter         intervalFilterFunction
	}
	tests := []struct {
		name string
		args args
		want []JustInterval
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := justIntervalsFromMultipliers(tt.args.multiplierList, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("justIntervalsFromMultipliers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_multipliers(t *testing.T) {
	type args struct {
		base uint
	}
	tests := []struct {
		name string
		args args
		want [][]uint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := multipliers(tt.args.base); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("multipliers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newInterval(t *testing.T) {
	type args struct {
		numerator   uint
		denominator uint
	}
	tests := []struct {
		name string
		args args
		want JustInterval
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newInterval(tt.args.numerator, tt.args.denominator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortIntervals(t *testing.T) {
	type args struct {
		intervals []JustInterval
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortIntervals(tt.args.intervals)
		})
	}
}
