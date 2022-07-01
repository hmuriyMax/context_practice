package google

import (
	"context"
	"context_practice/pkg/search_ex/userip"
	"encoding/json"
	"log"
	"net/http"
)

type Result struct {
	Title string
	URL   string
}

type Results []Result

func Search(ctx context.Context, query string) (Results, error) {
	// Подготовливаем запрос API поиска Google.
	req, err := http.NewRequest("GET",
		"https://ajax.googleapis.com/ajax/services/search/web?v=1.0", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", query)

	// Если ctx передает IP-адрес пользователя,
	// перенаправляем его на сервер.
	// API Google используют
	// IP-адрес пользователя для различения запросов,
	// инициированных сервером
	// от запросов конечного пользователя.
	if userIP, ok := userip.FromContext(ctx); ok {
		q.Set("userip", userIP.String())
	}
	req.URL.RawQuery = q.Encode()

	var results Results
	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Обрабатываем JSON результат поиска.
		var data struct {
			ResponseData struct {
				Results []struct {
					TitleNoFormatting string
					URL               string
				}
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Println(err)
		}
		for _, res := range data.ResponseData.Results {
			results = append(results, Result{Title: res.TitleNoFormatting, URL: res.URL})
		}
		return nil
	})
	// httpDo ожидает возврата из предоставленного нами
	// замыкания, поэтому безопасно читать результаты здесь.
	return results, err
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Запускаем HTTP-запрос в goroutine
	// и передаем ответ в f.
	c := make(chan error, 1)
	req = req.WithContext(ctx)
	go func() { c <- f(http.DefaultClient.Do(req)) }()
	select {
	case <-ctx.Done():
		<-c // Ожидаем пока f вернется.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
