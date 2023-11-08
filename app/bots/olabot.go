package bots

import (
	ms "github.com/ar2rworld/botsservice/app/messageservice"
	"github.com/ar2rworld/botsservice/app/bot"
)

type OlaBot struct {
	bot.BaseBot
}

func (*OlaBot) GetName() string {
	return "olabot"
}

func (*OlaBot) HandleUpdate(u *ms.Update) (ms.MessageReply, error) {
	if u.GetText() == "a" {
		return ms.MessageReply{ Text: "b", ChatID: u.GetChatID(), UserID: u.GetUserID()}, nil
	}
	return ms.MessageReply{}, nil
}

func NewOlaBot() *OlaBot {
	return &OlaBot{}
}
