package folder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
		fmt.Fprintf(w, string(b))
		return
	}
	token := r.Header.Get("Authorization")
	db, err := sql.Open("mysql", "crypithmusr:cDP9gNEQmUQt7qXbzU7XJ3Xz4mmcMf@tcp(127.0.0.1:3306)/crypithm")
	if err != nil {
		dta := Response{"Failed"}
		b, _ = json.Marshal(dta)
		fmt.Fprintf(w, string(b))
		return
	}
	defer db.Close()
	if err != nil {
		dta := Response{"Failed"}
		b, _ = json.Marshal(dta)
		fmt.Fprintf(w, string(b))
		return
	}
	action := r.FormValue("action")
	rows, err := db.Query("SELECT uid FROM user WHERE token=?", token)
	defer rows.Close()
	if err != nil {
		dta := Response{"Failed"}
		b, _ = json.Marshal(dta)
		fmt.Fprintf(w, string(b))
		return
	}
	if !rows.Next() {
		dta := Response{"Failed"}
		b, _ = json.Marshal(dta)
		fmt.Fprintf(w, string(b))
		return
	}
	if action == "create" {
		curentdirindex := r.FormValue("curentdirindex")
		foldername := r.FormValue("name")
		var uid string
		rows.Scan(&uid)
		folderId := randstring(15)
		ins, e := db.Query("INSERT INTO folder (name, userid, date, size, parent, id) VALUES (?,?,?,?,?,?)", foldername, uid, time.Now().Format("2006-01-02 15:04:05"), 0, curentdirindex, folderId)
		if e != nil {
			fmt.Println(e)
		}
		defer ins.Close()
		dta := SucceedResp{"Success", folderId}
		b, _ = json.Marshal(dta)
	} else if action == "delete" {

	} else if action == "move" {

		itemList := r.FormValue("targetObjs")
		target := r.FormValue("target")
		var arr []string
		var uid string
		rows.Scan(&uid)
		_ = json.Unmarshal([]byte(itemList), &arr)
		query1, args, err := sqlx.In("UPDATE folder SET parent=? WHERE id IN (?) AND userid=?", target, arr, uid)
		if err != nil {
			dta := Response{"Failed"}
			b, _ = json.Marshal(dta)
		}
		query2, args2, err := sqlx.In("UPDATE files SET directory=? WHERE id IN (?) AND userid=?", target, arr, uid)
		if err != nil {
			dta := Response{"Failed"}
			b, _ = json.Marshal(dta)
		}
		_, e := db.Exec(query1, args...)
		_, e = db.Exec(query2, args2...)
		if e != nil {
			dta := Response{"Failed"}
			b, _ = json.Marshal(dta)
		}
		dta := Response{"Success"}
		b, _ = json.Marshal(dta)
	}
	fmt.Fprintf(w, "%s", string(b))
}
func randstring(length int) string {

	var fin []string
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ1234567890"
	chars := strings.Split(str, "")
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		fin = append(fin, chars[rand.Intn(26*2+10)])
	}
	return strings.Join(fin, "")
}
