package main

import (
	"fmt"
	"net/http"

	"./preupload"
	"./upload"
)

func main() {
	http.HandleFunc("/api/upload", upload.Uploadhandle)
	http.HandleFunc("/api/pre", preupload.Prehandle)
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err)
	}
}
