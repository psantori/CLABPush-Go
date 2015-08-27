package server

// APNHandler handles the request for devices using APNS.
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"

	"github.com/contactlab/clabpush-go/utils"
)

// APNHandler handles request for the apn path.
type APNHandler struct {
	Handler                // Anonymous inner Handler.
	regexp  *regexp.Regexp // Regexp to match the URL.
}

// NewAPNHandler return an handler with the given application context.
func NewAPNHandler(ctx *AppContext) *APNHandler {
	regexp := regexp.MustCompile(`/apn/devices/(?P<token>[\w]*)`)
	return &APNHandler{Handler: Handler{Ctx: ctx}, regexp: regexp}
}

func (handler *APNHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Check out GCMHandler implementation of this very same method for explanations.

	key := r.Header.Get("Authorization")

	if !handler.Ctx.ValidateAuthToken(key) {
		handler.RenderHTTPStatus(w, http.StatusUnauthorized, nil) // 401
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
				handler.RenderHTTPStatus(w, http.StatusBadRequest, err) // 400
				return
			}

			userInfo, err := data.UserInfoAsString()
			if err != nil {
				handler.RenderHTTPStatus(w, http.StatusBadRequest, err)
				return
			}

			if _, err := handler.Ctx.db.Exec("REPLACE INTO devices (token, vendor, data, app_id, language) VALUES ($1, 'APN', $2, $3, $4)", token, userInfo, data.CLab.AppID, data.CLab.Language); err != nil {
				handler.RenderHTTPStatus(w, http.StatusInternalServerError, err) // 500
				return
			}

			handler.RenderHTTPStatus(w, http.StatusCreated, nil) // 201
			return

		} else if method == "DELETE" {

			if _, err := handler.Ctx.db.Exec("DELETE FROM devices WHERE token = $1", token); err != nil {
				handler.RenderHTTPStatus(w, http.StatusInternalServerError, err) // 500
				return
			}

			handler.RenderHTTPStatus(w, http.StatusOK, nil) // 200
			return
		}
	}

	handler.RenderHTTPStatus(w, http.StatusBadRequest, nil) // 400
}
