package share

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-redis/redis"
)

type Defaultresp struct {
	Message string
}
type Linkedresp struct {
	Message string
	Link    string
}
type Keyresp struct {
	Message string
	Key     string
	Token   string
	Id      string
}

func randstring(length int) string {

	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

func ShareHandler(w http.ResponseWriter, r *http.Request) {
	if !(r.Method == "POST" || r.Method == "PUT") {
		message, _ := json.Marshal(Defaultresp{"Error"})
		fmt.Fprintf(w, string(message))
		return
	}
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
		return
	}
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	var uid string
	rows.Scan(&uid)
	var link string
	if r.Method == "POST" {
		originId := r.FormValue("id")
		err := db.QueryRow("SELECT encLink from shared WHERE uid=? AND originId=?", uid, originId).Scan(&link)
		if err == sql.ErrNoRows {
			var key string
			db.QueryRow("SELECT blobkey from files WHERE userid=? AND id=?", uid, originId).Scan(&key)
			shareToken := randstring(20)
			Id := randstring(10)
			e := rdb.Set(ctx, shareToken, Id+","+originId, time.Minute*3).Err()
			if e != nil {
				message, _ := json.Marshal(Defaultresp{"Error"})
				fmt.Fprintf(w, string(message))
				return
			}
			message, _ := json.Marshal(Keyresp{"Create", key, shareToken, Id})
			fmt.Fprintf(w, string(message))
			return
		} else if err == nil {
			message, _ := json.Marshal(Linkedresp{"Success", link})
			fmt.Fprintf(w, string(message))
			return
		} else {
			message, _ := json.Marshal(Defaultresp{"Error"})
			fmt.Fprintf(w, string(message))
			return
		}
	} else {
		var DisplayUser bool
		if r.FormValue("du") == "true" {
			DisplayUser = true
		} else {
			DisplayUser = false
		}
		Link := r.FormValue("link")
		Name := r.FormValue("name")
		shareId := r.FormValue("shareId")
		Token := r.FormValue("token")
		Key := r.FormValue("Key")
		val, e := rdb.Get(ctx, Token).Result()
		RedisId := strings.Split(val, ",")[0]
		if RedisId == shareId {
			originId := strings.Split(val, ",")[1]
			var Size string
			db.QueryRow("SELECT size from files WHERE userid=? AND id=?", uid, originId).Scan(&Size)

			_, e = db.Exec("INSERT INTO shared (mappedId, originId, fileName, fileSize, sharedKey, uid, showName, encLink) values (?,?,?,?,?,?,?,?)", shareId, originId, Name, Size, Key, uid, DisplayUser, Link)
			if e != nil {
				message, _ := json.Marshal(Defaultresp{"Error"})
				fmt.Fprintf(w, string(message))
				return
			}
			message, _ := json.Marshal(Defaultresp{"Success"})
			fmt.Fprintf(w, string(message))
			return
		} else {
			message, _ := json.Marshal(Defaultresp{"Error"})
			fmt.Fprintf(w, string(message))
			return
		}
	}

}
