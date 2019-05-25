package api

import (
	"github.com/boothgames/nightfury/pkg/db"
	"github.com/boothgames/nightfury/pkg/nightfury"
	"github.com/gin-gonic/gin"
	"net/http"
)

const hintContextKey = "hint"

func listHints(c *gin.Context) {
	repository := db.DefaultRepository()
	hints, err := nightfury.NewHintsFromRepo(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hints)
}

func createHint(c *gin.Context) {
	hint := nightfury.Hint{}
	repository := db.DefaultRepository()
	err := c.ShouldBindJSON(&hint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = hint.Save(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, hint)
}

func uploadHints(c *gin.Context) {
	var hints []nightfury.Hint
	repository := db.DefaultRepository()
	err := c.ShouldBindJSON(&hints)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, hint := range hints {
		err = hint.Save(repository)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusCreated, hints)
}

func populateHint(c *gin.Context) {
	hintTitle := c.Param("id")
	repository := db.DefaultRepository()
	hint, err := nightfury.NewHintFromRepoWithName(repository, hintTitle)
	if err != nil {
		if entryNotFoundErr, ok := err.(db.EntryNotFound); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": entryNotFoundErr.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set(hintContextKey, hint)
}

func readHint(c *gin.Context) {
	hint, _ := c.Get(hintContextKey)
	c.JSON(http.StatusOK, hint)
}

func updateHint(c *gin.Context) {
	hint, _ := c.Get(hintContextKey)
	hintToBeUpdated := nightfury.Hint{}
	currentHint := hint.(nightfury.Hint)
	err := c.ShouldBindJSON(&hintToBeUpdated)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := currentHint.DetectChangeInTitle(hintToBeUpdated); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repository := db.DefaultRepository()
	err = hintToBeUpdated.Save(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hintToBeUpdated)
}

func deleteHint(c *gin.Context) {
	hint, _ := c.Get(hintContextKey)
	hintToBeDeleted := hint.(nightfury.Hint)
	repository := db.DefaultRepository()
	err := hintToBeDeleted.Delete(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
