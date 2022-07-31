package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Defaultresp struct {
	Message string
}
type Fileresponse struct {
	Message string
	Files   []interface{}
	Folders []interface{}
}

type Filedata struct {
	Id   string
	Name string
	Size string
	Dir  string
}

type AppendedFileData struct {
	Folders  []interface{}
	Username string
	Message  string
}

type Folderdata struct {
	Id    string
	Name  string
	Index string
	Date  string
}

func Datahandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		message, _ := json.Marshal(Defaultresp{"Inallowed Method"})
		fmt.Fprintf(w, string(message))
		return
	} else {
		db, err := sql.Open("mysql", "crypithmusr:cDP9gNEQmUQt7qXbzU7XJ3Xz4mmcMf@tcp(127.0.0.1:3306)/crypithm")
		if err != nil {
			log.Fatal(err)
		}
		token := r.Header.Get("Authorization")
		defer db.Close()
		rows, _ := db.Query("SELECT uid, username FROM user WHERE token=?", token)
		defer rows.Close()
		if !rows.Next() {
			message, _ := json.Marshal(Defaultresp{"Error"})
			fmt.Fprintf(w, string(message))
		} else {
			var uid, username string
			rows.Scan(&uid, &username)
			if r.FormValue("action") == "getOnlyFolder" {
				folderRows, _ := db.Query("SELECT name,id, date, parent FROM folder WHERE userid=?", uid)
				//Folderdata
				var folders AppendedFileData
				for folderRows.Next() {
					var folderjson Folderdata
					folderRows.Scan(&folderjson.Name, &folderjson.Id, &folderjson.Date, &folderjson.Index)
					folders.Folders = append(folders.Folders, folderjson)
				}
				folders.Username = username
				folders.Message = "Success"
				returnJSONarr, _ := json.Marshal(folders)
				fmt.Fprintf(w, string(returnJSONarr))
			} else {
				fileRows, _ := db.Query("SELECT size,name,id, directory FROM files WHERE userid=?", uid)
				var returnData Fileresponse
				for fileRows.Next() {
					var fileJson Filedata
					fileRows.Scan(&fileJson.Size, &fileJson.Name, &fileJson.Id, &fileJson.Dir)
					returnData.Files = append(returnData.Files, fileJson)
				}
				folderRows, _ := db.Query("SELECT name,id, date, parent FROM folder WHERE userid=?", uid)
				//Folderdata
				for folderRows.Next() {
					var folderjson Folderdata
					folderRows.Scan(&folderjson.Name, &folderjson.Id, &folderjson.Date, &folderjson.Index)
					returnData.Folders = append(returnData.Folders, folderjson)
				}
				returnData.Message = "Success"
				returnJSONarr, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnJSONarr))
			}
		}
	}
}
