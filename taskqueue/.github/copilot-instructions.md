Problem: Implementation in Go
Design and implement a robust and efficient task queue system in Go. This system should be capable of handling asynchronous tasks, allowing for the decoupling of task submission and execution.

Requirements: 

Task Definition
Define a ‘Task’ that represents a unit of work.
Each task should have a function or method that encapsulates the task’s execution logic.
Tasks might have associated metadata (e.g., priority, task ID, creation time).

Task Queue
Implement a queue data structure to store tasks.
The queue should support concurrent access for task submission and execution.

Task Submission
Provide a mechanism to submit tasks to the queue.

Task Execution
Implement a pool of worker goroutines to execute tasks from the queue.
Workers should continuously retrieve and execute tasks.
Handle potential errors during task execution gracefully.

Shutdown
Provide a mechanism to gracefully shut down the task queue, ensuring that all pending tasks are completed or handled appropriately.
