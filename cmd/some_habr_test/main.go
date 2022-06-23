package main

import (
	"context"
	"log"
	"time"
)

const (
	TCancel  = 4
	TContext = 2
	TWork    = 5
)

func sleepFor(fromFunction string, dur int, res chan bool) {
	defer func() {
		log.Printf("%s: sleepFor() complete\n", fromFunction)
	}()

	log.Printf("%s: started sleep for %ds\n", fromFunction, dur)
	time.Sleep(time.Duration(dur) * time.Second)
	log.Printf("%s: sleep finished arter %ds\n", fromFunction, dur)

	if res != nil {
		res <- true
	}
}

// Функция, выполняющая медленную работу с использованием контекста
// Заметьте, что контекст - это первый аргумент
func doWorkContext2(ctx context.Context, ch chan bool) {

	// Выполнение (прим. пер.: отложенное выполнение) действий по очистке
	// Созданных контекстов больше нет
	// Следовательно, отмена не требуется
	defer func() {
		log.Println("doWorkContext2() complete")
		ch <- true
	}()

	// Создаем канал
	sleepFinished := make(chan bool)

	// Запускаем выполнение медленной задачи в горутине
	// Передаем канал для коммуникаций
	go sleepFor("doWorkContext2", TWork, sleepFinished)

	// Используем select для выхода по истечении времени жизни контекста
	select {
	case <-ctx.Done():
		// Если контекст отменен, выбирается этот случай
		// Это случается, если заканчивается таймаут doWorkContext или
		// doWorkContext или main вызывает cancelFunction
		log.Println("doWorkContext2: cancel-func call detected")

	case <-sleepFinished:
		// Этот вариант выбирается, когда работа завершается до отмены контекста
		log.Println("doWorkContext2: sleep finished successfully")
	}
}

// Вспомогательная функция, которая в реальности может использоваться для разных целей
// Здесь она просто вызывает одну функцию
func doWorkContext1(ctx context.Context) {
	// От контекста с функцией отмены создаём производный контекст с таймаутом 2 секунды
	ctxWithTimeout, cancelFunction := context.WithTimeout(ctx, TContext*time.Second)

	// Функция отмены для освобождения ресурсов после завершения функции
	defer func() {
		log.Println("doWorkContext1() complete")
		cancelFunction()
	}()

	// Создаем канал и вызываем функцию контекста
	// Можно также использовать группы ожидания для этого конкретного случая,
	// поскольку мы не используем возвращаемое значение, отправленное в канал
	ch := make(chan bool)
	go doWorkContext2(ctxWithTimeout, ch)

	// Используем select для выхода при истечении контекста
	select {
	case <-ctx.Done():
		// Этот случай выбирается, когда переданный в качестве аргумента контекст уведомляет о завершении работы
		// В данном примере это произойдёт, когда в main будет вызвана cancelFunction
		log.Println("doWorkContext1: cancel-func call detected")

	case <-ch:
		// Этот вариант выбирается, когда работа завершается до отмены контекста
		log.Println("doWorkContext1: dWC2 returned successfully")
	}
}

func main() {
	// Создаем контекст background
	ctx := context.Background()
	// Производим контекст с отменой
	ctxWithCancel, cancelFunction := context.WithCancel(ctx)

	// Отложенная функция вызывает функцию отмены
	defer func() {
		log.Println("Main Defer: canceling context")
		cancelFunction()
	}()

	// Отмена контекста после TCancel секунд.
	// Если это происходит, все производные от него контексты должны завершиться
	go func() {
		sleepFor("Main", TCancel, nil)
		cancelFunction()
	}()

	// Выполнение работы
	doWorkContext1(ctxWithCancel)
}
