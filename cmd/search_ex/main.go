package main

import (
	"context"
	"context_practice/pkg/search_ex/google"
	"context_practice/pkg/search_ex/userip"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func handleSearch(w http.ResponseWriter, req *http.Request) {
	// ctx - это Context для этого обработчика.
	// Вызов отмены закрывает
	// канал ctx.Done,
	// который является сигналом отмены для запросов
	// запущенных этим обработчиком.
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	timeout, err := time.ParseDuration(req.FormValue("timeout"))
	if err == nil {
		// У запроса есть timeout,
		// поэтому создаем контекст, который
		// автоматически отменяется по истечении
		// времени ожидания.
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	// Отмена ctx, как только вернется handleSearch.
	defer cancel()

	// Проверяем поисковый запрос.
	query := req.FormValue("q")
	if query == "" {
		http.Error(w, "no query", http.StatusBadRequest)
		return
	}

	// Сохраняем IP-адрес пользователя в ctx
	// для использования кодом в других пакетах.
	userIP, err := userip.FromRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx = userip.NewContext(ctx, userIP)

	// Запустить поиск Google и распечатать результаты.
	start := time.Now()
	results, err := google.Search(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	elapsed := time.Since(start)

	if err := Execute(w, struct {
		Results          google.Results
		Timeout, Elapsed time.Duration
	}{
		Results: results,
		Timeout: timeout,
		Elapsed: elapsed,
	}); err != nil {
		log.Print(err)
		return
	}
}

func Execute(w http.ResponseWriter, s struct {
	Results          google.Results
	Timeout, Elapsed time.Duration
}) error {
	marshal, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = w.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", handleSearch)
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
