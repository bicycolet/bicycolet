package database

import "fmt"

type ConnectionInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Memory   bool
}

func (c ConnectionInfo) String() string {
	if c.Memory {
		return ":memory:"
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}
