package lambda

import (
	"encoding/json"

	"github.com/Golos1/faas_akt"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/go-playground/validator/v10"
	"github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/goaktpb"
)

// Actor which invokes a lambda function, taking parameters to the function as a JSON string and
// optionally returning the last 4 KB of the function logs as well. The consumer is responsible
// for initializing the lambdaClient and ensuring it has the permissions neccesary to invoke functions.
// The parameters are validated against the supplied type parameter.
type LambdaActor[T any] struct {
	functionName string
	getLogs      bool
	client       *lambda.Client
}

// Initializes a new Actor. getLogs being true only fetches the last 4 KB, not the whole logs of a function call.
func NewLambdaActor[T any](functionName string, client *lambda.Client, getLogs bool) LambdaActor[T] {
	return LambdaActor[T]{
		functionName: functionName,
		client:       client,
		getLogs:      getLogs,
	}
}

// This Actor doesn't have any setup for PreStart to do.
func (actor LambdaActor[T]) PreStart(ctx *actor.Context) error {
	return nil
}

// This Actor doesn't have any teardown for PostStop to do.
func (actor LambdaActor[T]) PostStop(ctx *actor.Context) error {
	return nil
}

// If the supplied parameters fail validation, the
// The Result of the invocation will be turned into a json string and sent back.
func (actor LambdaActor[T]) Receive(ctx *actor.ReceiveContext) {
	var args T
	switch ctx.Message().(type) {
	case *goaktpb.PostStart:
	case *faas_akt.Params:
		parameters := ctx.Message().ProtoReflect().Get(ctx.Message().ProtoReflect().Descriptor().Fields().ByName("JsonParamString")).String()
		arg_bytes := []byte(parameters)
		err := json.Unmarshal(arg_bytes, &args)
		if err != nil {
			ctx.Logger().Error("Error parsing parameters: ", err)
			ctx.Unhandled()
		}
		validator := validator.New()
		err = validator.Struct(args)
		if err != nil {
			ctx.Logger().Error("Error validating schema: ", err)
			ctx.Unhandled()
		}
		logType := types.LogTypeNone
		if actor.getLogs {
			logType = types.LogTypeTail
		}
		result, err := actor.client.Invoke(ctx.Context(), &lambda.InvokeInput{
			FunctionName: &actor.functionName,
			LogType:      logType,
			Payload:      arg_bytes,
		})
		if err != nil {
			ctx.Logger().Error("Error Invoking Function: ", err)
			ctx.Unhandled()

		} else {
			ctx.Logger().Info("Function successfully invoked.")
			reply := new(faas_akt.Result)
			reply.JsonResultString = string(result.Payload)
			reply.Logs = *result.LogResult
			ctx.Logger().Info(reply.JsonResultString)
			ctx.Response(reply)
		}
	default:
		ctx.Unhandled()
	}
}
