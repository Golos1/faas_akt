package inngest

import (
	"encoding/json"

	"github.com/Golos1/faas_akt"
	"github.com/inngest/inngestgo"
	"github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/goaktpb"
)

func NewInngestActor[T any](client inngestgo.Client, eventName string) *InngestActor[T] {
	return &InngestActor[T]{
		InngestClient: client,
	}
}

type InngestActor[T any] struct {
	InngestClient inngestgo.Client
}

// PreStart implements actor.Actor. No setup is needed for this Actor
func (actor *InngestActor[T]) PreStart(ctx *actor.Context) error {
	ctx.ActorSystem().Logger().Info("PreStart for Inngest Actor")
	return nil
}

// Receive implements actor.Actor.
func (actor *InngestActor[T]) Receive(ctx *actor.ReceiveContext) {
	switch ctx.Message().(type) {
	case *goaktpb.PostStart:
	case *faas_akt.InngestEvent:
		ctx.Logger().Info("Received Message.")
		eventName := ctx.Message().ProtoReflect().Get(ctx.Message().ProtoReflect().Descriptor().Fields().ByName("EventName")).String()
		payload := ctx.Message().ProtoReflect().Get(ctx.Message().ProtoReflect().Descriptor().Fields().ByName("JsonParamString")).String()
		var payloadMap map[string]any
		err := json.Unmarshal([]byte(payload), &payloadMap)
		if err != nil {
			ctx.Logger().Error("Error parsing parameters: ", err)
			ctx.Unhandled()
		}
		ctx.Logger().Info("Parsed Message.")
		_, err = actor.InngestClient.Send(ctx.Context(), inngestgo.Event{
			Name: eventName,
			Data: payloadMap,
		})
		ctx.Logger().Info("Got OK from inngest.")
		if err != nil {
			ctx.Logger().Error("Error Sending event: ", err)
			ctx.Unhandled()
		} else {
			ctx.Response(&faas_akt.InvokedSuccessfully{})
			ctx.Logger().Info("Invoked Inngest successfully.")
		}
	default:
		ctx.Logger().Error("Wrong Type of message")
		ctx.Unhandled()
	}
}

// PostStop implements actor.Actor. No teardown is needed for this Actor.
func (actor *InngestActor[T]) PostStop(ctx *actor.Context) error {
	return nil
}
