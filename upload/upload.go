package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	} else {

		token := r.FormValue("token")

		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("partialFileDta")
		if err != nil {
			message, _ = json.Marshal(Response{"Error"})
			fmt.Fprintf(w, string(message))
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
		}

		uploadedBytes, _ := ioutil.ReadAll(file)
		//Get real filename from redis!(var token)

		fileName := val
		target, e := os.OpenFile("/storedblob/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if e != nil {
			message, _ = json.Marshal(Response{"Error"})
			fmt.Fprintf(w, string(message))
		}
		target.Write(uploadedBytes)
		if target.Close(); e != nil {
			message, _ = json.Marshal(Response{"Error"})
			fmt.Fprintf(w, string(message))
		} else {
			message, _ = json.Marshal(Response{"Success"})
		}
	}
	fmt.Fprintf(w, string(message))
}

//redis:token
