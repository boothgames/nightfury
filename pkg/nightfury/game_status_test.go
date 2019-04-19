package nightfury

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type gameStatusScenario struct {
	name           string
	status         Status
	expectedStatus Status
	isError        bool
	errMsg         string
}

func TestGameStatusFailed(t *testing.T) {
	failedGameStatusScenarios := []gameStatusScenario{
		{
			name:           "should be able to fail a in-progress game",
			status:         InProgress,
			expectedStatus: Failed,
		},
		{
			name:           "should not fail a ready game",
			status:         Ready,
			expectedStatus: Ready,
			isError:        true,
			errMsg:         "cannot fail from a Ready game",
		},
		{
			name:           "should not fail a failed game",
			status:         Failed,
			expectedStatus: Failed,
			isError:        true,
			errMsg:         "cannot fail from a Failed game",
		},
		{
			name:           "should not fail a not available game",
			status:         NotAvailable,
			expectedStatus: NotAvailable,
			isError:        true,
			errMsg:         "cannot fail from a NotAvailable game",
		}, {
			name:           "should not fail a completed game",
			status:         Completed,
			expectedStatus: Completed,
			isError:        true,
			errMsg:         "cannot fail from a Completed game",
		},
	}

	for _, scenario := range failedGameStatusScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			status := GameStatus{Name: "game", Status: scenario.status}

			actual, err := status.Failed()

			if scenario.isError {
				assert.Error(t, err)
				assert.Equal(t, scenario.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, GameStatus{Name: "game", Status: scenario.expectedStatus}, actual)
		})
	}
}

func TestGameStatusCompleted(t *testing.T) {
	completedGameStatusScenarios := []gameStatusScenario{
		{
			name:           "should be able to progress a in-progress game",
			status:         InProgress,
			expectedStatus: Completed,
		},
		{
			name:           "should not complete a ready game",
			status:         Ready,
			expectedStatus: Ready,
			isError:        true,
			errMsg:         "cannot complete from a Ready game",
		},
		{
			name:           "should not complete a failed game",
			status:         Failed,
			expectedStatus: Failed,
			isError:        true,
			errMsg:         "cannot complete from a Failed game",
		},
		{
			name:           "should not complete a not available game",
			status:         NotAvailable,
			expectedStatus: NotAvailable,
			isError:        true,
			errMsg:         "cannot complete from a NotAvailable game",
		}, {
			name:           "should not complete a completed game",
			status:         Completed,
			expectedStatus: Completed,
			isError:        true,
			errMsg:         "cannot complete from a Completed game",
		},
	}

	for _, scenario := range completedGameStatusScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			status := GameStatus{Name: "game", Status: scenario.status}

			actual, err := status.Completed()

			if scenario.isError {
				assert.Error(t, err)
				assert.Equal(t, scenario.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, GameStatus{Name: "game", Status: scenario.expectedStatus}, actual)
		})
	}
}

func TestGameStatusInProgress(t *testing.T) {
	inProgressGameStatusScenarios := []gameStatusScenario{
		{
			name:           "should be able to progress a in-progress game",
			status:         InProgress,
			expectedStatus: InProgress,
		},
		{
			name:           "should progress a ready game",
			status:         Ready,
			expectedStatus: InProgress,
		},
		{
			name:           "should not progress a failed game",
			status:         Failed,
			expectedStatus: Failed,
			isError:        true,
			errMsg:         "cannot progress from a Failed game",
		},
		{
			name:           "should not progress a not available game",
			status:         NotAvailable,
			expectedStatus: NotAvailable,
			isError:        true,
			errMsg:         "cannot progress from a NotAvailable game",
		}, {
			name:           "should not progress a completed game",
			status:         Completed,
			expectedStatus: Completed,
			isError:        true,
			errMsg:         "cannot progress from a Completed game",
		},
	}

	for _, scenario := range inProgressGameStatusScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			status := GameStatus{Name: "game", Status: scenario.status}

			actual, err := status.InProgress()

			if scenario.isError {
				assert.Error(t, err)
				assert.Equal(t, scenario.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, GameStatus{Name: "game", Status: scenario.expectedStatus}, actual)
		})
	}
}
