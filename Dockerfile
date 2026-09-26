FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/gopython ./cmd/gopython

FROM scratch

WORKDIR /app

COPY --from=builder /out/gopython /gopython
COPY examples ./examples

ENTRYPOINT ["/gopython"]
CMD ["examples/functions.py"]
