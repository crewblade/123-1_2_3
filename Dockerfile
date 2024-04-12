

FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN go build -tags migrate -o /bin/app ./cmd/app/main.go

FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
CMD ["/app"]