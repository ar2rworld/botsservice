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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/ar2rworld/botsservice/app/messageservice"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	// Contact the server and print out its response.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	// startMessage := tgbotapi.NewMessage(int64(-1001506079405), "Hello from " + name)
	// _, err = bot.Send(startMessage)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)


	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("Sending update as %s", name)

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
			log.Println(err)
		} else {
			for {
				mr, err := stream.Recv()
				if err == io.EOF {
						break
				}
				if err != nil {
						log.Fatalf("%v.SendUpdates(_) = _, %v", c, err)
				}
				if mr.GetText() == "" {
					continue
				}
				m := tgbotapi.NewMessage(mr.GetChatID(), mr.GetText())
				_, err = bot.Send(m)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
