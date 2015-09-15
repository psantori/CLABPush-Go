package server

// Context struct to share global data and objects.
import (
	"database/sql"
	"regexp"

	"github.com/contactlab/clabpush-go/utils"
)

// AppContext holds objects that would pollute the global space.
type AppContext struct {
	db            *sql.DB        // SQL database
	authKey       string         // Authorization key
	authKeyRegexp *regexp.Regexp // Auth key regexp
}

// NewAppContext return an application context.
func NewAppContext(db *sql.DB, authKey string) *AppContext {

	r := regexp.MustCompile(`Token (?P<token>[a-z0-9]*)`)
	return &AppContext{db: db, authKey: authKey, authKeyRegexp: r}
}

// ValidateAuthToken returns true if the key provided matches the one in context. String
// format must be 'key=<KEY_VALUE>'.
func (ctx *AppContext) ValidateAuthToken(key string) bool {

	if ctx.authKeyRegexp.MatchString(key) {
		captures := utils.FindStringNamedSubmatches(ctx.authKeyRegexp, key)
		keyValue := captures["token"]
		return ctx.authKey == keyValue
	}

	return false
}
