package exercises

import "testing"

// Expected implementation contract (to be provided by you):
func (l *SinglyLinkedList) Append(value int) {
	node := &Node{value: value}
	if l.head == nil {
		l.head = node
		return
	}

	currentNode := l.head
	for currentNode != nil {
		if currentNode.next == nil {
			currentNode.next = node
			break
		}
		currentNode = currentNode.next
	}
}
func (l *SinglyLinkedList) Prepend(value int) {
	node := &Node{value: value}
	if l.head == nil {
		l.head = node
		return
	}

	node.next = l.head
	l.head = node
}
func (l *SinglyLinkedList) Delete(value int) bool {
	if l.head == nil {
		return false
	}
	if l.head.value == value {
		l.head = l.head.next
		return true
	}

	node := l.head.next
	previous := l.head
	for node != nil {
		if node.value == value {
			previous.next = node.next
			return true
		}
		previous = node
		node = node.next
	}
	return false
}

func (l *SinglyLinkedList) ToSlice() []int {
	sliceList := []int{}
	if l.head == nil {
		return sliceList
	}

	node := l.head
	for node != nil {
		sliceList = append(sliceList, node.value)
		node = node.next
	}

	return sliceList
}

type Node struct {
	value int
	next  *Node
}

type SinglyLinkedList struct {
	head *Node
}

func TestSinglyLinkedListAppendAndPrepend(t *testing.T) {
	list := &SinglyLinkedList{}
	list.Append(2)
	list.Append(3)
	list.Prepend(1)

	got := list.ToSlice()
	want := []int{1, 2, 3}

	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("value mismatch at index %d: got %d, want %d", i, got[i], want[i])
		}
	}
}

func TestSinglyLinkedListDelete(t *testing.T) {
	list := &SinglyLinkedList{}
	list.Append(1)
	list.Append(2)
	list.Append(3)
	list.Append(4)

	if !list.Delete(1) {
		t.Fatalf("expected to delete head")
	}
	if !list.Delete(3) {
		t.Fatalf("expected to delete middle node")
	}
	if !list.Delete(4) {
		t.Fatalf("expected to delete tail")
	}
	if list.Delete(999) {
		t.Fatalf("did not expect delete to succeed for missing value")
	}

	got := list.ToSlice()
	want := []int{2}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("unexpected list after deletes: got %v, want %v", got, want)
	}
}
