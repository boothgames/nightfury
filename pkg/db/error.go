package db

// EntryNotFound represents a not found entity error
type EntryNotFound string

// Error returns the error string
func (e EntryNotFound) Error() string {
	return string(e)
}
