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
}

type Filedata struct {
	Id   string
	Name string
	Size string
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
		rows, _ := db.Query("SELECT uid FROM user WHERE token=?", token)
		defer rows.Close()
		if !rows.Next() {
			message, _ := json.Marshal(Defaultresp{"Error"})
			fmt.Fprintf(w, string(message))
		} else {
			var uid string
			rows.Scan(&uid)
			fileRows, _ := db.Query("SELECT size,name,id FROM files WHERE userid=?", uid)
			var returnData Fileresponse
			for fileRows.Next() {
				var fileJson Filedata
				fileRows.Scan(&fileJson.Size, &fileJson.Name, &fileJson.Id)
				returnData.Files = append(returnData.Files, fileJson)
			}
			returnData.Message = "Success"
			returnJSONarr, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnJSONarr))
		}
	}
}
