# Fret Placement Calculator

Outputs where to place the frets on a fretboard of a stringed instrument, effectively where to stop the strings, for various tunings, including:

* Just Intonation
* Quarter-Comma Meantone
* Extended Meantone
* Any Equal Temperament
* Pythagorean

It's agnostic of the actual open string tuning, tension of the string, type of instrument, etc.

# Examples

Compute just intonation for a scale length of 546mm:

`go run main.go -length=546`

Compute Pythagorean tuning for a scale length of 21 inches:

`go run main.go -length=21 -temper=pythagorean`

Computer extended quarter-comma meantone for a scale length of 546mm:

`go run main.go -length=546 -temper=meantone -extend=true`

Compute sixth-comma meantone:

`go run main.go -length=546 -temper=meantone -temper-by=0.166666667`

Compute 31 EDO:

`go run main.go -length=546 -temper=equal -divisions=31`

Compute fret placements for a saz/baglama:

`go run main.go -length=546 -temper=saz`


