package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type Response struct {
	StatusMessage string
}

func Downloader(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		token := r.FormValue("token")

		var ctx = context.Background()

		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		val, e := rdb.Get(ctx, "view"+token).Result()
		if e != nil {
			return
		}
		fileName := val
		file, err := os.Open("/storedblob/" + fileName)
		if err != nil {
			return
		}
		fileStat, e := file.Stat()
		var realBytes []byte
		startbyte, e1 := strconv.Atoi(r.Header.Get("StartRange"))
		endbyte, e2 := strconv.Atoi(r.Header.Get("EndRange"))
		if e1 != nil || e2 != nil {
			return
		}
		if int64(endbyte) > fileStat.Size() {
			if int64(startbyte) > fileStat.Size() {
				return
			} else {
				realBytes = make([]byte, fileStat.Size()-int64(startbyte))
			}
		} else {
			realBytes = make([]byte, endbyte-startbyte)
		}
		defer file.Close()
		if startbyte == 0 {
			file.Seek(0, 0)
		} else if startbyte > 0 {
			file.Seek(int64(startbyte), 0)
		}
		_, e = io.ReadFull(file, realBytes)
		rdb.Set(ctx, "view"+token, fileName, time.Minute*2)

		fmt.Fprintf(w, "%s", realBytes)
	}
}
