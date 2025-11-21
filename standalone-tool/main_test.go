package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	b, _ := io.ReadAll(r)
	return string(b)
}

func Test_ShouldPlaceEqualTemperamentFretAtExpectedDistance(t *testing.T) {
	out := captureOutput(func() { equalTemperament(540.0, 12, 1) })
	if !strings.Contains(out, "Place fret # 1 at 30.308") {
		t.Fatalf("expected first equal-temperament fret to be at 30.308; got output:\n%s", out)
	}
}

func Test_ShouldComputeSimpleRatioPosition(t *testing.T) {
	out := captureOutput(func() { computeScaleFromRatios(540.0, [][]uint{{2, 1}}) })
	if !strings.Contains(out, "around 270.000") {
		t.Fatalf("expected ratio 2:1 to produce position 270.000; got output:\n%s", out)
	}
}

func Test_ShouldReduceOctaveToWithinOneAndTwo(t *testing.T) {
	if v := octaveReduce(4.0); v != 1.0 {
		t.Fatalf("octaveReduce(4.0) = %v; want 1.0", v)
	}
	if v := octaveReduce(0.5); v != 1.0 {
		t.Fatalf("octaveReduce(0.5) = %v; want 1.0", v)
	}
	if v := octaveReduce(1.5); v != 1.5 {
		t.Fatalf("octaveReduce(1.5) = %v; want 1.5", v)
	}
}

func Test_ShouldPrintOctaveAndNoteNamesInMeantone(t *testing.T) {
	out := captureOutput(func() { meantone(540.0, 0.25, false) })
	if !strings.Contains(out, "Open string") {
		t.Fatalf("expected meantone output to include 'Open string'; got:\n%s", out)
	}
	if !strings.Contains(out, "Place octave fret at 270.0") {
		t.Fatalf("expected meantone output to include octave placement at 270.0; got:\n%s", out)
	}
	if !strings.Contains(out, "(Eb)") && !strings.Contains(out, "(D#)") {
		t.Fatalf("expected meantone output to include a named note like Eb or D#; got:\n%s", out)
	}
}

func Test_ShouldPrintSazWithExpectedNumberOfPlaces(t *testing.T) {
	out := captureOutput(func() { saz(540.0) })
	// computeScaleFromRatios prints one "Place" line per ratio; saz uses 17 ratios
	count := strings.Count(out, "Place ")
	if count != 17 {
		t.Fatalf("expected 17 'Place' lines for saz ratios; got %d; full output:\n%s", count, out)
	}
}
