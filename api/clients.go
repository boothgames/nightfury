package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
	"net/http"
)

func listClients(c *gin.Context) {
	repository := db.DefaultRepository()
	clients, err := nightfury.NewClientsFromRepo(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clients)
}
