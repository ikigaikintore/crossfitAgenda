package domain

type content struct {
	day, name string
}

type nodeExercise struct {
	element content
	next    *nodeExercise
	prev    *nodeExercise
}

type listNodes struct {
	head *nodeExercise
	tail *nodeExercise
	size int
}

func newListNodes() *listNodes {
	return &listNodes{
		size: 0,
	}
}

func (ln *listNodes) Insert(name, day string) {
	node := &nodeExercise{
		element: content{
			day:  day,
			name: name,
		},
	}

	if ln.head == nil {
		ln.head = node
		ln.tail = node
		ln.size++
		return
	}

	node.next = ln.head
	ln.head.prev = node
	ln.head = node
	ln.size++
}

func (ln *listNodes) Get() *nodeExercise {
	if ln.IsEmpty() {
		return nil
	}

	el := ln.tail
	ln.tail = ln.tail.prev
	if ln.tail != nil {
		ln.tail.next = nil
	}
	ln.size--
	return el
}

func (ln *listNodes) IsEmpty() bool {
	return ln == nil || ln.size <= 0
}
