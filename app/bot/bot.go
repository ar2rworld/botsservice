package bot

import (
	pb "github.com/ar2rworld/botsservice/app/messageservice"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bot interface {
	HandleUpdate(*pb.Update) (pb.MessageReply, error)
	GetName() string
	GetToken() string
	SetToken(string)
	SetDatabase(*mongo.Database)
	GetDatabase() (*mongo.Database)
}

type BaseBot struct {
	token string

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
