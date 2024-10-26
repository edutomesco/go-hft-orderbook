package internal

import (
	"sync"
)

// Doubly linked orders queue
// TODO: this should be compared with ring buffer queue performance
type ordersQueue struct {
	head *Order
	tail *Order
	size int
	mu   sync.Mutex
}

func NewOrdersQueue() ordersQueue {
	return ordersQueue{}
}

func (this *ordersQueue) Size() int {
	return this.size
}

func (this *ordersQueue) IsEmpty() bool {
	return this.size == 0
}

func (this *ordersQueue) Enqueue(o *Order) {
	this.mu.Lock()
	defer this.mu.Unlock()

	tail := this.tail
	this.tail = o
	if tail != nil {
		tail.Next = o
		o.Prev = tail
	}
	if this.head == nil {
		this.head = o
	}
	this.size++
}

func (this *ordersQueue) Dequeue() *Order {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.size == 0 {
		return nil
	}

	head := this.head
	if this.tail == this.head {
		this.tail = nil
	}

	this.head = this.head.Next
	this.size--
	return head
}

func (this *ordersQueue) Delete(o *Order) {
	this.mu.Lock()
	defer this.mu.Unlock()

	prev := o.Prev
	next := o.Next
	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}
	o.Next = nil
	o.Prev = nil

	this.size--

	if this.head == o {
		this.head = next
	}
	if this.tail == o {
		this.tail = prev
	}
}
