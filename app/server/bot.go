package main

import (
	pb "github.com/ar2rworld/botsservice/app/messageservice"
)

type Bot interface {
	HandleUpdate(*pb.Update) (pb.MessageReply, error)
	GetName() string
	GetToken() string
}