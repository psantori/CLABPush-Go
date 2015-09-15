package server

// GCMHandler handles request for devices using GCM.
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"

	"github.com/contactlab/clabpush-go/utils"
)

// GCMHandler handles request for the GCM path.
type GCMHandler struct {
	Handler                // Anonymous inner Handler.
	regexp  *regexp.Regexp // Regexp to match the URL.
}

// NewGCMHandler return an handler with the given application context.
func NewGCMHandler(ctx *AppContext) *GCMHandler {
	regexp := regexp.MustCompile(`/gcm/devices/(?P<token>[\w:-]*)`)
	return &GCMHandler{Handler: Handler{Ctx: ctx}, regexp: regexp}
}

func (handler *GCMHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get the authorization key entry from the header and validate it. Send status 401 if
	// the key is invalid.
	key := r.Header.Get("Authorization")
	if !handler.Ctx.ValidateAuthToken(key) {
		handler.RenderHTTPStatus(w, http.StatusUnauthorized, nil)
		return
	}

	path, _ := url.QueryUnescape(r.URL.Path) // We use the unescaped URL's path.

	if handler.regexp.MatchString(path) { // Match single device token.

		captures := utils.FindStringNamedSubmatches(handler.regexp, path)

		token := captures["token"]

		if r.Method == "PUT" {

			// Get the request body and unmarshal it.
			data := NewData()

			buffer := new(bytes.Buffer)
			buffer.ReadFrom(r.Body)
			decoder := json.NewDecoder(buffer)

			if err := decoder.Decode(data); err != nil {
				handler.RenderHTTPStatus(w, http.StatusBadRequest, err)
				return
			}

			userInfo, err := data.UserInfoAsString()
			if err != nil {
				handler.RenderHTTPStatus(w, http.StatusBadRequest, err)
				return
			}

			// Create/update the record in the database.
			if _, err := handler.Ctx.db.Exec("REPLACE INTO devices (token, vendor, user_info, app_id, language) VALUES ($1, 'GCM', $2, $3, $4)", token, userInfo, data.CLab.AppID, data.CLab.Language); err != nil {
				handler.RenderHTTPStatus(w, http.StatusInternalServerError, err)
				return
			}

			handler.RenderHTTPStatus(w, http.StatusCreated, nil)
			return

		} else if r.Method == "DELETE" {

			if _, err := handler.Ctx.db.Exec("DELETE FROM devices WHERE token = $1", token); err != nil {
				handler.RenderHTTPStatus(w, http.StatusInternalServerError, err)
				return
			}

			handler.RenderHTTPStatus(w, http.StatusOK, nil)
			return
		}
	}

	handler.RenderHTTPStatus(w, http.StatusBadRequest, nil)
}
