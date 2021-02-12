package tokens

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/tokens/timer"
	"github.com/bicycolet/bicycolet/internal/resilience/timers/wall"
	"github.com/pkg/errors"
)

// TokenType defines the token we want to use for resilience.
type TokenType int

const (
	// Provision token that describes what token to use.
	Provision TokenType = iota
)

// New creates a new encoding token based on the type.
func New(t TokenType, capacity int64, freq time.Duration) (token.Token, error) {
	switch t {
	case Provision:
		return provision.New(capacity, freq), nil
	default:
		return nil, errors.Errorf("invalid token type %q", t)
	}
}
