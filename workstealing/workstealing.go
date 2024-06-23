package workstealing

import (
	"fmt"
	"proj3-redesigned/utils"
	"sync/atomic"
)

type Task interface {
	Execute()
}

type TerminationTask struct{}

func (t *TerminationTask) Execute() {}

type NodeTask struct {
	Node  *utils.QuadNode
	Root  *utils.QuadNode
	Theta float64
	Dt    float64
}

func (nt *NodeTask) Execute() {
	if nt.Node.BodiesPtr != nil && nt.Node.BodiesPtr.NodeBodies != nil {
		nt.Node.CalculateForce(nt.Root, nt.Theta, nt.Dt)
	}
}

type BodyTask struct {
	Body *utils.Body
	Dt   float64
}

func (bt *BodyTask) Execute() {
	bt.Body.Update(bt.Dt)
}

type node struct {
	task    *Task
	next    atomic.Pointer[node]
	prev    atomic.Pointer[node]
	removed atomic.Bool // Atomic boolean for logical deletion
}

type Dequeue struct {
	head atomic.Pointer[node]
	tail atomic.Pointer[node]
}

func NewWorkStealingDequeue() *Dequeue {
	n := &node{}
	dq := &Dequeue{}
	dq.head.Store(n)
	dq.tail.Store(n)
	return dq
}

func (dq *Dequeue) PrintQueue() {
	curNode := dq.head.Load()
	for curNode != nil {
		if curNode.task != nil {
			task := *curNode.task
			fmt.Printf("Task: %+v\n", task)
		} else {
			fmt.Println("Empty node found")
		}
		curNode = curNode.next.Load()
	}
}

func (dq *Dequeue) Push(task Task) {
	newNode := &node{task: &task}
	for {
		tail := dq.tail.Load()
		if tail.next.CompareAndSwap(nil, newNode) {
			newNode.prev.Store(tail) // Link back to the old tail
			dq.tail.Store(newNode)   // Update the tail to the new node
			return
		}
	}
}
func (dq *Dequeue) Pop() (Task, bool) {
	for {
		head := dq.head.Load()
		if head == nil || head.removed.Load() {
			return nil, false
		}
		next := head.next.Load()
		if next == nil {
			return nil, false
		}

		if !head.removed.Load() && head.removed.CompareAndSwap(false, true) {
			if dq.head.CompareAndSwap(head, next) {
				next.prev.Store(nil)
				return *next.task, true
			}
		}
	}
}

func (dq *Dequeue) Steal() (Task, bool) {
	for {
		tail := dq.tail.Load()
		if tail == dq.head.Load() || tail.removed.Load() {
			return nil, false
		}
		prev := tail.prev.Load()
		if prev != nil && !tail.removed.Load() && tail.removed.CompareAndSwap(false, true) {
			if dq.tail.CompareAndSwap(tail, prev) {
				prev.next.Store(nil)
				return *tail.task, true
			}
		}
	}
}
