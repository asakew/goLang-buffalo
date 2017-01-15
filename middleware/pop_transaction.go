package middleware

import (
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

// PopTransaction is a piece of Buffalo middleware that wraps each
// request in a transaction that will automatically get committed or
// rolledback. It will also add a field to the log, "db", that
// shows the total duration spent during the reques making database
// calls.
var PopTransaction = func(db *pop.Connection) buffalo.MiddlewareFunc {
	return func(h buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			// wrap all requests in a transaction and set the length
			// of time doing things in the db to the log.

			err := db.Transaction(func(tx *pop.Connection) error {
				start := tx.Elapsed
				defer func() {
					finished := tx.Elapsed
					elapsed := time.Duration(finished - start)
					c.LogField("db", elapsed)
				}()
				c.Set("tx", tx)
				return h(c)
			})
			// find out if there is an underlying http error and return it rather than returning
			// the wrapped transaction error
			switch rootCause := errors.Cause(err).(type) {
			case buffalo.HTTPError:
				return rootCause
			default:
				return err
			}
		}
	}
}
