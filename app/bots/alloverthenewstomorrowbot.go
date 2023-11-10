package bots

import (
	"context"
	"fmt"
	"strings"

	"github.com/ar2rworld/botsservice/app/bot"
	ms "github.com/ar2rworld/botsservice/app/messageservice"
)

type AllOverTheNewsTomorrowBot struct {
	bot.BaseBot
}


func (*AllOverTheNewsTomorrowBot) GetName() string {
	return "all_over_the_news_tomorrow_bot"
}

func (b *AllOverTheNewsTomorrowBot) HandleUpdate(u *ms.Update) (ms.MessageReply, error) {
	ctx := context.TODO()

	switch u.GetText() {
		case "b": {
			return ms.MessageReply{ ChatID: u.GetChatID(), UserID: u.GetUserID(), Text: "c" }, nil
		}
		case "/help": {
			return ms.MessageReply{ Text: `/tomorrow to get a random prediction
				/subscribe - to subsribe and to unsubscribe from daily notification
				/whoAreYouByDima dayIndex monthIndex //integers please`,
				ChatID: u.GetChatID(),
				UserID: u.GetUserID(),
			}, nil
		}
		case "/tomorrow": {
			phrase, err := GetRandomPhrase(ctx, b.GetDatabase())
			if err != nil {
				return ms.MessageReply{}, fmt.Errorf("Error GetRandomPhrase: %v", err)
			}

			return ms.MessageReply{
				Text: phrase,
				ChatID: u.GetChatID(),
				UserID: u.GetUserID(),
			}, nil
		}
		default:
			args, found := strings.CutPrefix(u.GetText(), "/whoAreYouByDima")
			if ! found {
				return ms.MessageReply{}, nil
			}

			who, err := GetWhoAreYouByDima(ctx, b.GetDatabase(), args)
			if err != nil {
				return ms.MessageReply{}, err
			}

			return ms.MessageReply{ Text: who, ChatID: u.GetChatID(), UserID: u.GetUserID() }, nil
	}
}

func NewAllOverTheNewsTomorrowBot() *AllOverTheNewsTomorrowBot {
	return &AllOverTheNewsTomorrowBot{}
}
