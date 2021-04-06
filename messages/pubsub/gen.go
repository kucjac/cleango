package pubsub

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -destination=./../mockbus/bus.go -package=mockbus . Publisher,Subscriber
