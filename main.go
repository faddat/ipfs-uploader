//This came from:
//https://tutorialedge.net/golang/go-file-upload-tutorial/
//https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
//https://github.com/ipfs/js-ipfs/tree/master/examples/ipfs-101
// But I eventuually gave up on actuuallly using the client librarries from IPFS and decidded to just POST it to a locally exposed IPFS gateway, that makes the brain hurt less and doesn't have any securrityy tradeoffs.


package main

import (
	"bytes"
	"fmt"
	shell "github.com/ipfs/go-ipfs-http-api"
	"io/ioutil"
	"mime/multipart"
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

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("cheating", "cheating")
	if err != nil {
		fmt.Println(err)
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}



	part.Write(fileBytes)

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)


	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":6969", nil)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
}
