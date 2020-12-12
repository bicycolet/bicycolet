package info

// Envelope represents the structure for the server
type Envelope struct {
	Environment Environment            `json:"environment"`
	Config      map[string]interface{} `json:"config"`
}

// Environment defines the server environment for the daemon
type Environment struct {
	Addresses              []string `json:"addresses"`
	Certificate            string   `json:"certificate"`
	CertificateFingerprint string   `json:"certificate-fingerprint"`
	CertificateKey         string   `json:"certificate-key,omitempty"`
	Server                 string   `json:"server"`
	ServerPid              int      `json:"server-pid"`
	ServerVersion          string   `json:"server-version"`
	ServerClustered        bool     `json:"server-clustered"`
	ServerName             string   `json:"server-name"`
}
