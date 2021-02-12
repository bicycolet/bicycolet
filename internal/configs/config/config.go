package config

// Config defines config access.
type Config interface {
	// Change the values of this configuration Map.
	//
	// Return a map of key/value pairs that were actually changed. If
	// some keys fail to apply, details are included in the returned
	// ErrorList.
	Change(map[string]interface{}) (map[string]string, error)

	// Dump the current configuration held by this Map.
	//
	// Keys that match their default value will not be included in the dump. Also,
	// if a key has its Hidden attribute set to true, it will be rendered as
	// "true", for obfuscating the actual value.
	Dump(bool) (map[string]interface{}, error)

	// Key returns the value of the given key, which must be of type String.
	Key(string) (string, error)
}
