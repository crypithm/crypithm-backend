package viewShared

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Defaultresp struct {
	Message string
}

type Namedresp struct {
	Username string
	Name     string
	Size     string
	Key      string
	Token    string
	Rqid string
}

type Noname struct {
	Name  string
	Size  string
	Key   string
	Token string
	Rqid string
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

func SharedHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "140.238.219.8:6379",
			Password: "69GUaedM9MNApmU5wugCz5T7gdBa6K",
			DB:       0,
		})
		id := r.FormValue("id")
		db, err := sql.Open("mysql", "crypithmusr:cDP9gNEQmUQt7qXbzU7XJ3Xz4mmcMf@tcp(127.0.0.1:3306)/crypithm")
		if err != nil {
			message, _ := json.Marshal(Defaultresp{"ConnError"})
			fmt.Fprintf(w, string(message))
			return
		}

		var originId, fileName, fileSize, sharedKey, uid string
		var showName bool
		err = db.QueryRow("SELECT originId, fileName, fileSize, sharedKey,uid, showName FROM shared WHERE mappedId=?", id).Scan(&originId, &fileName, &fileSize, &sharedKey, &uid, &showName)
		if err != nil {

			message, _ := json.Marshal(Defaultresp{"QueryError"})
			fmt.Fprintf(w, string(message))
			return
		}

		var savedName string
		err = db.QueryRow("SELECT savedname FROM files WHERE userid=? AND id=?", uid, originId).Scan(&savedName)

		token := randstring(20)

		var ctx = context.Background()
		e := rdb.Set(ctx, "view"+token, savedName, time.Minute*1).Err()
		if e != nil {
			message, _ := json.Marshal(Defaultresp{"RedisError"})
			fmt.Fprintf(w, string(message))
			return
		}
		if showName == true {
			var username string
			err = db.QueryRow("SELECT username from user WHERE uid=?", uid).Scan(&username)
			message, _ := json.Marshal(Namedresp{username, fileName, fileSize, sharedKey, token, "st-ch1"})
			fmt.Fprintf(w, string(message))
			return
		} else {
			message, _ := json.Marshal(Noname{fileName, fileSize, sharedKey, token, "st-ch1"})
			fmt.Fprintf(w, string(message))
			return
		}
	} else {
		message, _ := json.Marshal(Defaultresp{"Invalid Method"})
		fmt.Fprintf(w, string(message))
		return
	}
}
