package testutils

import (
	"os"
)

// DBCreds are the required database credentials
type DBCreds struct {
	Port     int
	Dbname   string
	User     string
	Password string
	Host     string
}

// GetTestDBConnection returns a new gorm Database connection for the local docker instance
func GetTestDBConnection() DBCreds {

	var creds DBCreds
	if os.Getenv("CI") == "true" {
		creds = DBCreds{
			Port:     5432,
			Dbname:   "circle_test",
			User:     "root",
			Password: "root",
			Host:     "localhost",
		}
	} else {
		creds = DBCreds{
			Port:     5432,
			Dbname:   "civil_crawler",
			User:     "docker",
			Password: "docker",
			Host:     "localhost",
		}
	}

	return creds

}
