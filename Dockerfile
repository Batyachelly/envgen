FROM golang:1.20.3-alpine3.17 as builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o envgen cmd/envgen/main.go


FROM scratch

WORKDIR /root/

COPY --from=builder /app/envgen .