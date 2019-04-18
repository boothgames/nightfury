package pkg

import "gitlab.com/jskswamy/nightfury/pkg/db"

var gamesBucketName = "games"

// Game represents the game
type Game struct {
	Name        string
	Instruction string
	Type        string
}

// NewClientFromRepoWithName return all the client from db
func NewGameFromRepoWithName(repo db.Repository, name string) (Game, error) {
	game := Game{}
	err := repo.Fetch(gamesBucketName, name, &game)
	return game, err
}

// ID returns the identifiable name for client
func (g Game) ID() string {
	return g.Name
}

// Save saves the client information to db
func (g Game) Save(repo db.Repository) error {
	return repo.Save(gamesBucketName, g)
}

// GameStatus represents the game current status
type GameStatus struct {
	Name   string
	Status string
}

// GameStatuses represents the collection game current status
type GameStatuses map[string]GameStatus
