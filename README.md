# postmand
[![Build Status](https://github.com/allisson/postmand/workflows/release/badge.svg)](https://github.com/allisson/postmand/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/postmand)](https://goreportcard.com/report/github.com/allisson/postmand)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/postmand)
[![Docker Image](https://img.shields.io/docker/cloud/build/allisson/postmand)](https://hub.docker.com/r/allisson/postmand)

Simple webhook delivery system powered by Golang and PostgreSQL.

## Features

- Simple rest api with only three endpoints (webhooks/deliveries/delivery-attempts).
- Select the status codes that are considered valid for a delivery.
- Control the maximum amount of delivery attempts and delay between these attempts (min and max backoff).
- Locks control of worker deliveries using PostgreSQL SELECT FOR UPDATE SKIP LOCKED.
- Sending the X-Hub-Signature header if the webhook is configured with a secret token.
- Simplicity, it does the minimum necessary, it will not have authentication/permission scheme among other things, the idea is to use it internally in the cloud and not leave exposed.

## Quickstart

Let's start with the basic concepts, we have three main entities that we must know to start:

- Webhook: The configuration of the webhook.
- Delivery: The content sent to a webhook.
- Delivery Attempt: An attempt to deliver the content to the webhook.

### Run the server

To run the server it is necessary to have a database available from postgresql, in this example we will consider that we have a database called postmand running in localhost with user and password equal to user.

#### Docker

```bash
docker run --rm --env POSTMAND_DATABASE_URL='postgres://user:password@host.docker.internal:5432/postmand?sslmode=disable' allisson/postmand migrate # create database schema
```

```bash
docker run -p 8000:8000 -p 8001:8001 --env POSTMAND_DATABASE_URL='postgres://user:password@host.docker.internal:5432/postmand?sslmode=disable' allisson/postmand server # run the server
```

#### Local

```bash
git clone https://github.com/allisson/postmand
cd postmand
cp local.env .env # and edit .env
make run-migrate # create database schema
make run-server # run the server
```

### Create a new webhook

The fields delivery_attempt_timeout/retry_min_backoff/retry_max_backoff are in seconds.

```bash
curl --location --request POST 'http://localhost:8000/v1/webhooks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Httpbin Post",
    "url": "https://httpbin.org/post",
    "content_type": "application/json",
    "valid_status_codes": [
        200,
        201
    ],
    "secret_token": "my-secret-token",
    "active": true,
    "max_delivery_attempts": 5,
    "delivery_attempt_timeout": 1,
    "retry_min_backoff": 10,
    "retry_max_backoff": 60
}'
```

```javascript
{
  "id":"a6e9a525-ac5a-488c-b118-bd7327ce6d8d",
  "name":"Httpbin Post",
  "url":"https://httpbin.org/post",
  "content_type":"application/json",
  "valid_status_codes":[
    200,
    201
  ],
  "secret_token":"my-secret-token",
  "active":true,
  "max_delivery_attempts":5,
  "delivery_attempt_timeout":1,
  "retry_min_backoff":10,
  "retry_max_backoff":60,
  "created_at":"2021-03-08T20:41:25.433671Z",
  "updated_at":"2021-03-08T20:41:25.433671Z"
}
```

### Create a new delivery

```bash
curl --location --request POST 'http://localhost:8000/v1/deliveries' \
--header 'Content-Type: application/json' \
--data-raw '{
    "webhook_id": "a6e9a525-ac5a-488c-b118-bd7327ce6d8d",
    "payload": "{\"success\": true}"
}'
```

```javascript
{
  "id":"bc76122c-e56b-45c7-8dc3-b80a861191d5",
  "webhook_id":"a6e9a525-ac5a-488c-b118-bd7327ce6d8d",
  "payload":"{\"success\": true}",
  "scheduled_at":"2021-03-08T20:43:49.986771Z",
  "delivery_attempts":0,
  "status":"pending",
  "created_at":"2021-03-08T20:43:49.986771Z",
  "updated_at":"2021-03-08T20:43:49.986771Z"
}
```

###  Run the worker

The worker is responsible to delivery content to the webhooks.

#### Docker

```bash
docker run --env POSTMAND_DATABASE_URL='postgres://user:pass@host.docker.internal:5432/postmand?sslmode=disable' allisson/postmand worker
{"level":"info","ts":1615236411.115703,"caller":"service/worker.go:74","msg":"worker-started"}
{"level":"info","ts":1615236411.1158803,"caller":"http/server.go:60","msg":"http-server-listen-and-server"}
{"level":"info","ts":1615236411.687701,"caller":"service/worker.go:42","msg":"worker-delivery-attempt-created","id":"d72719d6-5a79-4df7-a2c2-2029ab0e1848","webhook_id":"a6e9a525-ac5a-488c-b118-bd7327ce6d8d","delivery_id":"bc76122c-e56b-45c7-8dc3-b80a861191d5","response_status_code":200,"execution_duration":547,"success":true}
```

#### Local

```bash
make run-worker
go run cmd/postmand/main.go worker
{"level":"info","ts":1615236411.115703,"caller":"service/worker.go:74","msg":"worker-started"}
{"level":"info","ts":1615236411.1158803,"caller":"http/server.go:60","msg":"http-server-listen-and-server"}
{"level":"info","ts":1615236411.687701,"caller":"service/worker.go:42","msg":"worker-delivery-attempt-created","id":"d72719d6-5a79-4df7-a2c2-2029ab0e1848","webhook_id":"a6e9a525-ac5a-488c-b118-bd7327ce6d8d","delivery_id":"bc76122c-e56b-45c7-8dc3-b80a861191d5","response_status_code":200,"execution_duration":547,"success":true}
```

### Get deliveries

```bash
curl --location --request GET 'http://localhost:8000/v1/deliveries?webhook_id=a6e9a525-ac5a-488c-b118-bd7327ce6d8d'
```

```javascript
{
  "deliveries":[
    {
      "id":"bc76122c-e56b-45c7-8dc3-b80a861191d5",
      "webhook_id":"a6e9a525-ac5a-488c-b118-bd7327ce6d8d",
      "payload":"{\"success\": true}",
      "scheduled_at":"2021-03-08T20:43:49.986771Z",
      "delivery_attempts":1,
      "status":"succeeded",
      "created_at":"2021-03-08T20:43:49.986771Z",
      "updated_at":"2021-03-08T20:46:51.674623Z"
    }
  ],
  "limit":50,
  "offset":0
}
```

### Get delivery

```bash
curl --location --request GET 'http://localhost:8000/v1/deliveries/bc76122c-e56b-45c7-8dc3-b80a861191d5'
```

```javascript
{
  "id":"bc76122c-e56b-45c7-8dc3-b80a861191d5",
  "webhook_id":"a6e9a525-ac5a-488c-b118-bd7327ce6d8d",
  "payload":"{\"success\": true}",
  "scheduled_at":"2021-03-08T20:43:49.986771Z",
  "delivery_attempts":1,
  "status":"succeeded",
  "created_at":"2021-03-08T20:43:49.986771Z",
  "updated_at":"2021-03-08T20:46:51.674623Z"
}
```

### Get delivery attempts

```bash
curl --location --request GET 'http://localhost:8000/v1/delivery-attempts?delivery_id=bc76122c-e56b-45c7-8dc3-b80a861191d5'
```

```javascript
{
  "delivery_attempts":[
    {
      "id":"d72719d6-5a79-4df7-a2c2-2029ab0e1848",
      "webhook_id":"a6e9a525-ac5a-488c-b118-bd7327ce6d8d",
      "delivery_id":"bc76122c-e56b-45c7-8dc3-b80a861191d5",
      "raw_request":"POST /post HTTP/1.1\r\nHost: httpbin.org\r\nContent-Type: application/json\r\nX-Hub-Signature: 3fc5d4b8ff4efb404be24faf543667d29902d6a1306bd0c1ef2084497300cee9\r\n\r\n{\"success\": true}",
      "raw_response":"HTTP/2.0 200 OK\r\nContent-Length: 538\r\nAccess-Control-Allow-Credentials: true\r\nAccess-Control-Allow-Origin: *\r\nContent-Type: application/json\r\nDate: Mon, 08 Mar 2021 20:46:51 GMT\r\nServer: gunicorn/19.9.0\r\n\r\n{\n  \"args\": {}, \n  \"data\": \"{\\\"success\\\": true}\", \n  \"files\": {}, \n  \"form\": {}, \n  \"headers\": {\n    \"Accept-Encoding\": \"gzip\", \n    \"Content-Length\": \"17\", \n    \"Content-Type\": \"application/json\", \n    \"Host\": \"httpbin.org\", \n    \"User-Agent\": \"Go-http-client/2.0\", \n    \"X-Amzn-Trace-Id\": \"Root=1-60468d3b-36d312777a03ec3e1c564e3b\", \n    \"X-Hub-Signature\": \"3fc5d4b8ff4efb404be24faf543667d29902d6a1306bd0c1ef2084497300cee9\"\n  }, \n  \"json\": {\n    \"success\": true\n  }, \n  \"origin\": \"191.35.122.74\", \n  \"url\": \"https://httpbin.org/post\"\n}\n",
      "response_status_code":200,
      "execution_duration":547,
      "success":true,
      "error":"",
      "created_at":"2021-03-08T20:46:51.680846Z"
    }
  ],
  "limit":50,
  "offset":0
}
```

### Get delivery attempt

```bash
curl --location --request GET 'http://localhost:8000/v1/delivery-attempts/d72719d6-5a79-4df7-a2c2-2029ab0e1848'
```

```javascript
{
  "id":"d72719d6-5a79-4df7-a2c2-2029ab0e1848",
  "webhook_id":"a6e9a525-ac5a-488c-b118-bd7327ce6d8d",
  "delivery_id":"bc76122c-e56b-45c7-8dc3-b80a861191d5",
  "raw_request":"POST /post HTTP/1.1\r\nHost: httpbin.org\r\nContent-Type: application/json\r\nX-Hub-Signature: 3fc5d4b8ff4efb404be24faf543667d29902d6a1306bd0c1ef2084497300cee9\r\n\r\n{\"success\": true}",
  "raw_response":"HTTP/2.0 200 OK\r\nContent-Length: 538\r\nAccess-Control-Allow-Credentials: true\r\nAccess-Control-Allow-Origin: *\r\nContent-Type: application/json\r\nDate: Mon, 08 Mar 2021 20:46:51 GMT\r\nServer: gunicorn/19.9.0\r\n\r\n{\n  \"args\": {}, \n  \"data\": \"{\\\"success\\\": true}\", \n  \"files\": {}, \n  \"form\": {}, \n  \"headers\": {\n    \"Accept-Encoding\": \"gzip\", \n    \"Content-Length\": \"17\", \n    \"Content-Type\": \"application/json\", \n    \"Host\": \"httpbin.org\", \n    \"User-Agent\": \"Go-http-client/2.0\", \n    \"X-Amzn-Trace-Id\": \"Root=1-60468d3b-36d312777a03ec3e1c564e3b\", \n    \"X-Hub-Signature\": \"3fc5d4b8ff4efb404be24faf543667d29902d6a1306bd0c1ef2084497300cee9\"\n  }, \n  \"json\": {\n    \"success\": true\n  }, \n  \"origin\": \"191.35.122.74\", \n  \"url\": \"https://httpbin.org/post\"\n}\n",
  "response_status_code":200,
  "execution_duration":547,
  "success":true,
  "error":"",
  "created_at":"2021-03-08T20:46:51.680846Z"
}
```

### Health check

The health check server is running on port defined by envvar POSTMAND_HEALTH_CHECK_HTTP_PORT (defaults to 8001).

```bash
curl --location --request GET 'http://localhost:8001/healthz'
```

```javascript
{
  "success":true
}
```

### Environment variables

All environment variables is defined on file local.env.

## How to build docker image

```
docker build -f Dockerfile -t postmand .
```

