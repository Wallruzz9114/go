# BUILD environment
# -----------------
FROM golang:1.12-alpine as build-env
WORKDIR /bookey

RUN apk update && apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/app ./cmd/app \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate


# Deployment environment
# ----------------------
FROM alpine
RUN apk update && apk add --no-cache bash

COPY --from=build-env /bookey/bin/app /bookey/
COPY --from=build-env /bookey/bin/migrate /bookey/
COPY --from=build-env /bookey/migrations /bookey/migrations

EXPOSE 8080
CMD ["/bookey/app"]