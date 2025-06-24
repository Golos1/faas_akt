package lambda

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Golos1/faas_akt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/log"
)

type AddParams struct {
	A int `validate:"required"`
	B int `validate:"required"`
}
type AddResult struct {
	Result int `json:"result"`
}

func TestLambdaActor(t *testing.T) {

	ctx := context.Background()
	logger := log.DefaultLogger
	actorSystem, _ := actor.NewActorSystem(
		"TestLambda",
		actor.WithLogger(logger),
	)
	err := actorSystem.Start(ctx)
	if err != nil {
		logger.Error(err)
		t.Error("Failed to start Actor System", err)
	}
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Error(err)
		t.Error("Failed to create Lambda client", err)
	}
	client := lambda.NewFromConfig(sdkConfig)
	pid, err := actorSystem.Spawn(ctx, "Lambda", NewLambdaActor[AddParams]("Add", client, true))
	if err != nil {
		logger.Error(err)
		t.Error("Failed to spawn actor", err)
	}
	params := new(faas_akt.Params)
	numberBytes, err := json.Marshal(AddParams{A: 2, B: 3})
	if err != nil {
		logger.Error(err)
		t.Error("Failed to spawn actor", err)
	}
	params.JsonParamString = string(numberBytes)

	response, err := actor.Ask(ctx, pid, params, time.Second)
	if err != nil {
		logger.Error(err)
		logger.Info(response)
		t.Error("Failed to message actor", err)
	}
	switch response.(type) {
	case *faas_akt.Result:
		descriptor := response.ProtoReflect().Descriptor().Fields().ByName("JsonResultString")
		jsonResult := response.ProtoReflect().Get(descriptor).String()
		structResult := new(AddResult)
		json.Unmarshal([]byte(jsonResult), structResult)
		if structResult.Result != 5 {
			logger.Info(jsonResult)
			logger.Info(structResult)
			t.Error("Add should have returned 5", structResult.Result)
		}
		descriptor = response.ProtoReflect().Descriptor().Fields().ByName("Logs")
		logs := response.ProtoReflect().Get(descriptor).String()
		if logs == "" {
			t.Error("Empty Logs")
		}
	}
}
