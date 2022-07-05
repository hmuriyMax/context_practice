package main

import (
	"fmt"
	"log"
	"sync"
)

type Muflex struct {
	mutex sync.Mutex
}

func (m *Muflex) Mutexer(f func() error) error {
	if &m.mutex == nil {
		m.mutex = sync.Mutex{}
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return f()
}

func main() {
	wg := sync.WaitGroup{}
	mu := Muflex{}
	sum := 0
	n := 1000
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(el int) {
			err := mu.Mutexer(func() error {
				sum += el
				return nil
			})
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Sum is %d\n", sum)
}
