package pubsub

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -destination=./../mockps/bus.go -package=mockps . Publisher,Subscriber
