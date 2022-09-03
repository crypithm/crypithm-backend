package main

import (
	"fmt"
	"net/http"

	"./data"
	"./folder"
	"./predown"
	"./preupload"
	"./share"
	"./viewShared"
)

func main() {
	http.HandleFunc("/api/pre", preupload.Prehandle)
	http.HandleFunc("/api/dta", data.Datahandle)
	http.HandleFunc("/api/folder", folder.Handlefolder)
	http.HandleFunc("/api/predown", predown.Predown)
	http.HandleFunc("/api/share", share.ShareHandler)
	http.HandleFunc("/api/viewshared", viewShared.SharedHandle)
	err := http.ListenAndServe(":22048", nil)
	if err != nil {
		fmt.Println(err)
	}
}
