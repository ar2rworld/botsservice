package bots

import (
	"github.com/ar2rworld/botsservice/app/bot"
	ms "github.com/ar2rworld/botsservice/app/messageservice"
)

type AllOverTheNewsTomorrowBot struct {
	bot.BaseBot
}


func (*AllOverTheNewsTomorrowBot) GetName() string {
	return "all_over_the_news_tomorrow_bot"
}

func (b *AllOverTheNewsTomorrowBot) GetToken() string {
	return b.token
}

func (b *AllOverTheNewsTomorrowBot) HandleUpdate(u *ms.Update) (ms.MessageReply, error) {
	if u.GetText() == "b" {
		return ms.MessageReply{ Text: "c", ChatID: u.GetChatID(), UserID: u.GetUserID() }, nil
	}
	return ms.MessageReply{}, nil
}

func NewAllOverTheNewsTomorrowBot() *AllOverTheNewsTomorrowBot {
	return &AllOverTheNewsTomorrowBot{}
}
