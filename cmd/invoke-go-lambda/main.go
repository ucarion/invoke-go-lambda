package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/rpc"
	"os"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

var port = flag.Int("port", 8001, "Port that the lambda is listening on. Should be equal to the _LAMBDA_SERVER_PORT env var of the Golang lambda running.")
var stdinIsPayload = flag.Bool("stdin-is-payload", false, "Use STDIN as the payload of the request. All other invoke parameters will be left empty")

func main() {
	flag.Parse()

	client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		panic(err)
	}

	args := messages.InvokeRequest{}
	if *stdinIsPayload {
		buf := bytes.Buffer{}
		if _, err := io.Copy(&buf, os.Stdin); err != nil {
			panic(err)
		}

		args.Payload = buf.Bytes()
	} else {
		// Assume stdin contains a JSON-encoded InvokeRequest.
		if err := json.NewDecoder(os.Stdin).Decode(&args); err != nil {
			panic(err)
		}
	}

	var reply messages.InvokeResponse
	if err := client.Call("Function.Invoke", args, &reply); err != nil {
		panic(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(reply); err != nil {
		panic(err)
	}
}
