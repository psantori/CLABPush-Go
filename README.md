# CLABPush-Go

ContactLab Push Notifications sample backend app for client push notifications management.

This repository contains a sample web app written in Go for [ContactLab](http://www.contactlab.com) push notifications management. It works in conjunction with the [CLABPush-Objective-C](https://github.com/contactlab/CLABPush-Objective-C) or [CLABPush-Swift](https://github.com/contactlab/CLABPush-Swift) sample code for iOS and [CLABPush-Android](https://github.com/contactlab/CLABPush-Android) for Android.

The project is composed by simple Go web application that will serve the requests from the Android and iOS applications and store the registration data in a database. In the package are included an exporter tool to export database record in csv format and an uploader tool to upload said file on the [ContactLab](http://www.contactlab.com) backend.

To facilitate testing, the project also includes a sample directory with a ready configuration file and an SQLite database. By launching the `clabpush-go` application from this directory, you have the sample app up and running.

## Disclaimer
This project does not represent a final product and is provided *as is* for demonstration only. For additional information or support, get in touch with [ContactLab](http://www.contactlab.com).

## Installation

First pull the CLABPush-Go project with

```bash
go get github.com/contactlab/clabpush-go
```

Then install the module with

```bash
go install github.com/contactlab/clabpush-go
```

Done.

## Usage

### Web application

If you want to run the app *as is*, move inside the sample folder in the project directory

```bash
cd $GOPATH/src/github.com/contactlab/clabpush-go/sample
```

Inside you'll find a `config.json.sample` file with a few parameters for the web application and a `clabpush.db.sample` SQLite database with a single devices table with an handful of columns. In order to launch the web server, make a copy of those files removing the `.sample` extension and modify them accordingly with your preferences.

```bash
cp clabpush.db.sample clabpush.db
cp config.json.sample config.json
```

To start the app, just use the `clabpush-go` command.

```bash
clabpush-go
```

### Exporting to csv

The exporter is not automatically installed when you install the `clabpush-go` application, but you can sort it out with

```bash
go install github.com/contactlab/clabpush-go/exporter
```

Once you have the tool installed, assuming you are in the sample folder of the project, you can just do

```bash
exporter -in clabpush.db -out export.csv
```

You can also pass in an username and a password if the database is protected

```bash
exporter -in clabpush.db -out export.csv -user your_username -password your_password
```

You should see in the console something along these lines

```bash
Connecting to clabpush.db...
Retrieving records...
Opening export.csv for output...
Exporting records...
Done!
```

### Upload the csv file to ContactLab

Like the exporter, the uploader is not installed by default, however it depends on the sftp module, that needs to be fetched first.

```bash
go get github.com/pkg/sftp
```

Then you can install the uploader

```bash
go install github.com/contactlab/clabpush-go/uploader
```

Assuming again you are in the sample directory, you can use it with

```bash
uploader -in export.csv -user your_username -secret your_secret -address sftp.example.com:22 -directory incoming/csvfiles -file exported.csv
```

That is a lot of stuff, and it will roughly do the equivalent of this

```bash
sftp your_username:your_secret@sftp.example.com:incoming/csvfiles
put export.csv
rename export.csv exported.csv
! touch ok.xml
put ok.xml
! rm ok.xml
exit
```

In other words, it will attempt to log in as *username* on the *sftp.example.com* port *22* and copy the `export.csv` file as `exported.csv` in the remote directory `incoming/csvfiles` (that needs to be already there). It will also create an empty `ok.xml` file in the same directory.

## Acknowledgments

- [ContactLab](http://www.contactlab.com)
- [DIMENSION](http://www.dimension.it)
- Paolo Santori, ContactLab
- Nicolò Tosi, DIMENSION
- Matteo Gavagnin, DIMENSION [@macteo](http://twitter.com/macteo)
- Daniele Dalledonne, DIMENSION
