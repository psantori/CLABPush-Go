package server

// GCMHandler handles request for devices using GCM.
import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/contactlab/clabpush-go/utils"
)

// GCMHandler handles request for the GCM path.
type GCMHandler struct {
	ctx    *AppContext    // Shared AppContext.
	regexp *regexp.Regexp // Regexp to match the URL.
}

// NewGCMHandler return an handler with the given application context.
func NewGCMHandler(ctx *AppContext) *GCMHandler {
	regexp := regexp.MustCompile(`/gcm/devices/(?P<token>[\w:-]*)`)
	return &GCMHandler{ctx: ctx, regexp: regexp}
}

func (handler *GCMHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get the authorization key entry from the header and validate it. Send status 401 if
	// the key is invalid.
	key := r.Header.Get("Authorization")
	if !handler.ctx.ValidateAuthToken(key) {
		w.WriteHeader(http.StatusUnauthorized) // 401
		log.Println("401 - Unauthorized")
		return
	}

	path, _ := url.QueryUnescape(r.URL.Path) // We use the unescaped URL's path.

	if handler.regexp.MatchString(path) { // Match single device token.

		captures := utils.FindStringNamedSubmatches(handler.regexp, path)

		token := captures["token"]

		if r.Method == "PUT" {

			// Get the request body as string.
			buffer := new(bytes.Buffer)
			buffer.ReadFrom(r.Body)
			data := buffer.String()

			// Create/update the record in the database.
			query := "REPLACE INTO devices (token, vendor, data) VALUES ('" + token + "', 'GCM', '" + data + "')"
			_, err := handler.ctx.db.Exec(query)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				log.Println("500 - Internal server error")
				log.Println(err)
				return
			}
			w.WriteHeader(http.StatusCreated) // 201
			log.Println("201 - Created")
			return
		} else if r.Method == "DELETE" {

			query := "DELETE FROM devices WHETE token = '" + token + "'"
			_, err := handler.ctx.db.Exec(query)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				log.Println("500 - Internal server error")
				log.Println(err)
				return
			}

			w.WriteHeader(http.StatusCreated) // 201
			log.Println("200 - OK")
			return
		}
	}

	log.Println("400 - Bad Request")
	w.WriteHeader(http.StatusBadRequest) // 400
}
