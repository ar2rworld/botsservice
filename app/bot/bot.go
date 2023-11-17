package bot

import (
	"github.com/ar2rworld/botsservice/app/messagequeue"
	pb "github.com/ar2rworld/botsservice/app/messageservice"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bot interface {
	HandleUpdate(*pb.Update) (pb.MessageReply, error)
	GetName() string
	GetToken() string
	SetToken(string)
	SetDatabase(*mongo.Database)
	GetDatabase() (*mongo.Database)
	SetScheduler(*gocron.Scheduler)
	GetScheduler() *gocron.Scheduler
	SetupScheduler()
	SetQueue(*messagequeue.MessageQueue)
	GetQueue() *messagequeue.MessageQueue
}

type BaseBot struct {
	token string

	messagequeue *messagequeue.MessageQueue
	scheduler *gocron.Scheduler
	Database *mongo.Database
}

func (b *BaseBot) GetToken() string {
	return b.token
}

func (b *BaseBot) SetToken(t string) {
	b.token = t
}

func (b *BaseBot) GetDatabase() (*mongo.Database) {
	return b.Database
}

func (b *BaseBot) SetDatabase(db *mongo.Database) {
	b.Database = db
}

func (b *BaseBot) SetScheduler(s *gocron.Scheduler) {
	b.scheduler = s
}
func (b *BaseBot) GetScheduler() *gocron.Scheduler {
	return b.scheduler
}
func (b *BaseBot) SetupScheduler() {}

func (b *BaseBot) GetQueue() *messagequeue.MessageQueue {
	return b.messagequeue
}
func (b *BaseBot) SetQueue(q *messagequeue.MessageQueue) {
	b.messagequeue = q
}
