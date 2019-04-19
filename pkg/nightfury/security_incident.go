package nightfury

import (
	"encoding/json"
	"fmt"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"strings"
)

var securityIncidentBucketName = "securityIncidents"

// SecurityIncident represents the security incident
type SecurityIncident struct {
	Title    string `binding:"required"`
	Tag      string `binding:"required"`
	Content  string `binding:"required"`
	Takeaway string `binding:"required"`
}

// SecurityIncidents represents collection of games
type SecurityIncidents map[string]SecurityIncident

func convertHyphenToSpaces(s string) string {
	return strings.Replace(s, "-", " ", -1)
}

// NewSecurityIncidentsFromRepo returns all the clients from db
func NewSecurityIncidentsFromRepo(repo db.Repository) (interface{}, error) {
	return repo.FetchAll(securityIncidentBucketName, func(data []byte) (model db.Model, e error) {
		securityIncident := SecurityIncident{}
		err := json.Unmarshal(data, &securityIncident)
		securityIncident.Title = convertHyphenToSpaces(securityIncident.Title)
		return securityIncident, err
	})
}

// NewSecurityIncidentFromRepoWithName return all the client from db
func NewSecurityIncidentFromRepoWithName(repo db.Repository, name string) (SecurityIncident, error) {
	securityIncident := SecurityIncident{}
	ok, err := repo.Fetch(securityIncidentBucketName, name, &securityIncident)
	if err == nil {
		if ok {
			return securityIncident, nil
		}
		return securityIncident, db.EntryNotFound(fmt.Sprintf("securityIncident with name %v doesn't exists", name))
	}
	securityIncident.Title = convertHyphenToSpaces(securityIncident.Title)
	return securityIncident, err
}

// ID returns the identifiable name for client
func (si SecurityIncident) ID() string {
	return strings.Replace(si.Title, " ", "-", -1)
}

// Save saves the client information to db
func (si SecurityIncident) Save(repo db.Repository) error {
	return repo.Save(securityIncidentBucketName, si)
}

// Delete deletes the client information to db
func (si SecurityIncident) Delete(repo db.Repository) error {
	return repo.Delete(securityIncidentBucketName, si)
}

// DetectChangeInTitle will return error if title changes during update
func (si SecurityIncident) DetectChangeInTitle(incidentToBeUpdated SecurityIncident) error {
	if si.Title != convertHyphenToSpaces(incidentToBeUpdated.Title) {
		return fmt.Errorf("title '%v' cannot be different", si.Title)
	}
	return nil
}
