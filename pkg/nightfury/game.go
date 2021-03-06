package nightfury

import (
	"encoding/json"
	"fmt"
	"github.com/boothgames/nightfury/pkg/db"
)

var gamesBucketName = "games"

// Game represents the game
type Game struct {
	Name        string                 `json:"name" binding:"required"`
	Title       string                 `json:"title"`
	Instruction string                 `json:"instruction" binding:"required"`
	Type        string                 `json:"type" binding:"required"`
	Mode        string                 `json:"mode"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Games represents collection of games
type Games map[string]Game

// NewGameFromRepoWithName return all the client from db
func NewGameFromRepoWithName(repo db.Repository, name string) (Game, error) {
	game := Game{}
	ok, err := repo.Fetch(gamesBucketName, Slug(name), &game)
	if err == nil {
		if ok {
			return game, nil
		}
		return game, db.EntryNotFound(fmt.Sprintf("game with name %v doesn't exists", name))
	}
	return game, err
}

// NewGamesFromRepo returns all the clients from db
func NewGamesFromRepo(repo db.Repository) (interface{}, error) {
	return repo.FetchAll(gamesBucketName, func(data []byte) (model db.Model, e error) {
		client := Game{}
		err := json.Unmarshal(data, &client)
		return client, err
	})
}

// ID returns the identifiable name for client
func (g Game) ID() string {
	return Slug(g.Name)
}

// Save saves the client information to db
func (g Game) Save(repo db.Repository) error {
	return repo.Save(gamesBucketName, g)
}

// Delete deletes the client information to db
func (g Game) Delete(repo db.Repository) error {
	return repo.Delete(gamesBucketName, g)
}
