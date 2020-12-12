// +build linux,cgo

package responses

import (
	"database/sql"
	"os"

	"github.com/pkg/errors"
)

// SmartError returns the right error message based on err.
func SmartError(err error) Reply {
	if err == nil {
		return EmptySyncResponse()
	}

	switch errors.Cause(err) {
	case os.ErrNotExist, sql.ErrNoRows:
		if errors.Cause(err) != err {
			return NotFound(err)
		}
		return NotFound(nil)

	case os.ErrPermission:
		if errors.Cause(err) != err {
			return Forbidden(err)
		}
		return Forbidden(nil)

	default:
		return InternalError(err)
	}
}
