package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

const G = 6.67430e-11 // Gravitational constant

type Body struct {
	Name       string
	Positions  Vector2 // [x, y]
	Velocities Vector2
	Mass       float64
	Force      Vector2
}

type Bodies struct {
	NodeBodies []*Body
}

type QuadNode struct {
	Center    Vector2
	TotalMass float64
	Region    [2]Vector2   // min and max values for all 3 dimensions (x, y, and z).
	Children  [4]*QuadNode // up to 4 octonode children per node. Consider a 2x2x2 cube.
	BodiesPtr *Bodies
}

func (node *QuadNode) IsLeaf() bool {
	for _, child := range node.Children {
		if child != nil {
			return false
		}
	}
	return true
}

// Fixed function to read bodies from CSV
func ReadInput(filename string) (Bodies, float64, float64) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Correct delimiter for CSV files
	reader.TrimLeadingSpace = true

	var bodies Bodies
	var SimulationTimeInSeconds float64
	var GravitationalConstant float64

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // Exit the loop at end of file
		}
		if err != nil {
			fmt.Printf("Error reading input: %s\n", err)
			continue
		}

		if record[0] == "SimulationTime" {
			if len(record) >= 4 { // Check for at least 4 fields
				if simTime, err := strconv.ParseFloat(record[1], 64); err == nil {
					SimulationTimeInSeconds = simTime
				}
				if gravConst, err := strconv.ParseFloat(record[3], 64); err == nil {
					GravitationalConstant = gravConst
				}
			}
			continue
		} else {
			var body Body
			body.Name = record[0]
			if xPosition, err := strconv.ParseFloat(record[1], 32); err == nil {
				body.Positions.X = float64(xPosition)
			}
			if yPosition, err := strconv.ParseFloat(record[2], 32); err == nil {
				body.Positions.Y = float64(yPosition)
			}
			if xVelocity, err := strconv.ParseFloat(record[3], 32); err == nil {
				body.Velocities.X = float64(xVelocity)
			}
			if yVelocity, err := strconv.ParseFloat(record[4], 32); err == nil {
				body.Velocities.Y = float64(yVelocity)
			}
			if mass, err := strconv.ParseFloat(record[5], 32); err == nil {
				body.Mass = float64(mass)
			}
			bodies.NodeBodies = append(bodies.NodeBodies, &body)
		}
	}
	return bodies, SimulationTimeInSeconds, GravitationalConstant
}

type Vector2 struct {
	X, Y float64
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{v.X + other.X, v.Y + other.Y}
}

func (v Vector2) Subtract(other Vector2) Vector2 {
	return Vector2{v.X - other.X, v.Y - other.Y}
}

func (v Vector2) Dot(other Vector2) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vector2) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize returns a unit vector in the direction of v.
func (v Vector2) Normalize() Vector2 {
	mag := v.Magnitude()
	return Vector2{v.X / mag, v.Y / mag}
}

// Multiply returns the vector multiplied by a scalar.
func (v Vector2) Multiply(scalar float64) Vector2 {
	return Vector2{v.X * scalar, v.Y * scalar}
}

func updateForce(curNode *QuadNode, node *QuadNode, theta float64, dt float64) {

	if node == nil { // base case
		return
	}

	if curNode == node {
		return
	}

	for _, child := range node.Children {
		if child != nil && child.TotalMass > 0 && curNode != child { // Ensure the child is not nil before recursing

			distance := curNode.Center.Subtract(child.Center)
			magnitude := distance.Magnitude()
			s := node.NodeSize()
			if s/magnitude < theta || len(child.BodiesPtr.NodeBodies) == 1 {
				newForce := calculateGravitationalForce(*curNode, *child)
				curNode.BodiesPtr.NodeBodies[0].Force = curNode.BodiesPtr.NodeBodies[0].Force.Add(newForce)
				continue
			} else {
				updateForce(curNode, child, theta, dt)
			}
		}
	}
}

func (node *QuadNode) CalculateForce(root *QuadNode, theta float64, dt float64) {
	node.BodiesPtr.NodeBodies[0].Force = Vector2{0, 0}
	updateForce(node, root, theta, dt)
}

func (node *QuadNode) NodeSize() float64 {
	diagonal := node.Region[1].Subtract(node.Region[0])
	return diagonal.Magnitude()
}

func calculateDistance(pos1 [3]float64, pos2 [3]float64) float64 {

	x1, x2, x3 := pos1[0], pos1[1], pos1[2]
	y1, y2, y3 := pos2[0], pos2[1], pos2[2]

	dist := math.Sqrt(math.Pow(x1-y1, 2) + math.Pow(x2-y2, 2) + math.Pow(x3-y3, 2))

	return dist

}

func calculateGravitationalForce(node1 QuadNode, node2 QuadNode) Vector2 {
	r := node1.Center.Subtract(node2.Center)                                         // vector from body1 to body2
	distance := r.Magnitude()                                                        // scalar distance between bodies
	forceMagnitude := -G * node1.TotalMass * node2.TotalMass / (distance * distance) // magnitude of the gravitational force
	forceDirection := r.Normalize()                                                  // unit vector in the direction of the force
	return forceDirection.Multiply(forceMagnitude)                                   // vector representation of the force
}

// Assuming Vector2 and Body are already defined with basic methods

func (body *Body) Update(dt float64) {
	// Leapfrog Integration: update velocities at full step using accumulated force
	acceleration := body.Force.Multiply(1 / body.Mass)
	body.Velocities = body.Velocities.Add(acceleration.Multiply(dt))
	// Update position using new velocities
	changePos := body.Velocities.Multiply(dt)
	body.Positions = body.Positions.Add(changePos)
}

func Pop(n *[]*QuadNode) (*QuadNode, *[]*QuadNode) {
	if len(*n) == 0 {
		panic("cannot pop from an empty slice") // Or return an error instead of panicking
	}
	first, rest := (*n)[0], (*n)[1:]
	return first, &rest
}
