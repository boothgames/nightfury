package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
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
