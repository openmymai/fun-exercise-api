FROM golang:1.22.1-alpine as build-base

ARG DATABASE_URL

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -tags=unit ./...

RUN go build -o ./out/funx .


### ------------

FROM alpine:3.19
COPY --from=build-base /app/out/funx /app/funx

ENV DATABASE_URL=postgres://root:password@localhost:5432/wallet?sslmode=disable

CMD ["/app/funx"]