package handler

import "testing"

func Test_parseIntegerQueryParameter(t *testing.T) {
	type args struct {
		q        map[string]string
		key      string
		fallback int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid integer parameter",
			args: args{
				q:        map[string]string{"octaves": "3"},
				key:      "octaves",
				fallback: 1,
			},
			want: 3,
		},
		{
			name: "missing parameter uses fallback",
			args: args{
				q:        map[string]string{},
				key:      "octaves",
				fallback: 2,
			},
			want: 2,
		},
		{
			name: "invalid integer parameter uses fallback",
			args: args{
				q:        map[string]string{"octaves": "invalid"},
				key:      "octaves",
				fallback: 4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseIntegerQueryParameter(tt.args.q, tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("parseIntegerQueryParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}
