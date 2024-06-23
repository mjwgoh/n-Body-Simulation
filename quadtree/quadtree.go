// This package creates an quadtree based on some given inputs.
// Allowing us to use the Barnes Hut Algorithm for calculation of the movement of the bodies.

package quadtree

import (
	"fmt"
	"proj3-redesigned/utils"
)

func newOctreeNode(region [2]utils.Vector2) *utils.QuadNode {
	node := &utils.QuadNode{
		Region:    region,
		BodiesPtr: &utils.Bodies{},
		Children:  [4]*utils.QuadNode{}, // Initialize all children to nil explicitly
	}
	return node
}

// Helper function to calculate the new region for a child
func childRegion(index int, parentRegion [2]utils.Vector2) [2]utils.Vector2 {

	//Establish child region
	var newRegion [2]utils.Vector2

	//Obtain the midpoint for x, y, and z. mid: [x_mid, y_mid, z_mid]
	midX := (parentRegion[0].X + parentRegion[1].X) / 2
	midY := (parentRegion[0].Y + parentRegion[1].Y) / 2

	// Performs bitwise calculations. Note that 2^2 = 4.

	if (index & 1) == 1 {
		newRegion[0].X = midX
		newRegion[1].X = parentRegion[1].X
	} else {
		newRegion[0].X = parentRegion[0].X
		newRegion[1].X = midX
	}

	if (index & 2) == 2 {
		newRegion[0].Y = midY
		newRegion[1].Y = parentRegion[1].Y
	} else {
		newRegion[0].Y = parentRegion[0].Y
		newRegion[1].Y = midY
	}

	return newRegion
}

func insertBody(node *utils.QuadNode, body *utils.Body, root *utils.QuadNode) {
	node.TotalMass += body.Mass
	newWeightedX := body.Positions.X * body.Mass
	newWeightedY := body.Positions.Y * body.Mass

	// Update center of mass
	if node.TotalMass > 0 {
		node.Center.X = (node.Center.X*(node.TotalMass-body.Mass) + newWeightedX) / node.TotalMass
		node.Center.Y = (node.Center.Y*(node.TotalMass-body.Mass) + newWeightedY) / node.TotalMass
	}

	if node.IsLeaf() {
		node.BodiesPtr.NodeBodies = append(node.BodiesPtr.NodeBodies, body)
		if len(node.BodiesPtr.NodeBodies) >= 2 {
			subdivide(node, root)
		}
		return
	}

	// Insert the body into the appropriate child
	for i, child := range node.Children {
		if child != nil && isWithinRegion(body.Positions, child.Region) {
			insertBody(child, body, root)
			return
		} else if child == nil && isWithinRegion(body.Positions, childRegion(i, node.Region)) {
			node.Children[i] = newOctreeNode(childRegion(i, node.Region))
			insertBody(node.Children[i], body, root)
			return
		}
	}
}

func subdivide(node *utils.QuadNode, root *utils.QuadNode) {
	// Initialize child nodes
	for i := 0; i < 4; i++ {
		node.Children[i] = newOctreeNode(childRegion(i, node.Region))
	}

	// Redistribute bodies
	allBodies := node.BodiesPtr.NodeBodies // Copy to avoid modification issues during iteration
	node.BodiesPtr.NodeBodies = nil        // Clear the parent node's body list immediately to avoid duplication

	for _, body := range allBodies {
		// Insert each body into the appropriate child node
		inserted := false
		for _, child := range node.Children {
			if isWithinRegion(body.Positions, child.Region) {
				insertBody(child, body, root)
				inserted = true
				break
			}
		}
		// Optionally handle the case where no child is suitable (e.g., if the body is on a boundary)
		if !inserted {
			fmt.Println("Error: Body could not be inserted into any child node")
		}
	}
}

// Helper function to check if a position is within a region
func isWithinRegion(position utils.Vector2, region [2]utils.Vector2) bool {
	return position.X >= region[0].X && position.X < region[1].X &&
		position.Y >= region[0].Y && position.Y < region[1].Y
}

// Build the Octree
func BuildQuadTree(bodies []*utils.Body, region [2]utils.Vector2) *utils.QuadNode {
	root := newOctreeNode(region)
	for _, body := range bodies {
		insertBody(root, body, root)
	}
	return root
}
