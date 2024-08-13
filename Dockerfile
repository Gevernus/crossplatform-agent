ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}
WORKDIR /app
COPY . .
#RUN apk add --no-cache build-base
#RUN apk add --no-cache pkgconfig
#RUN apk add --no-cache libappindicator3-dev pkgconfig gtk+3.0-dev
RUN apt-get update && apt-get install -y \
    ayatana-indicator-application \
    libayatana-appindicator3-dev \
    pkg-config \
    libgtk-3-dev \
    && rm -rf /var/lib/apt/lists/*
    
RUN go mod download
RUN CC="gcc" CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o /bin/crossplatform-agent.exe ./cmd/agent
RUN CC="gcc" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/crossplatform-agent-linux ./cmd/agent