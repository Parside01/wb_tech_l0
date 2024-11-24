package workerpool

import "time"

type Job struct {
	ID   int
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type JobChannel chan Job
type JobQueue chan chan Job

type Worker struct {
	ID      int
	JobChan JobChannel
	Queue   JobQueue
	Break   chan struct{}
}

func New(ID int, jobChan JobChannel, queue JobQueue, breakChan chan struct{}) *Worker {
	return &Worker{
		ID:      ID,
		JobChan: jobChan,
		Queue:   queue,
		Break:   breakChan,
	}
}

func (w *Worker) Start() {

}
