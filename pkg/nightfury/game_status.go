package nightfury

import "fmt"

// Status represents different game status
type Status int

// String returns the string representation of status
func (status Status) String() string {
	return [...]string{"Ready", "InProgress", "Failed", "Completed"}[status]
}

const (
	// Ready represents game is available to play
	Ready Status = iota

	// InProgress represents game is in progress
	InProgress

	// Failed represents game has not been successfully completed
	Failed

	// Completed represents game has been successfully completed
	Completed
)

// GameStatus represents the game current status
type GameStatus struct {
	Name   string
	Status Status
}

// Failed mark the status as failed
func (g GameStatus) Failed() (GameStatus, error) {
	if g.Status == InProgress {
		g.Status = Failed
		return g, nil
	}
	return g, fmt.Errorf("cannot fail from a %v game", g.Status)
}

// Completed mark the status as completed
func (g GameStatus) Completed() (GameStatus, error) {
	if g.Status == InProgress {
		g.Status = Completed
		return g, nil
	}
	return g, fmt.Errorf("cannot complete from a %v game", g.Status)
}

// InProgress mark the status as progress
func (g GameStatus) InProgress() (GameStatus, error) {
	if g.Status == InProgress || g.Status == Ready {
		g.Status = InProgress
		return g, nil
	}
	return g, fmt.Errorf("cannot progress from a %v game", g.Status)
}

// GameStatuses represents the collection game current status
type GameStatuses map[string]GameStatus
