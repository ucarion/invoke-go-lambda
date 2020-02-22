# invoke-go-lambda

`invoke-go-lambda` is a CLI tool that helps you invoke locally-running Golang
AWS Lambda functions.

## Installation

Install `invoke-go-lambda` by running:

```bash
go install github.com/ucarion/invoke-go-lambda
```

## Usage

### Simple lambdas

To run it, first invoke your lambda and force it to run on a particular port of
your choosing using the `_LAMBDA_SERVER_PORT` env var:

```bash
# As an example, we'll use the "hello" lambda in this repo.
_LAMBDA_SERVER_PORT=8001 go run ./examples/hello/...
```

You can then invoke this lambda by running:

```bash
echo '{}' | invoke-go-lambda --port=8001
```

```json
{"Payload":"ImhlbGxvLCB3b3JsZCEi","Error":null}
```

That's the raw output. If you'd like to parse that payload, you can pipe into
`jq` and `base64` as so:

```bash
echo '{}' | invoke-go-lambda --port=8001 | jq -r .Payload | base64 -D
```

```text
"hello, world!"
```

### Lambdas requiring a payload

The `hello` lambda we run here was the extremely simple case, where the lambda
doesn't try to read the "payload" parameter. Let's try a more complicated
lambda:

```bash
# The "ping-pong" lambda in this repo reads a JSON payload.
#
# If you give a JSON with ping=true, it'll return a JSON with pong=true.
_LAMBDA_SERVER_PORT=8001 go run ./examples/ping-pong/...
```

If we try to call it like we did before, we get an error:

```bash
echo '{}' | invoke-go-lambda --port=8001
```

```json
{"Payload":null,"Error":{"Message":"unexpected end of JSON input","Type":"SyntaxError","StackTrace":null,"ShouldExit":false}}
```

This happens because the underlying Lambda protocol expects a `payload` field
whose value is an array of bytes, and those bytes have to encode for valid JSON.

Let's encode the simplest possible JSON -- `{}` -- as a two-byte array: (`123`
and `125` encode for `{` and `}`, respectively)

```bash
echo '{"payload":[123, 125]}' | invoke-go-lambda --port=8001 | jq -r .Payload | base64 -D
```

```json
{"ping":true,"pong":false}
```

This time, it worked. But this is a pain to use. For the common case where you
only care about the `payload` and nothing else, use `--stdin-is-payload`:

```bash
echo '{"ping": true}' | invoke-go-lambda --stdin-is-payload --port=8001 | jq -r .Payload | base64 -D
```

```json
{"ping":false,"pong":true}
```

The `--stdin-is-payload` parameter reads standard input, and uses that as the
payload.
