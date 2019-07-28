package root

// Server represents the structure for the server
type Server struct {
	Environment Environment            `json:"environment" yaml:"environment"`
	Config      map[string]interface{} `json:"config" yaml:"config"`
}

// Environment defines the server environment for the daemon
type Environment struct {
	Addresses     []string `json:"addresses" yaml:"addresses"`
	Server        string   `json:"server" yaml:"server"`
	ServerPid     int      `json:"server_pid" yaml:"server_pid"`
	ServerVersion string   `json:"server_version" yaml:"server_version"`
	ServerName    string   `json:"server_name" yaml:"server_name"`
}
