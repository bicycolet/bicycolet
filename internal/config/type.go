package config

var booleans = []string{
	"true", "false",
	"1", "0",
	"yes", "no",
	"on", "off",
}
var truthy = []string{
	"true",
	"1",
	"yes",
	"on",
}

func contains(key string, list []string) bool {
	for _, entry := range list {
		if entry == key {
			return true
		}
	}
	return false
}
