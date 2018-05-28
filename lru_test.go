package kv

import (
	"testing"
)

func Test_node_action(t *testing.T) {
	ll := NewLruList()

	data := [...]string{"tom", "jay", "anna"}

	for _, val := range data {
		item := &Item{
			value:      val,
			expireTime: 0,
		}
		ll.Set(val, item)
	}

	if ll.Head().data.value != data[len(data)-1] {
		t.Fatalf("Add the head node is failure, result %s, expect %s", ll.Head().data.value, data[len(data)-1])
	}

	if ll.Size() != int64(len(data)) {
		t.Fatalf("Add Link size is error, result %d, except %d", ll.Size(), len(data))
	}

	for _, val := range data {
		item := &Item{
			value:      val,
			expireTime: 0,
		}
		ll.Set(val, item)
	}

	if ll.Head().data.value != data[len(data)-1] {
		t.Fatalf("Set the head node is failure, result %s, expect %s", data[len(data)-1], ll.Head().data.value)
	}

	if ll.Size() != int64(len(data)) {
		t.Fatalf("Set Link size is error, result %d, except %d", ll.Size(), len(data))
	}

	ll.Del(data[len(data)-1])

	if ll.Head().data.value != data[len(data)-2] {
		t.Fatalf("Delete the node is failure, result %s, except %s", ll.Head().data.value, data[len(data)-2])
	}

	ll.Del(data[len(data)-2])

	if ll.Head().data.value != data[len(data)-3] {
		t.Fatalf("Delete the node is failure, result %s, except %s", ll.Head().data.value, data[len(data)-3])
	}

	ll.Del(data[len(data)-3])

	if ll.Head() != nil || ll.Tail() != nil {
		t.Fatalf("Delete the last node is failure, result %v, except %s", ll.Head(), "nil")
	}
}
