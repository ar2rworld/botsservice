/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/ar2rworld/botsservice/app/messageservice"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var addr = os.Getenv("MESSAGESERVICE_ADDRESS")
	if addr == "" {
		log.Fatal("Did not find MESSAGESERVICE_ADDRESS env var")
	}
	var port = os.Getenv("MESSAGESERVICE_PORT")
	if port == "" {
		log.Fatal("Did not find MESSAGESERVICE_PORT env var")
	}
	var name = os.Getenv("NAME")
	if name == "" {
		log.Fatal("Did not find NAME env var")
	}
	
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMessageServiceClient(conn)

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	t, err := c.Register(ctx, &pb.RegisterRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	if t.GetToken() != "" {
		log.Println("Token: ok")
	}

	bot, err := tgbotapi.NewBotAPI(t.GetToken())
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)

	g.Go(func() error {
		tiker := time.NewTicker(time.Second)
		
		for {
			select {
				case <- tiker.C:
					stream, err := c.SendUpdates(ctx, &pb.Updates{ Botname: name })
					if err != nil {
						return err
					}
					err = CheckUpdates(stream, bot)
					if err != nil {
						return err
					}
				case <- ctx.Done():
					return nil
			}
		}
	})

	g.Go(func() error {
		ticker := time.NewTicker(time.Second)

		for {
			select {
				case <- ticker.C:
					if err != nil {
						return err
					}

					updates, err := bot.GetUpdates(updateConfig)
					if err != nil {
						return err
					}

					for _, update := range updates {
						if update.UpdateID >= updateConfig.Offset {
							updateConfig.Offset = update.UpdateID + 1
						}

						if update.Message == nil {
							continue
						}
			
						u := &pb.Update{
							ChatID: update.Message.Chat.ID,
							UserID: update.Message.From.ID,
							Text: update.Message.Text,
						}
			
						stream, err := c.SendUpdates(
							ctx,
							&pb.Updates{
								Botname: name,
								Updates: []*pb.Update{ u },
							},
						)
						if err != nil {
							return err
						} else {
							err = CheckUpdates(stream, bot)
							if err != nil {
								return err
							}
						}
					}
				case <- ctx.Done():
					return nil
			}
		}
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("Client error: %v", err)
	}
}

func CheckUpdates(stream pb.MessageService_SendUpdatesClient, bot *tgbotapi.BotAPI) error {
	for {
		mr, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if mr.GetText() == "" {
			continue
		}
		m := tgbotapi.NewMessage(mr.GetChatID(), mr.GetText())
		_, err = bot.Send(m)
		if err != nil {
			return err
		}
	}
	return nil
}
