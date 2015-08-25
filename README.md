# CLABPush-Go
ContactLab Push Notifications sample backend app for client push notifications management.

The project is composed by simple Go web application that will serve the requests from
the Android and iOS libraries and store the registration data in a database. In the
package are included an exporter tool to export database record in csv format and an
uploader tool to upload said file on the Contactlab backend.

To facilitate testing, the project also includes a sample directory with a ready
configuration file and an SQLite database. By launching the clabpush-go application from
this directory, you have the sample app up and running.

# Installation

First pull the CLABPush-Go project with

    go get github.com/contactlab/clabpush-go

Then install the module with

    go install github.com/contactlab/clabpush.go

Done.

# Usage

## Web application

If you want to run the app as-is, move inside the sample folder in the
project directory

    cd $GOPATH/src/github.com/contactlab/clabpush-go/sample

Inside you'll find a config.json file with a few parameters for the web application
and a clabpush.db SQLite database with a single devices table with an handful of
columns.

    ➜  sample git:(master) ✗ ls
    clabpush.db config.json

To start the app, just launch the clabpush-go

    clabpush-go

## Exporting to csv

The exporter is not automatically installed when you install the clabpush-go application,
but you can sort it out with

    go install github.com/contactlab/clabpush-go/exporter

Once you have the tool installed, assuming you are in the sample folder of the project,
you can just do

    exporter -in clabpush.db -out export.csv

You can also pass in an username and a password if the database is protected

    exporter -in clabpush.db -out export.csv -user your_username -password your_password

You should see in the console something along these lines

    2015/08/25 17:59:56 Connecting to clabpush.db...
    2015/08/25 17:59:56 Retrieving records...
    2015/08/25 17:59:56 Opening export.csv for output...
    2015/08/25 17:59:56 Exporting records...
    2015/08/25 17:59:56 Done!

## Upload the csv file to Contactlab

Like the exporter, the uploader is not installed by default

    go install github.com/contactlab/clabpush-go/uploader

Assuming again you are in the sample directory, you can use it with

    uploader -in export.csv -user your_username -secret your_secret -address your.sftp.server -directory incoming/xmlfiles -file exported.csv

That is a lot of stuff, and it will roughly do the equivalent of this

    sftp your_username:your_secret@your.sftp.server:incoming/xmlfiles
    put export.csv
    rename export.csv exported.csv
    ! touch ok.xml
    put ok.xml
    ! rm ok.xml
    exit

In other words, it will attemp to log in as *your_username* on the *your.sftp.server.address*
and copy the *export.csv* file as *exported.csv* in the remote directory *incoming/xmlfiles*. It will
also create an empty *ok.xml* file in the same directory.

# Notes
