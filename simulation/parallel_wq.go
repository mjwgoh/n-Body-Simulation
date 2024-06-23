package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"proj3-redesigned/utils"
	"proj3-redesigned/workstealing"
	"strconv"
	"sync"
	"time"
)

func worker(queue *workstealing.Dequeue, queues []*workstealing.Dequeue, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task, ok := queue.Pop()
		if !ok {
			// No task was available in the local queue, try to steal from others
			for _, q := range queues {
				if q != queue { // Avoid stealing from self
					task, ok = q.Steal()
					if ok {
						break // Successfully stole a task
					}
				}
			}
			if !ok { // No tasks are left to steal
				return // Exit the worker as there are no more tasks
			}
		}

		if task == nil { // Check for nil task to prevent panic
			continue
		}

		// Use a type switch to handle different task types and check for termination
		switch t := task.(type) {
		case *workstealing.NodeTask:
			t.Execute() // Execute NodeTask
		case *workstealing.BodyTask:
			t.Execute() // Execute BodyTask
		case *workstealing.TerminationTask:
			return // Terminate this worker
		default:
			// Optionally handle unknown task types or log an error
		}
	}
}

func simulateWQParallel(root *utils.QuadNode, allBodies *utils.Bodies, dt float64, theta float64, numWorkers int) {
	if root == nil {
		return
	}

	var wg sync.WaitGroup
	queues := make([]*workstealing.Dequeue, numWorkers)
	for i := range queues {
		queues[i] = workstealing.NewWorkStealingDequeue() // Initialize new lock-free queues
	}

	// Enqueue node tasks
	nodeList := []*utils.QuadNode{root}
	nodeIndex := 0

	for len(nodeList) > 0 {
		var nextNodes []*utils.QuadNode
		for _, node := range nodeList {
			for _, child := range node.Children {
				if child != nil {
					nextNodes = append(nextNodes, child)
				}
			}
			if node.BodiesPtr != nil && node.BodiesPtr.NodeBodies != nil {
				queueIndex := nodeIndex % numWorkers
				queues[queueIndex].Push(&workstealing.NodeTask{Node: node, Root: root, Theta: theta, Dt: dt})
				nodeIndex++
			}
		}
		nodeList = nextNodes
	}

	// Checks - print all nodes: Check passed - all tasks are appended correctly!
	//queues[0].PrintQueue()
	//queues[1].PrintQueue()

	for i := range queues {
		queues[i].Push(&workstealing.TerminationTask{})
	}

	// Start workers for node tasks
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(queues[i], queues, &wg)
	}

	wg.Wait() // Wait for all workers to finish processing nodes

	var bodyWg sync.WaitGroup
	bodyQueues := make([]*workstealing.Dequeue, numWorkers)
	for i := range bodyQueues {
		bodyQueues[i] = workstealing.NewWorkStealingDequeue()
	}

	// Enqueue body tasks
	bodyIndex := 0
	for _, body := range allBodies.NodeBodies {
		queueIndex := bodyIndex % numWorkers
		bodyQueues[queueIndex].Push(&workstealing.BodyTask{Body: body, Dt: dt})
		bodyIndex++
	}

	// Enqueue termination tasks for body tasks
	for i := 0; i < numWorkers; i++ {
		bodyQueues[i].Push(&workstealing.TerminationTask{})
	}

	// Start workers for body tasks
	for i := 0; i < numWorkers; i++ {
		bodyWg.Add(1)
		go worker(bodyQueues[i], bodyQueues, &bodyWg)
	}

	bodyWg.Wait() // Wait for all workers to finish processing bodies
}

func WQParallel(inputLink string, numWorkers int) {

	startTime := time.Now() // Start timing
	parallelTime := 0

	root, bodies, simulationFrames, _ := BuildQuadTree(inputLink)

	dt := 0.01 // Time step size

	// Create and open a CSV file
	file, err := os.Create("wq_parallel_simulation_results.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{"Frame", "Body Name", "PosX", "PosY", "VelX", "VelY", "ForceX", "ForceY"}
	writer.Write(headers)

	for frame := 0; frame < int(simulationFrames); frame++ {

		parallelStart := time.Now()
		simulateWQParallel(root, bodies, dt, 0.5, numWorkers)
		parallelEnd := time.Now()

		parallelTimeDelta := parallelEnd.Sub(parallelStart)
		parallelTime += int(parallelTimeDelta.Microseconds())

		root = RebuildQuadTree(bodies)

		// Write positions and velocities for each body to the CSV
		for _, body := range bodies.NodeBodies {
			record := []string{
				strconv.Itoa(frame),
				body.Name,
				fmt.Sprintf("%f", body.Positions.X),
				fmt.Sprintf("%f", body.Positions.Y),
				fmt.Sprintf("%f", body.Velocities.X),
				fmt.Sprintf("%f", body.Velocities.Y),
				fmt.Sprintf("%f", body.Force.X),
				fmt.Sprintf("%f", body.Force.Y),
			}
			if err := writer.Write(record); err != nil {
				fmt.Println("Error writing to CSV:", err)
			}
		}
		writer.Flush() // Flush after each frame to ensure data is written

	}

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	sequentialTime := int(totalTime.Microseconds()) - parallelTime

	fmt.Printf("Sequential %d, Parallel %d\n", sequentialTime, parallelTime)
}
