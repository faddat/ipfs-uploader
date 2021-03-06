//This came from:
//https://tutorialedge.net/golang/go-file-upload-tutorial/
//https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
//https://github.com/ipfs/js-ipfs/tree/master/examples/ipfs-101
// So what this does is it provides an endpoint where users can upload files, and it does not write those files to disk.  Instead, pumps them into ipfs, which does write them to disk.

package main

import (
	"bytes"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"net/http"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	//The magic trick that makes it work.  Take from []byte too a byte buffer.
	z := bytes.NewBuffer(fileBytes)

	// Connect to ipfs API
	sh := shell.NewShell("localhost:5001")
	
	// Send the byte buffer to the ipfs API
	cid, err := sh.Add(z)
	
	// Output to terminal on server
	fmt.Println("added", cid)

	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, cid)
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":6969", nil)
}

func main() {
	fmt.Println("Starting Blurt's file upload service")
	setupRoutes()
}
