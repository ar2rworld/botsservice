package bots

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ar2rworld/botsservice/app/bot"
	ms "github.com/ar2rworld/botsservice/app/messageservice"
	"go.mongodb.org/mongo-driver/bson"
)

const subscribes = "subscribes"

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
				/subscribe - to subscribe and to unsubscribe from daily notification
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
		case "/subscribe": {
			chatID := u.Message.Chat.ID
			var chat bot.Chat

			err := b.GetDatabase().
				Collection(subscribes).
				FindOne(ctx, bson.D{{ Key: "id", Value: chatID }}).
				Decode(&chat)
			if err != nil && err.Error() == NoDocuments {
				// subscribe
				chat := &bot.Chat{
					Name: u.Message.Chat.Title,
					ID: u.Message.Chat.ID,
					IsActive: true,
				}
				_, err = b.GetDatabase().
					Collection(subscribes).
					InsertOne(ctx, chat.ToDoc())

				if err != nil {
					return ms.MessageReply{}, err
				}

				return ms.MessageReply{ Text: "Successfully subsribed", ChatID: chatID }, nil
			}
			if err != nil {
				return ms.MessageReply{}, err
			}

			if chat.IsActive {
				// should change status , update doc, send "successfully unsubscribed"

				chat.IsActive = false

				_, err = b.GetDatabase().
					Collection(subscribes).
					UpdateOne(ctx, bson.D{{ Key: "id", Value: chat.ID }}, bson.D{{ Key: "$set", Value: chat.ToDoc() }})
				
				if err != nil {
					return ms.MessageReply{}, err
				}

				return ms.MessageReply{ Text: "Successfully unsubscribed", ChatID: chatID }, nil
			} else {
				// should change status , update doc, send "successfully subscribed"
				
				chat.IsActive = true

				_, err = b.GetDatabase().
					Collection(subscribes).
					UpdateOne(ctx, bson.D{{ Key: "id", Value: chat.ID }}, bson.D{{ Key: "$set", Value: chat.ToDoc() }})

				return ms.MessageReply{ Text: "Successfully subscribed", ChatID: chatID }, nil
			}
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
	return ms.MessageReply{}, nil
}

func (b *AllOverTheNewsTomorrowBot) SetupScheduler() {
	b.GetScheduler().Every(1).Day().At("17:00").Do(func() {
		// query IsActive chats and send out messages
		ctx := context.TODO()
		
		c, err := b.GetDatabase().Collection(subscribes).Find(ctx, bson.D{{ Key: "isactive", Value: true }})
		if err != nil {
			log.Printf("Scheduler error: %v", err)
			return
		}

		var chats []bot.Chat
		err = c.All(ctx, &chats)
		if err != nil {
			log.Printf("Scheduler error: %v", err)
			return
		}

		for _, chat := range chats {

			phrase, err := GetRandomPhrase(ctx, b.GetDatabase())
			if err != nil {
				log.Printf("Error GetRandomPhrase: %v", err)
				return
			}

			b.GetQueue().Push(ms.MessageReply{
				Text: phrase,
				ChatID: chat.ID,
			})
		}
	})
}

func NewAllOverTheNewsTomorrowBot() *AllOverTheNewsTomorrowBot {
	return &AllOverTheNewsTomorrowBot{}
}
