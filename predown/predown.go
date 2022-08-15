package predown

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

type NormalResponse struct {
	StatusMessage string
	Token         string
	Size          int
	Name          string
	Blobkey       string
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
		fmt.Fprintf(w, string(message))
		return
	}
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		message, _ = json.Marshal(Response{"Error"})
		fmt.Fprintf(w, string(message))
		return
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
		return
	}
	var uid string
	rows.Scan(&uid)
	targetFileId := r.FormValue("id")
	rows, err = db.Query("SELECT savedname, size, name, blobkey FROM files WHERE id=? AND userid=?", targetFileId, uid)
	if err != nil {
		message, _ = json.Marshal(Response{"Error"})
		fmt.Fprintf(w, string(message))
		return
	}
	defer rows.Close()
	if rows.Next() {
		var Size int
		var savedName, Name, Blobkey string
		rows.Scan(&savedName, &Size, &Name, &Blobkey)
		var ctx = context.Background()

		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		fileToken := randstring(20)
		e := rdb.Set(ctx, "view"+fileToken, savedName, time.Minute*3).Err()
		if e != nil {
			message, _ = json.Marshal(Response{"redisError"})
			fmt.Fprintf(w, string(message))
			return
		}
		message, _ = json.Marshal(NormalResponse{"Success", fileToken, Size, Name, Blobkey})
	}
	fmt.Fprintf(w, string(message))
}
