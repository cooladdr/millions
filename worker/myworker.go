package worker

import (
	"fmt"
)

type Worker struct {
	WorkerPool chan Worker
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan Worker) Worker {
	return Worker{WorkerPool: workerPool, JobChannel: make(chan Job), quit: make(chan bool)}
}

func (w Worker) Start() {
	go func() {
		fmt.Println("start worker: ")
		for {
			w.WorkerPool <- w
			select {
			case job := <-w.JobChannel:
				if err := job.Payload.UploadToS3(); err != nil {
					fmt.Errorf("Error uploading to S3: %s", err.Error())
				}

			case <-w.quit:
				fmt.Println(" worker exit: ")
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
