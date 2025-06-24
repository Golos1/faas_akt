package faas_akt

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/log"
)

type AddParams struct {
	a int `validate:"required"`
	b int `validate:"required"`
}
type AddResult struct {
	result int `validate:"required"`
}

func TestLambdaActor(t *testing.T) {

	ctx := context.Background()
	logger := log.DefaultLogger
	actorSystem, _ := actor.NewActorSystem(
		"TestLambda",
	)
	if err := actorSystem.Start(ctx); err != nil {
		logger.Error(err)
		t.Error("Failed to start Actor System", err)
	}
	client := lambda.Client{}
	pid, err := actorSystem.Spawn(ctx, "Lambda", NewLambdaActor[AddParams]("Add", &client, true))
	if err != nil {
		logger.Error(err)
		t.Error("Failed to spawn actor", err)
	}
	params := new(Params)
	numberBytes, err := json.Marshal(AddParams{a: 2, b: 3})
	if err != nil {
		logger.Error(err)
		t.Error("Failed to spawn actor", err)
	}
	params.JsonParamString = string(numberBytes)
	response, err := actor.Ask(ctx, pid, params, time.Minute)
	if err != nil {
		logger.Error(err)
		logger.Info(response)
		t.Error("Failed to message actor", err)
	}
	switch response.(type) {
	case *Result:
		descriptor := response.ProtoReflect().Descriptor().Fields().ByName("JsonParamString")
		jsonResult := response.ProtoReflect().Get(descriptor)
		var structResult AddResult
		json.Unmarshal(jsonResult.Bytes(), &structResult)
		if structResult.result != 5 {
			t.Error("Add should have returned 5", structResult.result, 5)
		}
	}
}
