package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gamezop/db/mongodb"
	"gamezop/db/redis"
	"gamezop/kafka"
	"gamezop/protos"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

type GamezopPlayers struct {
	Id         string `json:"user_id,omitempty"`
	PlayerName string `json:"player_name,omitempty"`
	Email      string `json:"email,omitempty"`
}

// get all players details
func (*Server) GetPlayersDetailsService(ctx context.Context, in *gamezoppb.Empty) (*gamezoppb.GamezopResponse, error) {
	//get all players info from mongodb
	res, err := mongodb.CreateCollection("all_players").Find(context.Background(), nil)
	if err != nil {
		return &gamezoppb.GamezopResponse{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong.",
		}, err
	} else {
		var items []GamezopPlayers
		for res.Next(nil) {
			item := GamezopPlayers{}
			if err := res.Decode(&item); err != nil {
				return nil, status.Errorf(
					codes.Aborted,
					fmt.Sprintln("Data can't decode.", err))
			}
			//items = append(items, &item)
			items = append(items, item)
		}
		data, _ := json.Marshal(items)
		return &gamezoppb.GamezopResponse{
			Code:    http.StatusOK,
			Message: "Success.",
			Data:    string(data),
		}, nil
	}
}
