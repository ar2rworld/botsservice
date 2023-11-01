package bots

import ms "github.com/ar2rworld/botsservice/app/messageservice"

type AllOverTheNewsTomorrowBot struct {
	token string
}


func (*AllOverTheNewsTomorrowBot) GetName() string {
	return "alloverthenewstomorrowbot"
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

func NewAllOverTheNewsTomorrowBot(t string) *AllOverTheNewsTomorrowBot {
	return &AllOverTheNewsTomorrowBot{token: t}
}
