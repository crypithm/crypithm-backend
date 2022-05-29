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

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	StatusMessage string
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

func Predown(w http.ResponseWriter, r *http.Request) {
	var message []byte
	if r.Method != "POST" {
		var resp Response
		resp.StatusMessage = "Inallowed Method"
		message, _ = json.Marshal(resp)
	} else {
		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			message, _ = json.Marshal(Response{"Error"})
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
			message, _ = json.Marshal(Response{"Error"})
			fmt.Fprintf(w, string(message))
		} else {
			var uid string
			rows.Scan(&uid)
			targetFileId := r.FormValue("id")
			rows, err := db.Query("SELECT savedname FROM files WHERE id=? AND userid=?", targetFileId, uid)
			if err != nil {
				message, _ = json.Marshal(Response{"Error"})
				fmt.Fprintf(w, string(message))
			}
			defer rows.Close()
			if rows.Next() {
				var savedName string
				rows.Scan(&savedName)
				var ctx = context.Background()

				rdb := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Password: "",
					DB:       0,
				})
				fileToken := randstring(20)
				e := rdb.Set(ctx, fileToken, savedName, time.Minute*3).Err()
				if e != nil {
					message, _ = json.Marshal(Response{"redisError"})
					fmt.Fprintf(w, string(message))
				}
				message, _ = json.Marshal(Response{fileToken})
			}
		}

	}
	fmt.Fprintf(w, string(message))
}
