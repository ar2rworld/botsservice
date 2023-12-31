package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ar2rworld/botsservice/app/bot"
	"github.com/ar2rworld/botsservice/app/bots"
	"github.com/ar2rworld/botsservice/app/db"
	mq "github.com/ar2rworld/botsservice/app/messagequeue"

	pb "github.com/ar2rworld/botsservice/app/messageservice"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/mongo"
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

	DBClient, err := db.NewDBClient()
	if err != nil {
		log.Fatal(err)
	}

	scheduler := gocron.NewScheduler(time.UTC)

	server := newServer()
	server.Scheduler = scheduler
	server.AdminID = adminID
	server.DBClient = *DBClient
	err = server.AddBot(bots.NewOlaBot(), bot.BotConfig{})
	if err != nil {
		log.Fatalf("Error adding a bot: %v", err)
	}
	err = server.AddBot(bots.NewAllOverTheNewsTomorrowBot(), bot.BotConfig{DatabaseRequired: true, SchedulerRequired: true})
	if err != nil {
		log.Fatalf("Error adding a bot: %v", err)
	}

	scheduler.StartAsync()

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

	DBClient mongo.Client
	AdminID int64
	Scheduler *gocron.Scheduler
}

func newServer() *server {
	return &server{ bq: map[string]BotQueue{} }
}

func (s *server) AddBot (b bot.Bot, config bot.BotConfig) error {
	name := b.GetName()

	token := os.Getenv(fmt.Sprintf("%s_token", name))
	if token == "" {
		return fmt.Errorf("Missing %s_token in env", name)
	}

	if config.DatabaseRequired {
		dbName := os.Getenv(fmt.Sprintf("%s_db", name))
		if dbName == "" {
			return fmt.Errorf("Missing %s_db", name)
		}

		b.SetDatabase(s.DBClient.Database(dbName))
	}

	if config.SchedulerRequired {
		b.SetScheduler(s.Scheduler)
		b.SetupScheduler()
	}

	b.SetToken(token)

	q := mq.NewMessageQueue(name)
	b.SetQueue(&q)

	s.bq[name] = BotQueue{ Bot: b, Queue: q}
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
