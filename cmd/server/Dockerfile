FROM golang:1.14 as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd cmd
COPY openapi openapi
COPY pkg pkg

RUN GOOS=linux CGO_ENABLED=0 go build -o /api ./cmd/server

FROM scratch

WORKDIR /app

COPY --from=build /api .

ENTRYPOINT [ "/app/api" ]
