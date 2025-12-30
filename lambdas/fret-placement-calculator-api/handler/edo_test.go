package handler

import (
	"reflect"
	"testing"
)

func TestEquallyDividedOctave_division(t *testing.T) {
	type fields struct {
		NumberOfDivisions uint
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Division
	}{
		{
			name: "division returns ratio of 2 for the twelfth division of 12-EDO",
			fields: fields{
				NumberOfDivisions: 12,
			},
			args: args{
				i: 12,
			},
			want: Division{Ratio: 2, Cents: 1200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := EquallyDividedOctave{
				NumberOfDivisions: tt.fields.NumberOfDivisions,
			}
			if got := o.division(tt.args.i); got != tt.want {
				t.Errorf("division() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquallyDividedOctave_divisionAsCents(t *testing.T) {
	type fields struct {
		NumberOfDivisions uint
	}
	type args struct {
		i float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "division returns 1200 cents for the twelfth division of 12-EDO",
			fields: fields{
				NumberOfDivisions: 12,
			},
			args: args{
				i: 2.0,
			},
			want: 1200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := EquallyDividedOctave{
				NumberOfDivisions: tt.fields.NumberOfDivisions,
			}
			if got := o.divisionInCents(tt.args.i); got != tt.want {
				t.Errorf("divisionInCents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquallyDividedOctave_divisions(t *testing.T) {
	type fields struct {
		NumberOfDivisions uint
	}
	tests := []struct {
		name   string
		fields fields
		want   []Division
	}{
		{
			name:   "divisions returns the correct divisions for 12-EDO",
			fields: fields{NumberOfDivisions: 12},
			want: []Division{
				{1.0594630943592953, 100},
				{1.122462048309373, 200},
				{1.189207115002721, 300},
				{1.2599210498948732, 400},
				{1.3348398541700344, 500},
				{1.414213562373095, 600},
				{1.4983070768766815, 700},
				{1.5874010519681994, 800},
				{1.6817928305074292, 900},
				{1.7817974362806788, 1000},
				{1.887748625363387, 1100},
				{2, 1200},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := EquallyDividedOctave{
				NumberOfDivisions: tt.fields.NumberOfDivisions,
			}
			if got := o.divisions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("divisions() = %v, want %v", got, tt.want)
			}
		})
	}
}
