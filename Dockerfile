ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN CC="gcc" CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o /bin/crossplatform-agent-macos ./cmd/agent
RUN CC="gcc" CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o /bin/crossplatform-agent.exe ./cmd/agent
RUN CC="gcc" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/crossplatform-agent-linux ./cmd/agent