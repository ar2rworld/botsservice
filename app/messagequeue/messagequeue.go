package messagequeue

import (
	"container/list"
	"fmt"

	pb "github.com/ar2rworld/botsservice/app/messageservice"
)

func NewMessageQueue(n string) MessageQueue {
	return MessageQueue{
		name: n,
		queue: list.New(),
	}
}

type MessageQueue struct {
	name string
	queue *list.List
}

func (q *MessageQueue) GetName() string {
	return q.name
}

func (q *MessageQueue) Pop() (pb.MessageReply, error) {
	if q.Len() == 0 {
		return pb.MessageReply{}, fmt.Errorf("Queue is empty")
	}
	el := q.queue.Front()
	q.queue.Remove(el)
	m, ok := el.Value.(pb.MessageReply)
	
	if !ok {
		return pb.MessageReply{}, fmt.Errorf("Could not cast to Message")
	}
	
	return m, nil
}
func (q *MessageQueue) Push(m pb.MessageReply) {
	q.queue.PushBack(m)
}
func (q *MessageQueue) Len() int {
	return q.queue.Len()
}
