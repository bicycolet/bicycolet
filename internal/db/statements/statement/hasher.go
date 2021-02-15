package statement

// Hasher defines a way to hash sql statements.
type Hasher interface {
	Hash(string) string
}
