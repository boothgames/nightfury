package nightfury

import (
	"encoding/json"
	"fmt"
	"github.com/boothgames/nightfury/pkg/db"
	"strings"
)

var hintBucketName = "hints"

// Hint represents the hint
type Hint struct {
	Title    string   `json:"title" binding:"required"`
	Tag      []string `json:"tag"  binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Takeaway string   `json:"takeaway" binding:"required"`
}

// Hints represents collection of games
type Hints map[string]Hint

func convertHyphenToSpaces(s string) string {
	return strings.Replace(s, "-", " ", -1)
}

// NewHintsFromRepo returns all the clients from db
func NewHintsFromRepo(repo db.Repository) (interface{}, error) {
	return repo.FetchAll(hintBucketName, func(data []byte) (model db.Model, e error) {
		hint := Hint{}
		err := json.Unmarshal(data, &hint)
		hint.Title = convertHyphenToSpaces(hint.Title)
		return hint, err
	})
}

// NewHintFromRepoWithName return all the client from db
func NewHintFromRepoWithName(repo db.Repository, name string) (Hint, error) {
	hint := Hint{}
	ok, err := repo.Fetch(hintBucketName, Slug(name), &hint)
	if err == nil {
		if ok {
			return hint, nil
		}
		return hint, db.EntryNotFound(fmt.Sprintf("hint with name %v doesn't exists", name))
	}
	hint.Title = convertHyphenToSpaces(hint.Title)
	return hint, err
}

// ID returns the identifiable name for client
func (hint Hint) ID() string {
	return Slug(hint.Title)
}

// Save saves the client information to db
func (hint Hint) Save(repo db.Repository) error {
	return repo.Save(hintBucketName, hint)
}

// Delete deletes the client information to db
func (hint Hint) Delete(repo db.Repository) error {
	return repo.Delete(hintBucketName, hint)
}

// DetectChangeInTitle will return error if title changes during update
func (hint Hint) DetectChangeInTitle(incidentToBeUpdated Hint) error {
	if hint.Title != convertHyphenToSpaces(incidentToBeUpdated.Title) {
		return fmt.Errorf("title '%v' cannot be different", hint.Title)
	}
	return nil
}
