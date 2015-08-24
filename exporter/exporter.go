package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"
)

// inputPath is the path to the database.
var inputPath = flag.String("in", "clabpush.db", "the database to export")

// outputPath is the path to the output csv file.
var outputPath = flag.String("out", "out.csv", "the output file")

// username is the database username to use (UNUSED).
var username = flag.String("usr", "", "the database user")

// password is the database password to use (UNUSED).
var password = flag.String("pwd", "", "the database password")

// databaseName is the name for the connection pool. You can ignore this.
var databaseName = flag.String("dbname", "clabpush.exporter.SQLITE", "the database name")

// Record represents a row in the database.
type Record struct {
	Token  string
	Vendor string
}

// ToCSV return a string slice to pass to the csv package writing functions.
func (r *Record) ToCSV() []string {
	var record []string
	record = append(record, r.Token)
	record = append(record, r.Vendor)
	return record
}

// NewRecord returns a new Record.
func NewRecord() *Record {
	return new(Record)
}

func main() {

	// Parse the flags from the command line.
	flag.Parse()

	// Prepare the db connection pool.
	sqlite3conn := []*sqlite3.SQLiteConn{}
	sql.Register(*databaseName,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				sqlite3conn = append(sqlite3conn, conn)
				return nil
			},
		})

	// Open a connection to the database.
	log.Printf("Connecting to %s...\n", *inputPath)
	db, err := sql.Open(*databaseName, *inputPath)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	// Collect all the data we need from the database.
	log.Println("Retrieving records...")
	rows, err := db.Query("SELECT token, vendor FROM devices")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer rows.Close()

	// Open the output file.
	log.Printf("Opening %s for output...", *outputPath)
	file, err := os.OpenFile(*outputPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a CSV writer and dump the records in it.
	writer := csv.NewWriter(file)
	log.Println("Exporting records...")
	for rows.Next() {
		record := NewRecord()
		rows.Scan(&record.Token, &record.Vendor)
		if err := writer.Write(record.ToCSV()); err != nil {
			log.Fatalln(err)
		}
	}
	writer.Flush()
	log.Println("Done!")
}
