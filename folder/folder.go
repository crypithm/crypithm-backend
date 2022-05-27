package folder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	StatusMessage string
}

type SucceedResp struct {
	StatusMessage string
	Id            string
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
				folderId := randstring(15)
				ins, e := db.Query("INSERT INTO folder (name, folderindex, userid, date, size, parent, id) VALUES (?,?,?,?,?,?,?)", foldername, finindx, uid, time.Now().Format("2006-01-02 15:04:05"), 0, strings.Split(curentdirindex, " ")[0], folderId)
				if e != nil {
					fmt.Println(e)
				}
				defer ins.Close()
				dta := SucceedResp{"Success", folderId}
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
func randstring(length int) string {

	var fin []string
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ1234567890"
	chars := strings.Split(str, "")
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length+1; i++ {
		fin = append(fin, chars[rand.Intn(26*2+10)])
	}
	return strings.Join(fin, "")
}
