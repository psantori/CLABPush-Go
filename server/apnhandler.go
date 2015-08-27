package server

// APNHandler handles the request for devices using APNS.
import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/contactlab/clabpush-go/utils"
)

// APNHandler handles request for the apn path.
type APNHandler struct {
	ctx    *AppContext
	regexp *regexp.Regexp
}

// NewAPNHandler return an handler with the given application context.
func NewAPNHandler(ctx *AppContext) *APNHandler {
	regexp := regexp.MustCompile(`/apn/devices/(?P<token>[\w]*)`)
	return &APNHandler{ctx: ctx, regexp: regexp}
}

func (handler *APNHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Check out GCMHandler implementation of this very same method for explanations.

	key := r.Header.Get("Authorization")

	if !handler.ctx.ValidateAuthToken(key) {
		w.WriteHeader(http.StatusUnauthorized) // 401
		log.Println("401 - Unauthorized")
		return
	}

	method := r.Method

	path, _ := url.QueryUnescape(r.URL.Path)

	if handler.regexp.MatchString(path) { // Match single device token

		captures := utils.FindStringNamedSubmatches(handler.regexp, path)

		token := captures["token"]

		if method == "PUT" {

			// Attempt to retrieve body object.
			data := NewData()

			buffer := new(bytes.Buffer)
			buffer.ReadFrom(r.Body)
			decoder := json.NewDecoder(buffer)
			if err := decoder.Decode(data); err != nil {
				w.WriteHeader(http.StatusBadRequest) // 400
				log.Println("400 - Bad Request")
				log.Println(err)
				return
			}

			if _, err := handler.ctx.db.Exec("REPLACE INTO devices (token, vendor, data, app_id, language) VALUES ($1, 'APN', $2, $3, $4)", token, data.UserInfo, data.CLab.AppID, data.CLab.Language); err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				log.Println("500 - Internal server error")
				log.Println(err)
				return
			}

			w.WriteHeader(http.StatusCreated) // 201
			log.Println("201 - Created")
			return

		} else if method == "DELETE" {

			if _, err := handler.ctx.db.Exec("DELETE FROM devices WHERE token = $1", token); err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				log.Println("500 - Internal server error")
				log.Println(err)
				return
			}

			w.WriteHeader(http.StatusOK) // 200
			log.Println("200 - OK")
			return
		}
	}

	log.Println("400 - Bad Request")
	w.WriteHeader(http.StatusBadRequest) // 400
}
