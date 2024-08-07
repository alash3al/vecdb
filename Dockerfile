FROM golang:1.22-alpine As builder

WORKDIR /vecdb/

RUN apk update && apk add git upx

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /usr/bin/vecdb ./cmd/

RUN upx -9 /usr/bin/vecdb

FROM alpine

WORKDIR /vecdb/

COPY --from=builder /usr/bin/vecdb /usr/bin/vecdb