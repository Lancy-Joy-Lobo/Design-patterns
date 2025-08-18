package main

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Task represents a unit of work with priority and metadata
type Task struct {
	ID           int
	Priority     int
	CreationTime time.Time
	Execute      func() error
}

// PriorityQueue implements a priority queue for tasks
type PriorityQueue struct {
	tasks []*Task
	mu    sync.Mutex
}

// NewPriorityQueue creates a new priority queue
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		tasks: make([]*Task, 0),
	}
	heap.Init(pq)
	return pq
}

// Len returns the length of the queue (heap.Interface)
func (pq *PriorityQueue) Len() int {
	return len(pq.tasks)
}

// Less compares tasks by priority (heap.Interface)
func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.tasks[i].Priority > pq.tasks[j].Priority // Higher priority first
}

// Swap swaps two tasks (heap.Interface)
func (pq *PriorityQueue) Swap(i, j int) {
	pq.tasks[i], pq.tasks[j] = pq.tasks[j], pq.tasks[i]
}

// Push adds a task to the queue (heap.Interface)
func (pq *PriorityQueue) Push(x interface{}) {
	pq.tasks = append(pq.tasks, x.(*Task))
}

// Pop removes and returns the highest priority task (heap.Interface)
func (pq *PriorityQueue) Pop() interface{} {
	n := len(pq.tasks)
	if n == 0 {
		return nil
	}
	task := pq.tasks[n-1]
	pq.tasks = pq.tasks[:n-1]
	return task
}

// SafePush safely adds a task to the priority queue
func (pq *PriorityQueue) SafePush(task *Task) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(pq, task)
}

// SafePop safely removes and returns the highest priority task
func (pq *PriorityQueue) SafePop() *Task {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	
	if pq.Len() == 0 {
		return nil
	}
	
	task := heap.Pop(pq)
	if task == nil {
		return nil
	}
	return task.(*Task)
}

// TaskQueue manages task execution with worker goroutines
type TaskQueue struct {
	queue          *PriorityQueue
	workers        int
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	shutdown       chan struct{}
	tasksSubmitted int64
	tasksCompleted int64
	tasksFailed    int64
}

// NewTaskQueue creates a new task queue with specified number of workers
func NewTaskQueue(workers int) *TaskQueue {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskQueue{
		queue:    NewPriorityQueue(),
		workers:  workers,
		ctx:      ctx,
		cancel:   cancel,
		shutdown: make(chan struct{}),
	}
}

// Start initializes and starts worker goroutines
func (tq *TaskQueue) Start() {
	for i := 0; i < tq.workers; i++ {
		tq.wg.Add(1)
		go tq.worker(i + 1)
	}
}

// SubmitTask adds a new task to the queue
func (tq *TaskQueue) SubmitTask(task *Task) error {
	select {
	case <-tq.ctx.Done():
		return errors.New("task queue is shutting down")
	default:
		tq.queue.SafePush(task)
		atomic.AddInt64(&tq.tasksSubmitted, 1)
		return nil
	}
}

// worker is the main worker goroutine that processes tasks
func (tq *TaskQueue) worker(id int) {
	defer tq.wg.Done()
	
	log.Printf("Worker %d started\n", id)
	
	for {
		select {
		case <-tq.ctx.Done():
			log.Printf("Worker %d shutting down\n", id)
			return
		default:
			task := tq.queue.SafePop()
			if task == nil {
				// No tasks available, wait a bit before checking again
				time.Sleep(100 * time.Millisecond)
				continue
			}
			
			tq.executeTask(task, id)
		}
	}
}

// executeTask executes a single task and handles errors
func (tq *TaskQueue) executeTask(task *Task, workerID int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Worker %d: Task %d panicked: %v\n", workerID, task.ID, r)
			atomic.AddInt64(&tq.tasksFailed, 1)
		}
	}()
	
	log.Printf("Worker %d executing task %d (priority: %d)\n", workerID, task.ID, task.Priority)
	
	if err := task.Execute(); err != nil {
		log.Printf("Worker %d: Task %d failed: %v\n", workerID, task.ID, err)
		atomic.AddInt64(&tq.tasksFailed, 1)
	} else {
		log.Printf("Worker %d: Task %d completed successfully\n", workerID, task.ID)
		atomic.AddInt64(&tq.tasksCompleted, 1)
	}
}

// Shutdown gracefully shuts down the task queue
func (tq *TaskQueue) Shutdown(timeout time.Duration) error {
	log.Println("Initiating task queue shutdown...")
	
	// Cancel context to signal workers to stop
	tq.cancel()
	
	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		tq.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		log.Println("All workers shut down gracefully")
		close(tq.shutdown)
		return nil
	case <-time.After(timeout):
		log.Println("Shutdown timeout exceeded, forcing shutdown")
		close(tq.shutdown)
		return errors.New("shutdown timeout exceeded")
	}
}

// Stats returns current statistics about the task queue
func (tq *TaskQueue) Stats() (submitted, completed, failed int64) {
	return atomic.LoadInt64(&tq.tasksSubmitted),
		   atomic.LoadInt64(&tq.tasksCompleted),
		   atomic.LoadInt64(&tq.tasksFailed)
}

// WaitForShutdown blocks until the task queue has been shut down
func (tq *TaskQueue) WaitForShutdown() {
	<-tq.shutdown
}

// Example usage and demonstration
func main() {
	const (
		numWorkers = 3
		numTasks   = 20
	)
	
	// Create and start task queue
	tq := NewTaskQueue(numWorkers)
	tq.Start()
	
	// Submit tasks with different priorities
	for i := 0; i < numTasks; i++ {
		taskID := i
		task := &Task{
			ID:           taskID,
			Priority:     getPriority(taskID),
			CreationTime: time.Now(),
			Execute: func() error {
				// Simulate work
				workDuration := time.Duration(100+taskID*10) * time.Millisecond
				time.Sleep(workDuration)
				
				// Simulate occasional failures
				if taskID%7 == 0 {
					return fmt.Errorf("simulated error for task %d", taskID)
				}
				
				return nil
			},
		}
		
		if err := tq.SubmitTask(task); err != nil {
			log.Printf("Failed to submit task %d: %v\n", taskID, err)
		}
	}
	
	log.Printf("Submitted %d tasks\n", numTasks)
	
	// Let tasks run for a while
	time.Sleep(5 * time.Second)
	
	// Print stats before shutdown
	submitted, completed, failed := tq.Stats()
	log.Printf("Stats - Submitted: %d, Completed: %d, Failed: %d, Pending: %d\n",
		submitted, completed, failed, submitted-completed-failed)
	
	// Graceful shutdown with timeout
	if err := tq.Shutdown(10 * time.Second); err != nil {
		log.Printf("Shutdown error: %v\n", err)
	}
	
	// Final stats
	submitted, completed, failed = tq.Stats()
	log.Printf("Final Stats - Submitted: %d, Completed: %d, Failed: %d\n",
		submitted, completed, failed)
}

// getPriority assigns priority based on task ID
func getPriority(taskID int) int {
	switch {
	case taskID%5 == 0:
		return 10 // Highest priority
	case taskID%3 == 0:
		return 7  // High priority
	case taskID%2 == 0:
		return 5  // Medium priority
	default:
		return 2  // Low priority
	}
}