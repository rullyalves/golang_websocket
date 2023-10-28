package queue

import "time"

func StartWorker(processQueue ProcessMessage, processMessage MessageProcessor) func(queueName string) {
	return func(queueName string) {
		for {
			err := processQueue(queueName, 10, processMessage)
			if err != nil {
				time.Sleep(2 * time.Second)
			}
		}
	}
}
