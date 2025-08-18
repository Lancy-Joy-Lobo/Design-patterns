package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	TaskId       int
	CreationTime time.Time
	Priority     int
	Execute      func() error
}

type TaskQueue struct {
	Tasks   chan Task
	wg      sync.WaitGroup
	Workers int
}

func NewTaskQueue(totalTasks, workers int) *TaskQueue {
	tq := &TaskQueue{
		Tasks:   make(chan Task, totalTasks),
		wg:      sync.WaitGroup{},
		Workers: workers,
	}

	for i := 0; i < workers; i++ {
		tq.wg.Add(1)
		go tq.ExecuteTasks(i + 1)
	}

	return tq
}

func (tq *TaskQueue) Submit(t Task) {
	tq.Tasks <- t
}

func (tq *TaskQueue) ExecuteTasks(workerId int) {
	defer tq.wg.Done()
	fmt.Printf("Worker %d started\n", workerId)
	for t := range tq.Tasks {
		err := t.Execute()
		if err != nil {
			fmt.Printf("Worker %d failed to execute task %d: %v\n", workerId, t.TaskId, err)
		} else {
			fmt.Printf("Worker %d executed task %d\n", workerId, t.TaskId)
		}
	}
	fmt.Printf("Worker %d shutting down\n", workerId)
}

func (tq *TaskQueue) ShutDown() {
	fmt.Println("Shutting down task queue...")
	close(tq.Tasks)
	tq.wg.Wait()
	fmt.Println("All workers have shut down.")
}

func main() {

	totalTasks := 20
	workers := 3

	tq := NewTaskQueue(totalTasks, workers)

	// Use a separate WaitGroup for task submission to ensure all tasks are submitted
	// before we consider shutting down. This is good practice if the main function
	// could potentially exit before submission is complete.
	var submissionWg sync.WaitGroup
	for i := 1; i <= totalTasks; i++ {
		submissionWg.Add(1)
		go func(taskId int) {
			defer submissionWg.Done()
			task := Task{
				TaskId:       taskId,
				CreationTime: time.Now(),
				Priority:     taskId,
				Execute: func() error {
					// Simulate some work
					if taskId%5 == 0 {
						return errors.New("simulated error")
					}
					time.Sleep(1 * time.Second)
					return nil
				},
			}
			fmt.Printf("Submitting task %d\n", taskId)
			tq.Submit(task)
		}(i)
	}

	// Wait for all tasks to be submitted before shutting down
	submissionWg.Wait()
	fmt.Println("All tasks submitted.")

	// Now, shut down the queue. This will wait for all workers to finish.
	tq.ShutDown()

	fmt.Println("Main application finished.")
}
