package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

//Функция, которая запишет во второй параметр true, когда контекст закроется
func checkCtx(ctx context.Context, res *bool) {
	*res = false
	select {
	case <-ctx.Done():
		*res = true
	}
}

//Какой-то процесс, выводящий каждую секунду сообщение.
func process(ctx context.Context) {
	ctxClosed := false
	go checkCtx(ctx, &ctxClosed)

	//Выводим сообщение о том, как завершилась функция: сама или при закрытии контекста
	defer func() {
		if ctxClosed {
			fmt.Printf("Context cancelled: finishing")
		} else {
			fmt.Printf("Planned finish!")
		}
		wg.Done()
	}()

	//Сама работа функции
	for i := 0; i < 8; i++ {
		fmt.Printf("%ds passed\n", i)
		time.Sleep(time.Second)
		if ctxClosed {
			return
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancelFunction := context.WithCancel(ctx)

	wg.Add(1)
	//Запускаем процесс в горутине
	go process(ctx)
	//Ожидаем ввода пользователя
	_, err := fmt.Fscanln(os.Stdin)
	if err != nil {
		log.Fatalf(err.Error())
	}
	//После ввода вызываем функцию отмены. Таким образом, как только пользователь нажал enter, контекст отменяется
	cancelFunction()
	//Ожидаем завершения вейт-группы
	wg.Wait()
}
