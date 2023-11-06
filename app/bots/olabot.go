package bots

import ms "github.com/ar2rworld/botsservice/app/messageservice"

type OlaBot struct {
	token string
}

func (*OlaBot) GetName() string {
	return "olabot"
}

func (b *OlaBot) GetToken() string {
	return b.token
}

func (*OlaBot) HandleUpdate(u *ms.Update) (ms.MessageReply, error) {
	if u.GetText() == "a" {
		return ms.MessageReply{ Text: "b", ChatID: u.GetChatID(), UserID: u.GetUserID()}, nil
	}
	return ms.MessageReply{}, nil
}

func NewOlaBot(t string) *OlaBot {
	return &OlaBot{ token: t}
}
