package nightfury

import (
	"encoding/json"
	"fmt"
	"gitlab.com/jskswamy/nightfury/pkg/db"
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

// NewSecurityIncidentsFromRepo returns all the clients from db
func NewSecurityIncidentsFromRepo(repo db.Repository) (interface{}, error) {
	return repo.FetchAll(securityIncidentBucketName, func(data []byte) (model db.Model, e error) {
		securityIncident := SecurityIncident{}
		err := json.Unmarshal(data, &securityIncident)
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
	return securityIncident, err
}

// ID returns the identifiable name for client
func (si SecurityIncident) ID() string {
	return si.Title
}

// Save saves the client information to db
func (si SecurityIncident) Save(repo db.Repository) error {
	return repo.Save(securityIncidentBucketName, si)
}

// Delete deletes the client information to db
func (si SecurityIncident) Delete(repo db.Repository) error {
	return repo.Delete(securityIncidentBucketName, si)
}