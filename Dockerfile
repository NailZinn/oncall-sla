FROM golang:1.23.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./metrics ./metrics
COPY ./prober ./prober
COPY ./sla ./sla

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /oncall-sla

ENTRYPOINT ["/oncall-sla"]