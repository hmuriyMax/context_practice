package main

import (
	"context_practice/pkg/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func WaitHandler(w http.ResponseWriter, r *http.Request) {
	waitTime, err := strconv.Atoi(r.URL.Query().Get("time"))
	//Есть возможность напрямую работать с контекстом через r.Context()
	res := types.JSONResponse{ConvSuccess: err == nil}
	if res.ConvSuccess {
		time.Sleep(time.Duration(waitTime) * time.Second)
		res.WaitedFor = waitTime
	}
	marshal, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Printf("%s\n", string(marshal))
	_, err = w.Write(marshal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/waiter", WaitHandler)
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
