package music

import (
	"reflect"
	"testing"
)

func TestInterval_String(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "String representation of perfect fifth",
			fields: fields{
				Numerator:   3,
				Denominator: 2,
				Name:        "Perfect Fifth",
			},
			want: "3:2",
		},
		{
			name: "String representation of septimal minor seventh",
			fields: fields{
				Numerator:   7,
				Denominator: 6,
				Name:        "Septimal Minor Seventh",
			},
			want: "7:6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_add(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	type args struct {
		other Interval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Interval
	}{
		{
			name: "Adding a Synthonic Comma to a lesser major second produces a greater major second",
			fields: fields{
				Numerator:   10,
				Denominator: 9,
				Name:        "Lesser Major Second",
			},
			args: args{
				other: Interval{
					Numerator:   81,
					Denominator: 80,
					Name:        "Synthonic Comma",
				},
			},
			want: Interval{
				Numerator:   9,
				Denominator: 8,
				Name:        "Pythagorean (Greater) Major Second",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.add(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_greaterThan(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	type args struct {
		other Interval
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
				Numerator:   9,
				Denominator: 8,
				Name:        "Greater Major Second",
			},
			args: args{
				other: Interval{
					Numerator:   10,
					Denominator: 9,
					Name:        "Lesser Major Second",
				},
			},
			want: true,
		},
		{
			name: "Greater than test for lesser major second and greater major second",
			fields: fields{
				Numerator:   10,
				Denominator: 9,
				Name:        "Lesser Major Second",
			},
			args: args{
				other: Interval{
					Numerator:   9,
					Denominator: 8,
					Name:        "Greater Major Second",
				},
			},
			want: false,
		},
		{
			name: "Greater than test for equal intervals",
			fields: fields{
				Numerator:   9,
				Denominator: 8,
				Name:        "Greater Major Second",
			},
			args: args{
				other: Interval{
					Numerator:   9,
					Denominator: 8,
					Name:        "Greater Major Second",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.greaterThan(tt.args.other); got != tt.want {
				t.Errorf("greaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isDiminishedFifth(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Diminished fifth test for diminished fifth interval",
			fields: fields{
				Numerator:   64,
				Denominator: 45,
				Name:        "Diminished Fifth",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.isDiminishedFifth(); got != tt.want {
				t.Errorf("isDiminishedFifth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isEqualTo(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	type args struct {
		other Interval
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
				Numerator:   3,
				Denominator: 2,
				Name:        "Perfect Fifth",
			},
			args: args{
				other: Interval{
					Numerator:   3,
					Denominator: 2,
					Name:        "Perfect Fifth",
				},
			},
			want: true,
		},
		{
			name: "Equality test for unequal intervals",
			fields: fields{
				Numerator:   3,
				Denominator: 2,
				Name:        "Perfect Fifth",
			},
			args: args{
				other: Interval{
					Numerator:   4,
					Denominator: 3,
					Name:        "Perfect Fourth",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.isEqualTo(tt.args.other); got != tt.want {
				t.Errorf("isEqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isGreaterMajorSecond(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.IsGreaterMajorSecond(); got != tt.want {
				t.Errorf("isGreaterMajorSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isGreaterMinorSeventh(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.IsGreaterMinorSeventh(); got != tt.want {
				t.Errorf("isGreaterMinorSeventh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isLesserMajorSecond(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.IsLesserMajorSecond(); got != tt.want {
				t.Errorf("isLesserMajorSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isLesserMinorSeventh(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.IsLesserMinorSeventh(); got != tt.want {
				t.Errorf("isLesserMinorSeventh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isOctave(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.isOctave(); got != tt.want {
				t.Errorf("isOctave() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isPerfect(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Perfect interval test for perfect unison",
			fields: fields{
				Numerator:   1,
				Denominator: 1,
				Name:        "Perfect Unison",
			},
			want: true,
		},
		{
			name: "Perfect interval test for perfect forth",
			fields: fields{
				Numerator:   4,
				Denominator: 3,
				Name:        "Perfect Forth",
			},
			want: true,
		},
		{
			name: "Perfect interval test for perfect fifth",
			fields: fields{
				Numerator:   3,
				Denominator: 2,
				Name:        "Perfect Fifth",
			},
			want: true,
		},
		{
			name: "Perfect interval test for perfect octave",
			fields: fields{
				Numerator:   2,
				Denominator: 1,
				Name:        "Perfect Octave",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.isPerfect(); got != tt.want {
				t.Errorf("isPerfect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isPerfectFifth(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.isPerfectFifth(); got != tt.want {
				t.Errorf("isPerfectFifth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isPerfectFourth(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.isPerfectFourth(); got != tt.want {
				t.Errorf("isPerfectFourth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_isUnison(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.isUnison(); got != tt.want {
				t.Errorf("isUnison() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_lessThan(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	type args struct {
		other Interval
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
				Numerator:   10,
				Denominator: 9,
				Name:        "Lesser Major Second",
			},
			args: args{
				other: Interval{
					Numerator:   9,
					Denominator: 8,
					Name:        "Greater Major Second",
				},
			},
			want: true,
		},
		{
			name: "Return false if interval is greater than another interval",
			fields: fields{
				Numerator:   9,
				Denominator: 8,
				Name:        "Greater Major Second",
			},
			args: args{
				other: Interval{
					Numerator:   10,
					Denominator: 9,
					Name:        "Lesser Major Second",
				},
			},
			want: false,
		},
		{
			name: "Return false if interval is equal to another interval",
			fields: fields{
				Numerator:   10,
				Denominator: 9,
				Name:        "Lesser Major Second",
			},
			args: args{
				other: Interval{
					Numerator:   10,
					Denominator: 9,
					Name:        "Lesser Major Second",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.lessThan(tt.args.other); got != tt.want {
				t.Errorf("lessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_name(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.NameMatch(); got != tt.want {
				t.Errorf("name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_octaveReduce(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	tests := []struct {
		name   string
		fields fields
		want   Interval
	}{
		{
			name: "Octave reduce a major ninth to a major second",
			fields: fields{
				Numerator:   9,
				Denominator: 4,
				Name:        "Major Ninth",
			},
			want: Interval{
				Numerator:   9,
				Denominator: 8,
				Name:        "Pythagorean (Greater) Major Second",
			},
		},
		{
			name: "Octave reduce a reciprocal prefect fourth to a perfect fourth",
			fields: fields{
				Numerator:   2,
				Denominator: 3,
				Name:        "Perfect Fifth Inversion",
			},
			want: Interval{
				Numerator:   4,
				Denominator: 3,
				Name:        "Perfect Fourth",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.octaveReduce(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("octaveReduce() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestInterval_reciprocal(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	tests := []struct {
		name   string
		fields fields
		want   Interval
	}{
		{
			name: "Reciprocal of an  perfect fifth is perfect fourth",
			fields: fields{
				Numerator:   81,
				Denominator: 80,
				Name:        "Pythagorean Comma",
			},
			want: Interval{
				Numerator:   80,
				Denominator: 81,
				Name:        "Pythagorean Comma Reciprocal",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.reciprocal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reciprocal() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestInterval_simplify(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	tests := []struct {
		name   string
		fields fields
		want   Interval
	}{
		{
			name: "Simplify a major tenth to a major third",
			fields: fields{
				Numerator:   15,
				Denominator: 8,
				Name:        "Major Tenth",
			},
			want: Interval{
				Numerator:   5,
				Denominator: 4,
				Name:        "Major Third",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.simplify(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simplify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_sortWith(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	type args struct {
		j Interval
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.sortWith(tt.args.j); got != tt.want {
				t.Errorf("sortWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_subtract(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	type args struct {
		other Interval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Interval
	}{
		{
			name: "Subtracting a lesser major second from a greater major second produces a syntonic comma",
			fields: fields{
				Numerator:   9,
				Denominator: 8,
				Name:        "Greater Major Second",
			},
			args: args{
				other: Interval{
					Numerator:   10,
					Denominator: 9,
					Name:        "Lesser Major Second",
				},
			},
			want: Interval{
				Numerator:   81,
				Denominator: 80,
				Name:        "Synthonic Comma",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.Subtract(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subtract() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestInterval_toCents(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.toCents(); got != tt.want {
				t.Errorf("toCents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_toFloat(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
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
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
			}
			if got := i.ToFloat(); got != tt.want {
				t.Errorf("toFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_toPowerOf(t *testing.T) {
	type fields struct {
		Numerator   uint
		Denominator uint
		Name        string
	}
	type args struct {
		p int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Interval
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				Numerator:   tt.fields.Numerator,
				Denominator: tt.fields.Denominator,
				Name:        tt.fields.Name,
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
		want Interval
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
		want []Interval
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
		want []Interval
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
		want Interval
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
		intervals []Interval
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
