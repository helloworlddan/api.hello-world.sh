FROM golang:1.22-bookworm as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -v -o server

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY ./static /app/static
CMD ["/app/server"]
