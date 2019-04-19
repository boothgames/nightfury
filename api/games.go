package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/log"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

func listGames(c *gin.Context) {
	repository := db.DefaultRepository()
	games, err := nightfury.NewGamesFromRepo(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
}

func createGame(c *gin.Context) {
	game := nightfury.Game{}
	repository := db.DefaultRepository()
	err := c.ShouldBindJSON(&game)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = game.Save(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, game)
}

func populateGame(c *gin.Context) {
	gameName := c.Param("id")
	repository := db.DefaultRepository()
	game, err := nightfury.NewGameFromRepoWithName(repository, gameName)
	if err != nil {
		if entryNotFoundErr, ok := err.(db.EntryNotFound); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": entryNotFoundErr.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("game", game)
}

func readGame(c *gin.Context) {
	game, _ := c.Get("game")
	c.JSON(http.StatusOK, game)
}

func updateGame(c *gin.Context) {
	game, _ := c.Get("game")
	gameToBeUpdated := nightfury.Game{}
	currentGame := game.(nightfury.Game)
	err := c.ShouldBindJSON(&gameToBeUpdated)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if currentGame.Name != gameToBeUpdated.Name {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("name cannot be different")})
		return
	}

	repository := db.DefaultRepository()
	err = gameToBeUpdated.Save(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gameToBeUpdated)
}

func deleteGame(c *gin.Context) {
	game, _ := c.Get("game")
	gameToBeDeleted := game.(nightfury.Game)
	repository := db.DefaultRepository()
	err := gameToBeDeleted.Delete(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func gameConnected(session *melody.Session) {
	client, _, err := clientFromSession(session, func(id string) (client nightfury.Client, e error) {
		return nightfury.Client{}, fmt.Errorf("client not found")
	})
	if err != nil {
		logErr(err)
		err := session.Close()
		logErr(err)
		return
	}
	if !client.Available {
		log.Error("client not available, closing session")
		err := session.Close()
		logErr(err)
		return
	}
	game, repository, err := gameFromSession(session, func(id string) (nightfury.Game, error) {
		return nightfury.Game{Name: id}, nil
	})
	if err != nil {
		logErr(err)
		return
	}
	client.Add(*game)
	err = client.Save(repository)
	logErr(err)
}

func gameDisconnected(session *melody.Session) {
	client, _, err := clientFromSession(session, func(id string) (client nightfury.Client, e error) {
		return nightfury.Client{}, fmt.Errorf("client not found")
	})
	if err != nil {
		logErr(err)
		return
	}
	game, repository, err := gameFromSession(session, func(id string) (nightfury.Game, error) {
		return nightfury.Game{Name: id}, nil
	})
	if err != nil {
		logErr(err)
		return
	}
	client.Remove(*game)
	err = client.Save(repository)
	logErr(err)
}

func gameFromSession(session *melody.Session, notFoundFn func(name string) (nightfury.Game, error)) (*nightfury.Game, db.Repository, error) {
	if name, ok := gameName(session); ok {
		repository := db.DefaultRepository()
		game, err := nightfury.NewGameFromRepoWithName(repository, name)
		if _, ok := err.(db.EntryNotFound); ok {
			game, err := notFoundFn(name)
			if err != nil {
				return nil, nil, err
			}
			return &game, repository, nil
		} else if err != nil {
			return nil, nil, err
		}
		return &game, repository, nil
	}
	return nil, nil, fmt.Errorf("unable to parse game name from session")
}
