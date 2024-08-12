ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=windows GOARCH=amd64 go build -o /bin/crossplatform-agent.exe ./cmd/agent
RUN GOOS=linux GOARCH=amd64 go build -o crossplatform-agent ./cmd/agent