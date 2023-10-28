package utils

import "sync"

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
