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
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/contactlab/clabpush-go/server"
	"github.com/mattn/go-sqlite3"
)

// Config holds configuration file values.
type Config struct {
	Address string // Server address.
	Port    int    // Server port.
	AuthKey string // Authorization key.
	DbPath  string // DB file path.
	DbName  string // DB name.
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
	ctx := server.NewAppContext(db, config.AuthKey)
	defer db.Close()

	// Create mux and register the handlers.
	mux := http.NewServeMux()
	mux.Handle("/gcm/devices/", server.NewGCMHandler(ctx))
	mux.Handle("/apn/devices/", server.NewAPNHandler(ctx))

	// Attemp to start the server.
	address := fmt.Sprintf("%s:%d", config.Address, config.Port)
	log.Fatal(http.ListenAndServe(address, mux))
}
