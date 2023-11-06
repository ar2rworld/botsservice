package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ar2rworld/botsservice/app/bots"
	mq "github.com/ar2rworld/botsservice/app/messagequeue"

	pb "github.com/ar2rworld/botsservice/app/messageservice"
	"google.golang.org/grpc"
)

type Bot interface {
	HandleUpdate(*pb.Update) (pb.MessageReply, error)
	GetName() string
	GetToken() string
}

func main() {
	var addr = os.Getenv("MESSAGESERVICE_ADDRESS")
	if addr == "" {
		log.Fatal("Did not find MESSAGESERVICE_ADDRESS env var")
	}
	var port = os.Getenv("MESSAGESERVICE_PORT")
	if port == "" {
		log.Fatal("Did not find MESSAGESERVICE_PORT env var")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
  	log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	if os.Getenv("olabot_token") == "" {
		log.Fatal("Missing olabot_token in env")
	}
	if os.Getenv("all_over_the_news_tomorrow_bot_token") == "" {
		log.Fatal("Missing all_over_the_news_tomorrow_bot_token in env")
	}
	server := newServer()
	server.AddBot(bots.NewOlaBot(os.Getenv("olabot_token")))
	server.AddBot(bots.NewAllOverTheNewsTomorrowBot(os.Getenv("all_over_the_news_tomorrow_bot_token")))

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMessageServiceServer(grpcServer, server)

	log.Printf("Started gRPC at: %s:%s", addr, port)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving grpc: %v", err)
	}
}

type server struct {
	queue []mq.MessageQueue
	bots []Bot
	pb.UnimplementedMessageServiceServer
}

func newServer() *server {
	return &server{ bots: []Bot{}, queue: []mq.MessageQueue{}}
}

func (s *server) AddBot (b Bot) {
	s.bots = append(s.bots, b)
}

func (s *server) GetBot(name string) (Bot, error) {
	for _, b := range s.bots {
		if b.GetName() == name {
			return b, nil
		}
	}
	return nil, fmt.Errorf("Could not find bot: %s", name)
}

func (s *server) GetBotQueue(name string) (mq.MessageQueue, error) {
	var messagequeue mq.MessageQueue
	for _, mq := range s.queue {
		if mq.GetName() == name {
			return messagequeue, nil
		}
	}

	return mq.MessageQueue{}, fmt.Errorf("Could not find MessageQueue for %s", name)
}

func (s *server) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.TokenReply, error) {
	var bot Bot
	var n = r.GetName()
	for _, b := range s.bots {
		if b.GetName() == n {
			bot = b
		}
	}
	if bot == nil {
		return &pb.TokenReply{}, fmt.Errorf("Could not find bot: %s", n)
	}

	s.queue = append(s.queue, mq.NewMessageQueue(n))

	log.Printf("Registed: %s", n)
	return &pb.TokenReply{ Token: bot.GetToken() }, nil
}

func (s *server) SendUpdates(u *pb.Updates, stream pb.MessageService_SendUpdatesServer) error {
	bn := u.GetBotname()
	q, err := s.GetBotQueue(bn)
	if err != nil {
		log.Fatal(err)
	}
	for q.Len() > 0 {
		m, err := q.Pop()
		if err != nil {
			return err
		}

		if err = stream.Send(&m); err != nil {
			return err
		}
	}

	// TODO: check HandleUpdate on bot
	
	bot, err := s.GetBot(bn)
	if err != nil {
		return err
	}
	
	for _, update := range u.Updates {
		r, err := bot.HandleUpdate(update)
		if err != nil {
			return err
		}
		if err = stream.Send(&r); err != nil {
			return err
		}
	}
	
	return nil	
}
