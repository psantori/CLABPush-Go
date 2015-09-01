# CLABPush-Go
**ContactLab Push Notifications sample backend app for client push notifications management.**

This repository contains a simple GO web application for [ContactLab](http://www.contactlab.com) push notifications management. It works in conjunction with [CLABPush-Objective-C](https://github.com/contactlab/CLABPush-Objective-C) / [CLABPush-Swift](https://github.com/contactlab/CLABPush-Swift) for iOS and [CLABPush-Android](https://github.com/contactlab/CLABPush-Android) for Android.

The project explains how to manage device registration for iOS and Android and the process needed to synch data with [ContactLab](http://www.contactlab.com) backend.


## Disclaimer
This project does not represent a final product and is provided *as is* for demonstration use only. For additional information or support, get in touch with [ContactLab](http://www.contactlab.com).

## Install

First pull the CLABPush-Go project:

```bash
go get github.com/contactlab/clabpush-go
```

Then install the module:

```bash
go install github.com/contactlab/clabpush-go
```

## How to use

### Quick start

For a quick start, we provided a sample directory with a SQLite database and pre-configured settings.

Move to the project directory:

```bash
cd $GOPATH/src/github.com/contactlab/clabpush-go/sample
```

Copy and rename the sample files:

```bash
cp clabpush.db.sample clabpush.db
cp config.json.sample config.json
```

The SQLte database `clabpush.db` has the following columns:

- `id` an automatic incremental identifier
- `token` the device registration token
- `vendor` one of the following values:
 - `apn` for Apple Push Notification service (iOS)
 - `gcn` for Google Cloud Messaging (Android)
- `app_id` your mobile app Bundle Identifier (iOS) or Package Name (Android)
- `user_info` an optional JSON dictionary for user profiling

Edit `config.json` according to your preferences:

- `address` your web application IP address
- `port` your web application port
- `authKey` the authentication key to match your mobile app device registration
- `dbPath` yoour local SQLite database

To start the web application type:

```bash
clabpush-go
```

### Synch data with ContactLab

The easiest method to synch data with [ContactLab](http://www.contactlab.com) is by exporting the database in a comma-separated values (CSV) file, and uploading it to a SFTP server where [ContactLab](http://www.contactlab.com) backend will take care of.

### Exporting to CSV

To install the exporter type:

```bash
go install github.com/contactlab/clabpush-go/exporter
```

To start the exporter, if you are using our sample folder, type:

```bash
exporter -in clabpush.db -out export.csv
```
If you database requires authentication you can add username/password:  

```bash
exporter -in clabpush.db -out export.csv -user your_username -password your_password
```

If everything is correct you should see something like this:

```bash
Connecting to clabpush.db...
Retrieving records...
Opening export.csv for output...
Exporting records...
Done
```

### Upload the CSV

You can use your own favorite SFTP upload method or our simple uploader. To run the uploader you need to pull the SFTP module:

```bash
go get github.com/pkg/sftp
```

Then you can install the uploader:

```bash
go install github.com/contactlab/clabpush-go/uploader
```

To upload the CSV file, if you are in the sample folder, type:

```bash
uploader -in export.csv -user your_username -secret your_secret -address sftp.example.com:22 -directory incoming/daex -file exported.csv
```

This is the same as the following bash script:

```bash
sftp your_username:your_secret@sftp.example.com:incoming/daex
put export.csv
rename export.csv exported.csv
exit
```

In other words, it will attempt to log in as *username* on the *sftp.example.com* port *22* and copy the `export.csv` file as `exported.csv` in the remote directory `incoming/daex`.

### Schedule

According to your requirements, you can schedule a job to run periodically the exporter and uploader.

## Acknowledgments

- [ContactLab](http://www.contactlab.com)
- [DIMENSION](http://www.dimension.it)
- Paolo Santori, ContactLab
- Nicolò Tosi, DIMENSION
- Matteo Gavagnin, DIMENSION [@macteo](http://twitter.com/macteo)
- Daniele Dalledonne, DIMENSION
