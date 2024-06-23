package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"proj3-redesigned/utils"
	"strconv"
	"time"
)

func simulate(root *utils.QuadNode, allBodies *utils.Bodies, dt float64, theta float64) {
	if root == nil {
		return
	}

	nodeList := []*utils.QuadNode{root}

	for len(nodeList) > 0 {
		nextNodes := []*utils.QuadNode{}

		for _, curNode := range nodeList {
			if curNode.Children[0] != nil && curNode.TotalMass > 0 {
				for _, childNode := range curNode.Children {
					if childNode != nil {
						nextNodes = append(nextNodes, childNode)
					}
				}
			}

			if curNode.BodiesPtr != nil && curNode.BodiesPtr.NodeBodies != nil {
				curNode.CalculateForce(root, theta, dt)
			}
		}

		nodeList = nextNodes
	}

	for _, body := range allBodies.NodeBodies {
		body.Update(dt)
	}

}

func Sequential(inputLink string) {

	sequentialStart := time.Now()

	root, bodies, simulationFrames, _ := BuildQuadTree(inputLink)

	dt := 0.01 // Time step size

	// Create and open a CSV file
	file, err := os.Create("sequential_simulation_results.csv")
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
		simulate(root, bodies, dt, 0.5)
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

	sequentialEnd := time.Now()
	sequentialTime := int(sequentialEnd.Sub(sequentialStart).Microseconds())

	fmt.Printf("Sequential %d, Parallel %d\n", sequentialTime, 0)

}
