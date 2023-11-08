package bot

import pb "github.com/ar2rworld/botsservice/app/messageservice"

type Bot interface {
	HandleUpdate(*pb.Update) (pb.MessageReply, error)
	GetName() string
	GetToken() string
	SetToken(string)
}

type BaseBot struct {
	token string
}

func (b *BaseBot) GetToken() string {
	return b.token
}

func (b *BaseBot) SetToken(t string) {
	b.token = t
}
