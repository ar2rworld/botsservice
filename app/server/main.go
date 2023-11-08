package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/ar2rworld/botsservice/app/bot"
	"github.com/ar2rworld/botsservice/app/bots"
	mq "github.com/ar2rworld/botsservice/app/messagequeue"

	pb "github.com/ar2rworld/botsservice/app/messageservice"
	"google.golang.org/grpc"
)

type BotQueue struct {
	Bot bot.Bot
	Queue mq.MessageQueue
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

	if os.Getenv("ADMIN_ID") == "" {
		log.Fatal("Missing olabot_token in env")
	}
	adminID, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		log.Fatal("Error converting ADMIN_ID")
	}


	server := newServer()
	server.AdminID = adminID
	server.DBClient = *DBClient
	err = server.AddBot(bots.NewOlaBot())
	if err != nil {
		log.Fatalf("Error adding a bot: %v", err)
	}
	err = server.AddBot(bots.NewAllOverTheNewsTomorrowBot())
	if err != nil {
		log.Fatalf("Error adding a bot: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMessageServiceServer(grpcServer, server)

	log.Printf("Started gRPC at: %s:%s", addr, port)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving grpc: %v", err)
	}
}

type server struct {
	bq map[string]BotQueue
	pb.UnimplementedMessageServiceServer

	AdminID int64
}

func newServer() *server {
	return &server{ bq: map[string]BotQueue{} }
}

func (s *server) AddBot (b bot.Bot) error {
	name := b.GetName()

	token := os.Getenv(fmt.Sprintf("%s_token", name))
	if token == "" {
		return fmt.Errorf("Missing %s_token in env", name)
	}

	b.SetToken(token)
	s.bq[name] = BotQueue{ Bot: b, Queue: mq.NewMessageQueue(name)}
	return nil
}

func (s *server) GetBotQueue(name string) (BotQueue, error) {
	bq := s.bq[name]

	if bq.Bot == nil {
		return BotQueue{}, fmt.Errorf("Could not find bot: %s", name)
	}
	return bq, nil
}

func (s *server) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.TokenReply, error) {
	var name = r.GetName()
	
	bq, err := s.GetBotQueue(name)
	if err != nil {
		return &pb.TokenReply{}, fmt.Errorf("Error occured on GetBotQueue: %v", err)
	}

	log.Printf("Registed: %s", bq.Bot.GetName())

	bq.Queue.Push(pb.MessageReply{ Text: "Hello", ChatID: s.AdminID, UserID: s.AdminID })

	return &pb.TokenReply{ Token: bq.Bot.GetToken() }, nil
}

func (s *server) SendUpdates(u *pb.Updates, stream pb.MessageService_SendUpdatesServer) error {
	bn := u.GetBotname()
	bq, err := s.GetBotQueue(bn)
	if err != nil {
		log.Fatal(err)
	}
	for bq.Queue.Len() > 0 {
		m, err := bq.Queue.Pop()
		if err != nil {
			return err
		}

		if err = stream.Send(&m); err != nil {
			return err
		}
	}
	
	for _, update := range u.Updates {
		r, err := bq.Bot.HandleUpdate(update)
		if err != nil {
			return err
		}
		if err = stream.Send(&r); err != nil {
			return err
		}
	}
	
	return nil	
}
