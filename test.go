package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Actor struct {
	mailBox  chan any
	counter  int
	response chan any
}

func New() *Actor {

	act := &Actor{
		mailBox:  make(chan any),
		counter:  0,
		response: make(chan any),
	}

	go act.start()

	return act
}

func (a *Actor) send(message any) {
	a.mailBox <- message
}

func (a *Actor) start() {
	for {
		value := <-a.mailBox

		switch value {
		case "increment":
			a.counter++
		case "decrement":
			a.counter--
		case "get":
			a.response <- a.counter
		}
	}
}

type ThreadPool struct {
	work  chan func()
	group *sync.WaitGroup
	quit  chan any
}

func NewPool(workers int) *ThreadPool {

	tp := ThreadPool{
		work:  make(chan func()),
		group: &sync.WaitGroup{},
		quit:  make(chan any),
	}

	for i := 0; i < workers; i++ {
		go tp.start()
	}

	return &tp
}

func (r *ThreadPool) shutdown() {
	r.group.Wait()
}

func (r *ThreadPool) shutdownNow() {
	r.quit <- 1
}

func (r *ThreadPool) submit(action func()) {
	r.group.Add(1)
	r.work <- action
}

func (r *ThreadPool) start() {
	for {
		select {
		case value := <-r.work:
			value()
			r.group.Done()
		case <-r.quit:
			break
		}
	}
}

func mains() {

	tp := NewPool(10)

	counterActor := New()

	for i := 0; i < 100_000; i++ {
		tp.submit(func() {
			counterActor.send("increment")
		})
	}

	tp.shutdown()

	counterActor.mailBox <- "get"

	result := <-counterActor.response

	fmt.Println(result)

}

func test() {
	results, errorsList := WaitAll(func() (any, error) {
		log.Println("finish 1")
		return "a", nil
	}, func() (any, error) {
		log.Println("finish 2")
		return "a", nil
	}, func() (any, error) {
		log.Println("finish 3")
		return "a", nil
	}, func() (any, error) {
		time.Sleep(1 * time.Second)
		log.Println("finish 4")
		return nil, errors.New("erow")
	}, func() (any, error) {
		time.Sleep(1 * time.Second)
		log.Println("finish 5")
		return nil, errors.New("erowwwww")
	})

	fmt.Println(results)
	fmt.Println(errorsList)
}

func WaitAll(actions ...func() (any, error)) ([]any, []error) {

	var wg sync.WaitGroup

	length := len(actions)

	resultChannel := make(chan any, length)

	errorChannel := make(chan error, 1)

	for _, action := range actions {
		wg.Add(1)

		go func(action func() (any, error)) {
			defer wg.Done()

			var result, err = action()

			if err != nil {
				errorChannel <- err
				return
			}

			resultChannel <- result
		}(action)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
		close(errorChannel)
	}()

	var errorList []error
	for err := range errorChannel {
		errorList = append(errorList, err)
	}

	if len(errorList) > 0 {
		return nil, errorList
	}

	var results []any
	for result := range resultChannel {
		results = append(results, result)
	}

	return results, nil
}
