package music

import "testing"

func TestMusicalMode_IsDiatonic(t *testing.T) {
	tests := []struct {
		name string
		m    MusicalMode
		want bool
	}{
		{
			name: "Lydian is diatonic",
			m:    Lydian,
			want: true,
		},
		{
			name: "Athenian is not a diatonic mode",
			m:    "Athenian",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsDiatonic(); got != tt.want {
				t.Errorf("IsDiatonic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMusicalMode_String(t *testing.T) {
	tests := []struct {
		name string
		m    MusicalMode
		want string
	}{
		{
			name: "ModeLydian string representation",
			m:    Lydian,
			want: "Lydian",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
