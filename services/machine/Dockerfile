FROM golang:1.21-bookworm as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -v -o server # TODO disable CGO

FROM gcr.io/distroless/static
COPY --from=builder /app/server /app/server

CMD ["/app/server"]
