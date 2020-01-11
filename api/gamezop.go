package api

import (
	"context"
	"encoding/json"
	"gamezop/db/redis"
	"gamezop/kafka"
	"gamezop/protos"
	"net/http"

	"google.golang.org/grpc"
)

type Server struct {
}

func RegisterGamezopService(srv *grpc.Server) {
	gamezoppb.RegisterGamezopServiceServer(srv, &Server{})
}

// play game service
func (*Server) PlayGameService(ctx context.Context, in *gamezoppb.GamezopRequest) (*gamezoppb.GamezopResponse, error) {
	// write data into redis
	const prefix string = "player:"
	json, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	// set gamezop object in redis
	_, err = redis.GetClient().Do("SET", prefix+in.GetPlayerName(), json)
	if err != nil {
		return &gamezoppb.GamezopResponse{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong.",
		}, err
	}
	// send to messaging queue (any service) there is new data
	// kafka push data to consumer
	if err := kafka.Push(context.Background(), []byte("Gamezop"), []byte(prefix+in.GetPlayerName())); err != nil {
		return &gamezoppb.GamezopResponse{
			Code:    http.StatusInternalServerError,
			Message: "Kafka Can't push the Message.",
		}, err
	} else {
		return &gamezoppb.GamezopResponse{
			Code:    http.StatusOK,
			Message: "Success.",
			Data:    prefix + in.GetPlayerName(),
		}, nil
	}

}
