package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"proj3-redesigned/utils"
	"strconv"
	"sync"
	"time"
)

func workerCalculateForce(nodeQueue chan *utils.QuadNode, root *utils.QuadNode, theta float64, dt float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for node := range nodeQueue {
		if node.BodiesPtr != nil && node.BodiesPtr.NodeBodies != nil {
			node.CalculateForce(root, theta, dt)
		}
	}
}

func workerUpdateBodies(bodyQueue chan *utils.Body, dt float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for body := range bodyQueue {
		body.Update(dt)
	}
}

func simulateParallel(root *utils.QuadNode, allBodies *utils.Bodies, dt float64, theta float64, numWorkers int) {
	if root == nil {
		return
	}

	var wg sync.WaitGroup
	nodeQueue := make(chan *utils.QuadNode, 100)
	bodyQueue := make(chan *utils.Body, 100)

	// Start workers for force calculation
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go workerCalculateForce(nodeQueue, root, theta, dt, &wg)
	}

	// Enqueue nodes
	nodeList := []*utils.QuadNode{root}
	for len(nodeList) > 0 {
		var nextNodes []*utils.QuadNode
		for _, node := range nodeList {
			for _, child := range node.Children {
				if child != nil { // Add this check to ensure no nil values are enqueued
					nextNodes = append(nextNodes, child)
				}
			}
			if node.BodiesPtr != nil && node.BodiesPtr.NodeBodies != nil {
				nodeQueue <- node
			}
		}
		nodeList = nextNodes
	}
	close(nodeQueue) // Close the node queue after all nodes are enqueued
	wg.Wait()        // Wait for all force calculations to complete

	// Start workers for updating bodies
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go workerUpdateBodies(bodyQueue, dt, &wg)
	}

	// Enqueue bodies
	for _, body := range allBodies.NodeBodies {
		bodyQueue <- body
	}

	close(bodyQueue) // Close the body queue after all bodies are enqueued

	wg.Wait() // Wait for all updates to complete
}

func Parallel(inputLink string, numWorkers int) {

	startTime := time.Now() // Start timing
	parallelTime := 0

	root, bodies, simulationFrames, _ := BuildQuadTree(inputLink)

	dt := 0.01 // Time step size

	// Create and open a CSV file
	file, err := os.Create("parallel_simulation_results.csv")
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

		simulateParallel(root, bodies, dt, 0.5, numWorkers)

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
