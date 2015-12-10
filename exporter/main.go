package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-sqlite3"
)

// inputFile is the path to the database.
var inputFile = flag.String("in", "clabpush.db", "the database to export")

// outputFile is the path to the output csv file.
var outputFile = flag.String("out", "out.csv", "the output file")

// dbUser is the database dbUser to use (UNUSED).
var dbUser = flag.String("usr", "", "the database user")

// dbPassword is the database password to use (UNUSED).
var dbPassword = flag.String("pwd", "", "the database password")

// dbName is the name for the connection pool. You can ignore this.
var dbName = flag.String("dbname", "clabpush.exporter.SQLITE", "the database name")

// batchfile is the name for the sftp batch file that will be generated.
var batchFile = flag.String("bfile", "sftp_batch_file", "the sftp batch file for the upload")

var configFile = flag.String("cfg", "config.json", "the configuration file")

// Config holds the configuration options for the exporter.
type config struct {
	Fields []string
}

// LoadConfig returns a Config struct from a config file.
func LoadConfig(path string) *config {

	c := &config{}

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return c
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		log.Println(err)
	}

	return c
}

// Once the exporter tool is installed, just invoke the command with the
// appropriate parameters, like:
// exporter -in yourdatabase.db -out youroutputfile.csv
func main() {

	// Parse the flags from the command line.
	flag.Parse()

	cfg := LoadConfig(*configFile)

	ctx := &Context{}

	// Use FieldsExporter if there are export fields. Default to RawExporter.
	if len(cfg.Fields) > 0 {
		exporter := new(RawExporter)
		ctx.Exporter = exporter
	} else {
		exporter := NewFieldExporter(cfg.Fields)
		ctx.Exporter = exporter
	}

	// CSV Params.
	ctx.Out.CSVFile = *outputFile

	// DB params.
	ctx.In.Path = *inputFile
	ctx.In.User = *dbUser
	ctx.In.Password = *dbPassword

	// Prepare the db connection pool.
	sqlite3conn := []*sqlite3.SQLiteConn{}
	sql.Register("contactlab.push.exporter.SQLITE",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				sqlite3conn = append(sqlite3conn, conn)
				return nil
			},
		})

	// Open a connection to the database.
	log.Printf("Connecting to %s...\n", ctx.In.Path)
	db, err := sql.Open("contactlab.push.exporter.SQLITE", ctx.In.Path)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	ctx.In.Connection = db

	err = prepareCSV(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Done!")
}

// CSVRecord interface can be implemented to object serialized to CSV.
type CSVRecord interface {
	ToCSV(r *CSVRecord) []string
}

// Device represents a row in the database.
type Device struct {
	Token    string
	Vendor   string
	AppID    string
	Language string
	UserInfo string
}

// Exporter provide ToCSV method to export Device data to CSV format.
type Exporter interface {
	ToCSV(*Device) []string
}

// RawExporter export the fields of Device, including UserInfo as string.
type RawExporter struct {
}

// ToCSV implements Exporter interface ToCSV to export the UserInfo as string.
func (ex *RawExporter) ToCSV(d *Device) []string {
	var record []string
	record = append(record, d.Token)
	record = append(record, d.Vendor)
	record = append(record, d.AppID)
	record = append(record, d.Language)
	record = append(record, d.UserInfo)
	return record
}

// FieldExporter export the fields of Devices as they are, except for UserInfo
// where it attemps to extract specific values.
type FieldExporter struct {
	Fields [][]string
}

// NewFieldExporter creates a new FieldExporter with the given export fields.
func NewFieldExporter(fields []string) *FieldExporter {
	exporter := new(FieldExporter)
	for i, path := range fields {
		exporter.Fields[i] = strings.Split(path, ".")
	}
	return exporter
}

// findValueAtPath attempt to export the object at the path specified for the
// given info map.
func (ex *FieldExporter) findValueAtPath(path []string, info map[string]interface{}) string {
	m := info
	last := (len(path) - 1)
	for i, name := range path {
		if i == last {
			if s, ok := m[name].(string); ok {
				return s
			}
		} else {
			if next, ok := m[name].(map[string]interface{}); ok {
				m = next
			} else {
				return "" // Break
			}
		}
	}
	return ""
}

// ToCSV implementation that export specific fields from UserInfo to CSV.
func (ex *FieldExporter) ToCSV(r *Device) []string {
	var record []string
	record = append(record, r.Token)
	record = append(record, r.Vendor)
	record = append(record, r.AppID)
	record = append(record, r.Language)

	// If we have at least one field to export, we need to unmarhsal the UserInfo
	// to an arbitrary map[string]interface{}
	if len(ex.Fields) > 0 {

		var info map[string]interface{}
		err := json.Unmarshal([]byte(r.UserInfo), &info)

		if err != nil {
			log.Println(err)
			return record
		}

		for _, path := range ex.Fields {
			v := ex.findValueAtPath(path, info)
			record = append(record, v)
		}
	}

	return record
}

// ToCSV return a string slice to pass to the csv package writing functions.
func (r *Device) ToCSV() []string {
	var record []string
	record = append(record, r.Token)
	record = append(record, r.Vendor)
	record = append(record, r.AppID)
	record = append(record, r.Language)
	record = append(record, r.UserInfo)
	return record
}

// NewDevice returns a new Record.
func NewDevice() *Device {
	return new(Device)
}

// Context struct to avoid global variables pollution.
type Context struct {
	Out      Output
	In       Input
	Exporter Exporter
}

// Input holds database access info.
type Input struct {
	Path       string  // Path to the database
	User       string  // Database user
	Password   string  // Database password
	Connection *sql.DB // Reference to the database connection
}

// Output holds output information.
type Output struct {
	CSVFile string
	OKFile  string
}

func prepareCSV(ctx *Context) error {

	// Collect all the data we need from the database.
	log.Println("Retrieving records...")
	rows, err := ctx.In.Connection.Query("SELECT token, vendor, app_id, language, user_info FROM devices")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Open the output file.
	log.Printf("Opening %s for output...", ctx.Out.CSVFile)
	file, err := os.OpenFile(ctx.Out.CSVFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	// Create a CSV writer and dump the records in it.
	writer := csv.NewWriter(file)
	log.Println("Exporting records...")
	for rows.Next() {
		d := NewDevice()
		rows.Scan(&d.Token, &d.Vendor, &d.AppID, &d.Language, &d.UserInfo)
		if err := writer.Write(ctx.Exporter.ToCSV(d)); err != nil {
			return err
		}
	}

	writer.Flush()
	return nil
}
