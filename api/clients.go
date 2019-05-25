package api

import (
	"github.com/boothgames/nightfury/pkg/db"
	"github.com/boothgames/nightfury/pkg/nightfury"
	"github.com/gin-gonic/gin"
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
