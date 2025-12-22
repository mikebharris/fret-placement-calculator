Feature: Fret placement calculator API should return details of where to place frets

  Scenario: API defaults to using Ptolemy's Intense Diatonic Scale with the specified scale length
    When I request where to put the frets for a scale length of 540
    Then I am provided with the fret placements for 5-limit just intonation chromatic scale
