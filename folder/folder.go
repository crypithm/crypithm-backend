package folder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	StatusMessage string
}

func Handlefolder(w http.ResponseWriter, r *http.Request) {
	var b []byte
	if r.Method != "POST" {
		var resp Response
		resp.StatusMessage = "Inallowed Method"
		b, _ = json.Marshal(resp)
	} else {
		token := r.Header.Get("Authorization")
		db, err := sql.Open("mysql", "crypithmusr:cDP9gNEQmUQt7qXbzU7XJ3Xz4mmcMf@tcp(127.0.0.1:3306)/crypithm")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		if err != nil {
			fmt.Println(err)
		}
		action := r.FormValue("action")
		rows, err := db.Query("SELECT uid FROM user WHERE token=?", token)
		defer rows.Close()
		if err != nil {
			fmt.Println(err)
		}
		if rows.Next() {
			if action == "create" {
				curentdirindex := r.FormValue("curentdirindex")
				foldername := r.FormValue("name")
				var uid string
				rows.Scan(&uid)
				indx, _ := strconv.Atoi(strings.Split(curentdirindex, " ")[1])
				var finindx int
				if strings.Split(curentdirindex, " ")[0] == "/" {
					finindx = 0
				} else {
					finindx = indx
				}
				ins, e := db.Query("INSERT INTO folder (name, folderindex, userid, date, size, parent, id) VALUES (?,?,?,?,?,?,UUID())", foldername, finindx, uid, time.Now().Format("2006-01-02 15:04:05"), 0, strings.Split(curentdirindex, " ")[0])
				if e != nil {
					fmt.Println(e)
				}
				defer ins.Close()
				dta := Response{"Success"}
				b, _ = json.Marshal(dta)
			} else if action == "delete" {

			} else if action == "move" {
			}
		} else {
			dta := Response{"Failed"}
			b, _ = json.Marshal(dta)
		}
	}

	fmt.Fprintf(w, "%s", string(b))
}
