Feature: Fret placement calculator API should return details of where to place frets

  Scenario: API defaults to using Just Intonation with the specified scale length
    When I request where to put the frets for a scale length of 540
    Then I am provided with the fret placements for Just Intonation
