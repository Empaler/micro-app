package exercises

import "testing"

// Expected implementation contract (to be provided by you):
type DoublyLinkedList struct {
	head *DoublyNode
	tail *DoublyNode
}

type DoublyNode struct {
	value int
	next  *DoublyNode
	prev  *DoublyNode
}

func (l *DoublyLinkedList) Append(value int) {
	node := &DoublyNode{value: value}
	if l.head == nil {
		l.head = node
		l.tail = node
		return
	}

	currentNode := l.head
	for currentNode != nil {
		if currentNode.next == nil {
			node.prev = currentNode
			currentNode.next = node
			l.tail = node
			break
		}
		currentNode = currentNode.next
	}
}

func (l *DoublyLinkedList) Prepend(value int) {
	node := &DoublyNode{value: value}
	if l.head == nil {
		l.head = node
		l.tail = node
		return
	}

	l.head.prev = node
	node.next = l.head
	l.head = node
}

func (l *DoublyLinkedList) Remove(value int) bool {
	if l.head == nil {
		return false
	}

	if l.head == l.tail && l.head.value == value {
		l.head = nil
		l.tail = nil
		return true
	}

	if l.head.value == value {
		l.head = l.head.next
		l.head.prev = nil
		return true
	}

	currentNode := l.head.next
	for currentNode != nil {
		if currentNode.value == value {
			if currentNode.prev != nil {
				currentNode.prev.next = currentNode.next
			} else {
				l.head = currentNode.next
			}

			if currentNode.next != nil {
				currentNode.next.prev = currentNode.prev
			} else {
				l.tail = currentNode.prev
			}
			return true
		}
		currentNode = currentNode.next
	}

	return false
}

func (l *DoublyLinkedList) ForwardSlice() []int {
	fowardSlice := []int{}
	if l.head == nil {
		return fowardSlice
	}

	node := l.head
	for node != nil {
		fowardSlice = append(fowardSlice, node.value)
		node = node.next
	}

	return fowardSlice
}

func (l *DoublyLinkedList) BackwardSlice() []int {
	backwardSlice := []int{}
	if l.head == nil {
		return backwardSlice
	}

	node := l.tail
	for node != nil {
		backwardSlice = append(backwardSlice, node.value)
		node = node.prev
	}

	return backwardSlice
}

func TestDoublyLinkedListAppendAndTraverse(t *testing.T) {
	list := &DoublyLinkedList{}
	list.Append(1)
	list.Append(2)
	list.Append(3)

	forward := list.ForwardSlice()
	backward := list.BackwardSlice()

	wantForward := []int{1, 2, 3}
	wantBackward := []int{3, 2, 1}

	for i := range wantForward {
		if forward[i] != wantForward[i] {
			t.Fatalf("forward mismatch at index %d: got %d, want %d", i, forward[i], wantForward[i])
		}
	}
	for i := range wantBackward {
		if backward[i] != wantBackward[i] {
			t.Fatalf("backward mismatch at index %d: got %d, want %d", i, backward[i], wantBackward[i])
		}
	}
}

func TestDoublyLinkedListRemove(t *testing.T) {
	list := &DoublyLinkedList{}
	list.Append(1)
	list.Append(2)
	list.Append(3)
	list.Append(4)

	if !list.Remove(1) {
		t.Fatalf("expected to remove head")
	}
	if !list.Remove(3) {
		t.Fatalf("expected to remove middle")
	}
	if !list.Remove(4) {
		t.Fatalf("expected to remove tail")
	}
	if list.Remove(999) {
		t.Fatalf("did not expect remove to succeed for missing value")
	}

	got := list.ForwardSlice()
	want := []int{2}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("unexpected list after removes: got %v, want %v", got, want)
	}
}
