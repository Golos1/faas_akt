package inngest

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Golos1/faas_akt"
	"github.com/inngest/inngestgo"
	"github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/log"
)

type RandomJson struct {
	A int
	B int
	C string
}

func TestInngestActor(t *testing.T) {
	ctx := context.Background()
	logger := log.DefaultLogger
	actorSystem, _ := actor.NewActorSystem(
		"TestInngest",
		actor.WithLogger(logger),
	)
	err := actorSystem.Start(ctx)
	if err != nil {
		logger.Error(err)
		t.Error("Failed to start Actor System", err)
	}
	inngestKey := os.Getenv("INNGEST_KEY")
	inngestClient, err := inngestgo.NewClient(inngestgo.ClientOpts{EventKey: &inngestKey})
	if err != nil {
		logger.Error(err)
		t.Error("Error creating inngest client.", err)
	}
	inngestgo.CreateFunction(
		inngestClient, inngestgo.FunctionOpts{
			ID: "Spit-Back",
		},
		inngestgo.EventTrigger("echo", nil),
		func(ctx context.Context, input inngestgo.Input[map[string]any]) (any, error) {
			return input.Event.Data, nil
		},
	)
	pid, err := actorSystem.Spawn(ctx, "Inngest", NewInngestActor[map[string]string](inngestClient, "echo"))
	if err != nil {
		logger.Error(err)
		t.Error("Failed to spawn Actor", err)
	}
	jsonBytes, _ := json.Marshal(RandomJson{A: 1, B: 2, C: "stuff"})

	result, err := actor.Ask(ctx, pid, &faas_akt.InngestEvent{EventName: "echo", JsonParamString: string(jsonBytes)}, time.Second)
	if err != nil {
		logger.Error(err)
		t.Error("Failed to message Actor", err)
	}
	switch result.(type) {
	case *faas_akt.InvokedSuccessfully:
	default:
		t.Error("Failed to send event to inngest.")
	}
}
