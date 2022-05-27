package main

import (
	"fmt"
	"net/http"

	"./data"
	"./folder"
	"./preupload"
	"./upload"
)

func main() {
	http.HandleFunc("/api/upload", upload.Uploadhandle)
	http.HandleFunc("/api/pre", preupload.Prehandle)
	http.HandleFunc("/api/dta", data.Datahandle)
	http.HandleFunc("/api/folder", folder.Handlefolder)
	err := http.ListenAndServe(":22048", nil)
	if err != nil {
		fmt.Println(err)
	}
}
