package messagequeue

import (
	"testing"
	
	pb "github.com/ar2rworld/botsservice/app/messageservice"
)

func TestMessageQueue(t *testing.T) {
	q := NewMessageQueue()
	if q.Len() != 0 {
		t.Error("New queue should be empty")
	}
	_, err := q.Pop()
	if err == nil {
		t.Error("New queue pop should return error")
	}
	_, err = q.Pop()
	if err != nil && err.Error() != "Queue is empty" {
		t.Error("New queue pop should return \"Queue is empty\" error")
	}

	q.Push(&pb.MessageReply{Text: "1"})
	q.Push(&pb.MessageReply{Text: "2"})
	q.Push(&pb.MessageReply{Text: "3"})

	if q.Len() != 3 {
		t.Error("Queue should have 3 elements")
	}
	m, err := q.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if m.GetText() != "1" {
		t.Errorf(`Incorrect messageReply returned on pop: "%s" != "1"`, m.GetText())
	}
	if q.Len() != 2 {
		t.Error("Queue should have 2 elements")
	}
	m, err = q.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if m.GetText() != "2" {
		t.Errorf(`Incorrect messageReply returned on pop: "%s" != "2"`, m.GetText())
	}
	m, err = q.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if m.GetText() != "3" {
		t.Errorf(`Incorrect messageReply returned on pop: "%s" != "3"`, m.GetText())
	}
	if q.Len() != 0 {
		t.Error("Queue should have 0 elements")
	}
}