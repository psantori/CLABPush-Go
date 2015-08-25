package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// input is the path to the exported csv file to upload.
var input = flag.String("in", "out.csv", "the export file to upload")

// user name of the sftp account.
var user = flag.String("user", "", "the user")

// secret of the sftp account.
var secret = flag.String("secret", "", "the secret")

// address of the sftp server.
var address = flag.String("address", "", "the server address")

// directory is the directory to where upload the export file and the ok file.
var directory = flag.String("directory", "", "the remote directory")

// file is the remote name of the file that will be copied.
var file = flag.String("file", "", "the remote file")

func main() {

	// Attempt to open the input file
	flag.Parse()

	csv, err := os.OpenFile(*input, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalln(err)
	}
	defer csv.Close()
	reader := bufio.NewReader(csv)
	if err != nil {
		log.Fatalln(err)
	}

	// SSH connection.
	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			ssh.Password(*secret),
		},
	}

	conn, err := ssh.Dial("tcp", *address, config)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// SFTP client.
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatalln(err)
	}
	defer sftp.Close()

	// Upload over SFTP.
	sftpFile, err := sftp.OpenFile(filepath.Join(*directory, *file), os.O_WRONLY|os.O_CREATE)
	if err != nil {
		log.Fatalln(err)
	}
	defer sftpFile.Close()
	sftpFile.ReadFrom(reader)

	okFile, err := sftp.OpenFile(filepath.Join(*directory, "ok.xml"), os.O_WRONLY|os.O_CREATE)
	if err != nil {
		log.Fatalln(err)
	}
	defer okFile.Close()
}
