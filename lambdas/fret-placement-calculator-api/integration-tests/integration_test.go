package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/testcontainernetwork-go"
	"github.com/stretchr/testify/assert"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	var steps steps
	steps.t = t

	suite := godog.TestSuite{
		TestSuiteInitializer: func(ctx *godog.TestSuiteContext) {
			ctx.BeforeSuite(steps.startContainerNetwork)
			ctx.AfterSuite(steps.stopContainerNetwork)
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Step(`^I request where to put the frets for a scale length of (\d+)$`, steps.iRequestWhereToPutTheFretsForAScaleLengthOf)
			ctx.Step(`^I am provided with the fret placements for Just Intonation$`, steps.iAmProvidedWithTheFretPlacementsForJustIntonation)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

type steps struct {
	t                   *testing.T
	networkOfContainers testcontainernetwork.NetworkOfDockerContainers
	lambdaContainer     testcontainernetwork.LambdaDockerContainer
}

var responseFromLambda events.LambdaFunctionURLResponse

func (s *steps) startContainerNetwork() {
	s.lambdaContainer = testcontainernetwork.LambdaDockerContainer{
		Config: testcontainernetwork.LambdaDockerContainerConfig{
			Hostname:    "lambda",
			Executable:  "../main",
			Environment: map[string]string{},
		},
	}

	s.networkOfContainers =
		testcontainernetwork.NetworkOfDockerContainers{}.
			WithDockerContainer(&s.lambdaContainer)
	_ = s.networkOfContainers.StartWithDelay(5 * time.Second)
}

func (s *steps) stopContainerNetwork() {
	if err := s.networkOfContainers.Stop(); err != nil {
		log.Fatalf("stopping docker containers: %v", err)
	}
}

func (s *steps) iRequestWhereToPutTheFretsForAScaleLengthOf(scaleLength string) {
	s.invokeLambdaUsingRequest(events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": scaleLength}})
}

func (s *steps) invokeLambdaUsingRequest(request events.LambdaFunctionURLRequest) {
	requestJsonBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("marshalling lambda request %v", err)
	}
	response, err := http.Post(s.lambdaContainer.InvocationUrl(), "application/json", bytes.NewReader(requestJsonBytes))
	if err != nil {
		log.Fatalf("triggering lambda: %v", err)
	}

	if response.StatusCode != 200 {
		log.Fatalf("invoking Lambda: %d", response.StatusCode)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(response.Body); err != nil {
		log.Fatalf("reading response body: %v", err)
	}

	if err := json.Unmarshal(buf.Bytes(), &responseFromLambda); err != nil {
		log.Fatalf("unmarshalling response: %v", err)
	}
}

type Fret struct {
	Label    string  `json:"label"`
	Position float64 `json:"position"`
	Comment  string  `json:"comment,omitempty"`
}

type FretPlacements struct {
	System      string  `json:"system"`
	Description string  `json:"description, omitempty"`
	ScaleLength float64 `json:"scaleLength"`
	Frets       []Fret  `json:"frets"`
}

func (s *steps) iAmProvidedWithTheFretPlacementsForJustIntonation() error {
	assert.Equal(s.t, responseFromLambda.Headers["Content-Type"], "application/json")
	assert.Equal(s.t, http.StatusOK, responseFromLambda.StatusCode)

	fretPlacements := FretPlacements{}
	if err := json.Unmarshal([]byte(responseFromLambda.Body), &fretPlacements); err != nil {
		return fmt.Errorf("unmarshalling result: %s", err)
	}

	assert.Equal(s.t, "ji", fretPlacements.System)
	assert.Equal(s.t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(s.t, 7, len(fretPlacements.Frets))

	assert.Equal(s.t, "9:8", fretPlacements.Frets[0].Label)
	assert.Equal(s.t, 60.00, fretPlacements.Frets[0].Position)

	assert.Equal(s.t, "2:1", fretPlacements.Frets[6].Label)
	assert.Equal(s.t, 270.0, fretPlacements.Frets[6].Position)
	return nil
}
