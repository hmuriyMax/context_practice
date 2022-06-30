package main

import (
	"context"
	"context_practice/pkg/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func sendRequest(ctx context.Context, seconds int) {
	defer wg.Done()
	requestURL := fmt.Sprintf("http://localhost:8080/waiter?time=%d", seconds)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		log.Fatal(err)
	}
	client := http.DefaultClient
	//У клиента еще есть таймаут, при его истечении результат примерно такой же, как с контекстом
	//client.Timeout = 2 * time.Second
	var resp *http.Response

	resp, err = client.Do(r)

	if err != nil {
		log.Println(err)
		return
	}
	var result types.JSONResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Waited for %ds", result.WaitedFor)
}

func main() {
	ctx := context.Background()
	ctxWTO, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	wg.Add(1)
	go sendRequest(ctxWTO, 3)
	wg.Wait()
}
