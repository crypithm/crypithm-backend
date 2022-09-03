package preupload

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-redis/redis"
)

type Response struct {
	StatusMessage string
}

type SuccessResponse struct {
	StatusMessage string
	Rqid string
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

func Prehandle(w http.ResponseWriter, r *http.Request) {
	var message []byte
	var recievedVals [5]string
	if r.Method != "POST" {
		var resp Response
		resp.StatusMessage = "Inallowed Method"
		message, _ = json.Marshal(resp)
		fmt.Fprintf(w, string(message))
		return
	} else {
		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			message, _ = json.Marshal(Response{"NoAuthError"})
			fmt.Fprintf(w, string(message))
		}
		db, err := sql.Open("mysql", "crypithmusr:cDP9gNEQmUQt7qXbzU7XJ3Xz4mmcMf@tcp(127.0.0.1:3306)/crypithm")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query("SELECT uid FROM user WHERE token=?", token)
		defer rows.Close()
		if !rows.Next() {
			message, _ = json.Marshal(Response{"NoRowError"})
			fmt.Fprintf(w, string(message))
		}
		recievedVals[0] = r.FormValue("fileSize")
		recievedVals[1] = r.FormValue("fileName")
		recievedVals[2] = r.FormValue("chunkKey")
		recievedVals[3] = r.FormValue("id")
		recievedVals[4] = r.FormValue("dir")

		if !(len(recievedVals[3]) == 11 && (recievedVals[4] == "/ 0" || len(recievedVals[4]) == 15)) {
			message, _ = json.Marshal(Response{"Invalid Data Recieved"})
			fmt.Fprintf(w, string(message))
			return
		}
		for i := 0; i < len(recievedVals); i++ {
			if len(recievedVals[i]) == 0 {
				message, _ = json.Marshal(Response{"DataError"})
				fmt.Fprintf(w, string(message))
				break
			}
		}
		var uid string
		rows.Scan(&uid)

		fileName := randstring(16)

		_, e := db.Exec("INSERT INTO files (size, name,blobkey,id,directory,userid,savedname) values (?,?,?,?,?,?,?)", recievedVals[0], recievedVals[1], recievedVals[2], recievedVals[3], recievedVals[4], uid, fileName)
		if e != nil {
			message, _ = json.Marshal(Response{"DbError"})
			fmt.Fprintf(w, string(message))
			return
		}
		var ctx = context.Background()

		rdb := redis.NewClient(&redis.Options{
			Addr:     "140.238.219.8:6379",
			Password: "69GUaedM9MNApmU5wugCz5T7gdBa6Ka",
			DB:       0,
		})
		fileToken := randstring(20)
		e = rdb.Set(ctx, fileToken, fileName, time.Minute*3).Err()
		if e != nil {
			message, _ = json.Marshal(Response{"RdbError"})
			fmt.Fprintf(w, string(message))
			return
		}
		message, _ = json.Marshal(SuccessResponse{fileToken, "st-ch1"})
	}
	fmt.Fprintf(w, string(message))
}
