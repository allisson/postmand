#### development stage
FROM golang:1.20 AS builder

# set envvar
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE='on'

# set workdir
WORKDIR /code

# get project dependencies
COPY go.mod go.sum /code/
RUN go mod download

# copy files
COPY . /code

# generate binary
RUN go build -ldflags="-s -w" -o ./app ./cmd/postmand

#### final stage
FROM gcr.io/distroless/base:nonroot
COPY --from=builder /code/app /
COPY --from=builder /code/db/migrations /db/migrations
ENTRYPOINT ["/app"]
