package nightfury

import (
	"encoding/json"
	"fmt"
	"gitlab.com/jskswamy/nightfury/pkg/db"
)

var clientsBucketName = "clients"

// Client represent the client from where the games are started
type Client struct {
	Name         string
	Available    bool
	GameStatuses GameStatuses
}

// Clients represents the collection of Client
type Clients map[string]Client

// NewClient return a new instance of client with empty list of games
func NewClient(name string, available bool, game ...GameStatus) Client {
	games := GameStatuses{}
	for _, item := range game {
		games[item.Name] = item
	}
	return Client{
		Name:         name,
		Available:    available,
		GameStatuses: games,
	}
}

// NewClientFromRepoWithName return all the client from db
func NewClientFromRepoWithName(repo db.Repository, name string) (Client, error) {
	client := Client{}
	ok, err := repo.Fetch(clientsBucketName, name, &client)
	if err == nil {
		if ok {
			return client, nil
		}
		return client, db.EntryNotFound(fmt.Sprintf("client with name %v doesn't exists", name))

	}
	return client, err
}

// NewClientsFromRepo returns all the clients from db
func NewClientsFromRepo(repo db.Repository) (interface{}, error) {
	return repo.FetchAll(clientsBucketName, func(data []byte) (model db.Model, e error) {
		client := Client{}
		err := json.Unmarshal(data, &client)
		return client, err
	})
}

// ID returns the identifiable name for client
func (c Client) ID() string {
	return c.Name
}

// Add attaches a game to the client
func (c Client) Add(game Game) {
	c.GameStatuses[game.Name] = GameStatus{Name: game.Name, Status: Ready}
}

// Remove removes the game from the client
func (c Client) Remove(game Game) {
	delete(c.GameStatuses, game.Name)
}

// Connected marks the client as available
func (c Client) Connected() Client {
	c.Available = true
	return c
}

// Disconnected marks the client as unavailable
func (c Client) Disconnected() Client {
	c.Available = false
	return c
}

// Save saves the client information to db
func (c Client) Save(repo db.Repository) error {
	return repo.Save(clientsBucketName, c)
}

// Delete deletes the client information to db
func (c Client) Delete(repo db.Repository) error {
	return repo.Delete(clientsBucketName, c)
}
