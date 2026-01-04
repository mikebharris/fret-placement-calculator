package music

const (
	Lydian     MusicalMode = "Lydian"
	Ionian     MusicalMode = "Ionian"
	Mixolydian MusicalMode = "Mixolydian"
	Dorian     MusicalMode = "Dorian"
	Aeolian    MusicalMode = "Aeolian"
	Phrygian   MusicalMode = "Phrygian"
	Locrian    MusicalMode = "Locrian"
)

type MusicalMode string

func (m MusicalMode) IsDiatonic() bool {
	switch m {
	case Lydian, Ionian, Mixolydian, Dorian, Aeolian, Phrygian, Locrian:
		return true
	default:
		return false
	}
}
func (m MusicalMode) String() string {
	return string(m)
}
