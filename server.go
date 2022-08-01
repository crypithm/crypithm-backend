package main

import (
	"fmt"
	"net/http"

	"./data"
	"./download"
	"./folder"
	"./predown"
	"./preupload"
	"./share"
	"./upload"
	"./viewShared"
)

func main() {
	http.HandleFunc("/api/upload", upload.Uploadhandle)
	http.HandleFunc("/api/pre", preupload.Prehandle)
	http.HandleFunc("/api/dta", data.Datahandle)
	http.HandleFunc("/api/folder", folder.Handlefolder)
	http.HandleFunc("/api/predown", predown.Predown)
	http.HandleFunc("/api/download", download.Downloader)
	http.HandleFunc("/api/share", share.ShareHandler)
	http.HandleFunc("/api/viewshared", viewShared.SharedHandle)
	err := http.ListenAndServe(":22048", nil)
	if err != nil {
		fmt.Println(err)
	}
}
