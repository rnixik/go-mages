# build stage
FROM golang:1.22 AS builder

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mages ./cmd/go-mages.go

# final stage
FROM busybox
WORKDIR /app
COPY --from=builder /app/mages /app/
COPY web /app/web
COPY version /app/version
RUN sed -i -e 's#/ws#wss://mages.getid.org/ws#g' /app/web/index.html
EXPOSE 8009
ENTRYPOINT ["./mages", "--addr=:8009"]
