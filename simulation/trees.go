package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"proj3-redesigned/quadtree"
	"proj3-redesigned/utils"
)

func BuildQuadTree(inputLink string) (*utils.QuadNode, *utils.Bodies, float64, float64) {

	// Read Input
	cwd, _ := os.Getwd()
	fileName := fmt.Sprintf("simulation/data/%s.csv", inputLink)
	dataDir := filepath.Join(cwd, fileName)
	bodies, SimulationFrames, GravitationalConstant := utils.ReadInput(dataDir)

	var minimums, maximums utils.Vector2

	minimums.X, minimums.Y = math.MaxFloat64, math.MaxFloat64
	maximums.X, maximums.Y = -math.MaxFloat64, -math.MaxFloat64

	for _, body := range bodies.NodeBodies {

		minimums.X = math.Min(minimums.X, body.Positions.X)
		maximums.X = math.Max(maximums.X, body.Positions.X)

		minimums.Y = math.Min(minimums.Y, body.Positions.Y)
		maximums.Y = math.Max(maximums.Y, body.Positions.Y)

	}

	// Starting Region
	startRegion := [2]utils.Vector2{{minimums.X - 1, minimums.Y - 1},
		{maximums.X + 1, maximums.Y + 1},
	} // Example bounds
	root := quadtree.BuildQuadTree(bodies.NodeBodies, startRegion)

	return root, &bodies, SimulationFrames, GravitationalConstant

}

func RebuildQuadTree(bodies *utils.Bodies) *utils.QuadNode {

	// Read Input
	var minimums, maximums utils.Vector2

	minimums.X, minimums.Y = math.MaxFloat64, math.MaxFloat64
	maximums.X, maximums.Y = -math.MaxFloat64, -math.MaxFloat64

	for _, body := range bodies.NodeBodies {

		minimums.X = math.Min(minimums.X, body.Positions.X)
		maximums.X = math.Max(maximums.X, body.Positions.X)

		minimums.Y = math.Min(minimums.Y, body.Positions.Y)
		maximums.Y = math.Max(maximums.Y, body.Positions.Y)

	}

	// Starting Region
	startRegion := [2]utils.Vector2{{minimums.X - 1, minimums.Y - 1},
		{maximums.X + 1, maximums.Y + 1},
	} // Example bounds
	root := quadtree.BuildQuadTree(bodies.NodeBodies, startRegion)
	//fmt.Println("Successfully built Quadtree")

	root.TotalMass = 1

	return root

}
