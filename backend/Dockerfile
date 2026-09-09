FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/backend ./cmd/app

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /data/attachments \
    && chown -R app:app /data

COPY --from=build /out/backend /usr/local/bin/backend

USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/backend"]
