FROM golang:1-alpine AS builder
ENV CGO_ENABLED 0

RUN mkdir /build
WORKDIR /src

COPY . .

RUN go build -trimpath -o /build/web-server ./main.go

FROM gcr.io/distroless/static:nonroot AS release
LABEL maintainer "ghiloni@gmail.com"

WORKDIR /app

COPY --from=builder /build/web-server /app/web-server
EXPOSE 3000

ENTRYPOINT ["/app/web-server"]
