package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type Response struct {
	StatusMessage string
}

func Uploadhandle(w http.ResponseWriter, r *http.Request) {
	var message []byte
	if r.Method != "POST" {
		var resp Response
		resp.StatusMessage = "Inallowed Method"
		message, _ = json.Marshal(resp)
		fmt.Fprintf(w, string(message))
		return
	}
	token := r.FormValue("token")
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("partialFileDta")
	if err != nil {
		message, _ = json.Marshal(Response{"Error"})
		fmt.Fprintf(w, string(message))
		return
	}

	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	val, e := rdb.Get(ctx, token).Result()
	if e != nil {
		message, _ = json.Marshal(Response{"Error"})
		fmt.Fprintf(w, string(message))
		return
	}
	uploadedBytes, _ := ioutil.ReadAll(file)
	//Get real filename from redis!(var token)

	fileName := val
	rdb.Set(ctx, token, fileName, time.Minute*2)
	target, e := os.OpenFile("/storedblob/"+fileName, os.O_RDWR, 0777)
	if e != nil {
		message, _ = json.Marshal(Response{"Error"})
		fmt.Fprintf(w, string(message))
		return
	}
	startbyte, _ := strconv.Atoi(r.Header.Get("StartRange"))
	n, err := target.WriteAt(uploadedBytes, int64(startbyte))
	if n != len(uploadedBytes) {
		message, _ = json.Marshal(Response{"Error"})
		fmt.Fprintf(w, string(message))
		return
	}
	if err != nil {
		message, _ = json.Marshal(Response{"Error"})
		fmt.Fprintf(w, string(message))
		return
	}

	if target.Close(); e != nil {
		message, _ = json.Marshal(Response{"Error"})
		fmt.Fprintf(w, string(message))
	} else {
		message, _ = json.Marshal(Response{"Success"})
	}
	fmt.Fprintf(w, string(message))
}

//redis:token
