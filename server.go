// Copyright 2012-2015 ContactLab, Italy
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"database/sql"
	"encoding/json"
	"github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"regexp"
	"net/url"
	"bytes"
)

// Context struct to share global data and objects.
type AppContext struct {
	db *sql.DB 						// SQL database
	authKey string 					// Authorization key
	authKeyRegexp *regexp.Regexp 	// Auth key regexp
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
		captures := FindStringNamedSubmatches(ctx.authKeyRegexp, key)
		keyValue := captures["token"]
		return ctx.authKey == keyValue
	}
	
	return false
}

// APNHandler handles the request for devices using APNS.
type APNHandler struct {
	ctx *AppContext
	regexp *regexp.Regexp
}

// NewAPNHandler return an handler with the given application context.
func NewAPNHandler(ctx *AppContext) *APNHandler {
	regexp := regexp.MustCompile(`/apn/devices/(?P<token>[\w]*)`)
	return &APNHandler{ctx:ctx, regexp:regexp}
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
	
		captures := FindStringNamedSubmatches(handler.regexp, path)
		
		token := captures["token"]
		
		if method == "PUT" {
		
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(r.Body)
		data := buffer.String()
	
		query := "REPLACE INTO devices (token, vendor, data) VALUES ('" + token + "', 'APN', '" + data + "')"
			_, err := handler.ctx.db.Exec(query)
			if err == nil {
				w.WriteHeader(http.StatusCreated) // 201
				log.Println("201 - Created")
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError) // 500
				log.Println("500 - Internal server error")
				log.Println(err)
				return
			}
			
		} else if method == "DELETE" {
		
		query := "DELETE FROM device WHERE token = '" + token + "'"
		_, err := handler.ctx.db.Exec(query)
		if err == nil {
				w.WriteHeader(http.StatusOK) // 200
				log.Println("200 - OK")
				return
		} else {
				w.WriteHeader(http.StatusInternalServerError) // 500
				log.Println("500 - Internal server error")
				log.Println(err)
				return
		}			
	}
	}	
	
	log.Println("400 - Bad Request")
	w.WriteHeader(http.StatusBadRequest) // 400
}

// GCMHandler handles request for devices using GCM.
type GCMHandler struct {
	ctx *AppContext 		// Shared AppContext.
	regexp *regexp.Regexp 	// Regexp to match the URL.
}

// NewGCMHandler return an handler with the given application context.
func NewGCMHandler(ctx *AppContext) *GCMHandler {
	regexp := regexp.MustCompile(`/gcm/devices/(?P<token>[\w:-]*)`)
	return &GCMHandler{ctx:ctx, regexp:regexp}
}

// FindStringNamedSubmatches returns a map of the named capturing groups of the provided
// regexp in the given string.
// Based on the example at http://blog.kamilkisiel.net/blog/2012/07/05/using-the-go-regexp-package/.
func FindStringNamedSubmatches(r *regexp.Regexp, s string) map[string]string {

  submatches := make(map[string]string)

  match := r.FindStringSubmatch(s)
  if match == nil {
      return submatches
  }

  for i, name := range r.SubexpNames() {
      // Ignore the whole regexp match and unnamed groups
      if i == 0 || name == "" {
          continue
      }
      
      submatches[name] = match[i]

  }
  return submatches
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
	
		captures := FindStringNamedSubmatches(handler.regexp, path)
		
		token := captures["token"]
		
		if r.Method == "PUT" {
		
		// Get the request body as string.
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(r.Body)
		data := buffer.String()
	
		// Create/update the record in the database.
		query := "REPLACE INTO devices (token, vendor, data) VALUES ('" + token + "', 'GCM', '" + data + "')"
			_, err := handler.ctx.db.Exec(query)
			if err == nil {
				w.WriteHeader(http.StatusCreated) // 201
				log.Println("201 - Created")
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError) // 500
				log.Println("500 - Internal server error")
				log.Println(err)
				return
			}
		} else if r.Method == "DELETE" {
		
			query := "DELETE FROM devices WHETE token = '" + token + "'"
			_, err := handler.ctx.db.Exec(query)
			if err == nil {
				w.WriteHeader(http.StatusCreated) // 201
				log.Println("200 - OK")
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError) // 500
				log.Println("500 - Internal server error")
				log.Println(err)
				return			
			}
		}	
	}
	
	log.Println("400 - Bad Request")
	w.WriteHeader(http.StatusBadRequest) // 400
}

// Config holds configuration file values.
type Config struct {
	Address string	// Server address.
	Port int		// Server port.
	AuthKey string	// Authorization key.
	DbPath string	// DB file path.
	DbName string 	// DB name.
}

// LoadConfig returns a Config struct from a config file.
func LoadConfig(path string) *Config {

	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	
	c := &Config{}
	err := decoder.Decode(c)
	if err != nil {
		log.Fatal(err.Error())
	} 
	return c
}

func main() {

	// Load and validate config file.
	config := LoadConfig("config.json")
	
	// Register SQLite driver.
	sqlite3conn := []*sqlite3.SQLiteConn{}
	sql.Register(config.DbName,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				sqlite3conn = append(sqlite3conn, conn)
				return nil
			},
		})

	// Open SQLite database and assign it to the contxt.
	db, err := sql.Open(config.DbName, config.DbPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	
	// Create the application context.
	ctx := NewAppContext(db, config.AuthKey)
	defer db.Close()

	// Create mux and register the handlers.
	mux := http.NewServeMux()
	mux.Handle("/gcm/devices/", NewGCMHandler(ctx))
	mux.Handle("/apn/devices/", NewAPNHandler(ctx))

	// Attemp to start the server.
	address := fmt.Sprintf("%s:%d", config.Address, config.Port)
	log.Fatal(http.ListenAndServe(address, mux))
}
