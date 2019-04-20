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

// Delete deletes all the client information from db
func (c Clients) Delete(repo db.Repository) error {
	repository := db.DefaultRepository()
	for _, client := range c {
		if err := client.Delete(repository); err != nil {
			return fmt.Errorf("delete failed for client %v, error: %v", client.Name, err)
		}
	}
	return nil
}

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

// Status represents Client status based on game status
func (c Client) Status() Status {
	statusCount := map[Status]int{}
	for _, gameStatus := range c.GameStatuses {
		statusCount[gameStatus.Status]++
	}
	gamesCount := len(c.GameStatuses)
	if statusCount[Ready] == gamesCount {
		return Ready
	}
	if statusCount[Completed] == gamesCount {
		return Completed
	}
	if statusCount[Failed] >= 1 {
		return Failed
	}
	return InProgress
}

// Save saves the client information to db
func (c Client) Save(repo db.Repository) error {
	return repo.Save(clientsBucketName, c)
}

// Delete deletes the client information to db
func (c Client) Delete(repo db.Repository) error {
	return repo.Delete(clientsBucketName, c)
}

// Start starts the first ready game, returns error if game is already started
func (c Client) Start() (Game, error) {
	if c.Status() == Ready {
		return c.startNextGame()
	}
	return Game{}, fmt.Errorf("game already started")
}

// HasNext checks if there is any game to play
func (c Client) HasNext() bool {
	if c.Status() != Failed && c.GameStatuses.HasReadyGames() {
		return true
	}
	return false
}

// Next returns next ready game
func (c Client) Next() (Game, error) {
	if c.Status() == Ready {
		return Game{}, fmt.Errorf("game not yet started")
	}

	if c.Status() == Completed {
		return Game{}, fmt.Errorf("game completed")
	}

	if c.Status() == Failed {
		return Game{}, fmt.Errorf("game failed")
	}
	if c.GameStatuses.IsAnyGameInProgress() {
		return Game{}, fmt.Errorf("game already in progress")
	}
	return c.startNextGame()
}

func (c Client) startNextGame() (Game, error) {
	game, err := c.GameStatuses.ReadyGame()
	if err != nil {
		return game, err
	}

	repository := db.DefaultRepository()
	gameStatus, err := c.GameStatuses[game.Name].InProgress()
	if err != nil {
		return game, err
	}
	c.GameStatuses[game.Name] = gameStatus
	err = c.Save(repository)
	return game, err
}

// CompleteGame completes a given game
func (c Client) CompleteGame(game Game) error {
	repository := db.DefaultRepository()
	gameStatus, err := c.GameStatuses[game.Name].Completed()
	if err != nil {
		return err
	}
	c.GameStatuses[game.Name] = gameStatus
	err = c.Save(repository)
	return err
}

// FailGame completes a given game
func (c Client) FailGame(game Game) error {
	repository := db.DefaultRepository()
	gameStatus, err := c.GameStatuses[game.Name].Failed()
	if err != nil {
		return err
	}
	c.GameStatuses[game.Name] = gameStatus
	err = c.Save(repository)
	return err
}

// Reset resets state of all games
func (c Client) Reset() error {
	repository := db.DefaultRepository()
	for name := range c.GameStatuses {
		c.GameStatuses[name] = GameStatus{Name: name, Status: Ready}
	}
	return c.Save(repository)
}
