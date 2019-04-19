package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
	"net/http"
)

func listSecurityIncidents(c *gin.Context) {
	repository := db.DefaultRepository()
	games, err := nightfury.NewSecurityIncidentsFromRepo(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
}

func createSecurityIncident(c *gin.Context) {
	securityIncident := nightfury.SecurityIncident{}
	repository := db.DefaultRepository()
	err := c.ShouldBindJSON(&securityIncident)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = securityIncident.Save(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, securityIncident)
}

func populateSecurityIncident(c *gin.Context) {
	securityIncidentTitle := c.Param("id")
	repository := db.DefaultRepository()
	securityIncident, err := nightfury.NewSecurityIncidentFromRepoWithName(repository, securityIncidentTitle)
	if err != nil {
		if entryNotFoundErr, ok := err.(db.EntryNotFound); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": entryNotFoundErr.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("security_incident", securityIncident)
}

func readSecurityIncident(c *gin.Context) {
	securityIncident, _ := c.Get("security_incident")
	c.JSON(http.StatusOK, securityIncident)
}

func updateSecurityIncident(c *gin.Context) {
	securityIncident, _ := c.Get("security_incident")
	securityIncidentToBeUpdated := nightfury.SecurityIncident{}
	currentSecurityIncident := securityIncident.(nightfury.SecurityIncident)
	err := c.ShouldBindJSON(&securityIncidentToBeUpdated)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if currentSecurityIncident.Title != securityIncidentToBeUpdated.Title {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("title cannot be different").Error()})
		return
	}

	repository := db.DefaultRepository()
	err = securityIncidentToBeUpdated.Save(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, securityIncidentToBeUpdated)
}

func deleteSecurityIncident(c *gin.Context) {
	securityIncident, _ := c.Get("security_incident")
	securityIncidentToBeDeleted := securityIncident.(nightfury.SecurityIncident)
	repository := db.DefaultRepository()
	err := securityIncidentToBeDeleted.Delete(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
