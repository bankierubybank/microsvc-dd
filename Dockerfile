## Build
FROM golang:1.20.14-alpine3.19 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN go build -o /microsvc-dd

## Deploy
FROM gcr.io/distroless/base-debian10:debug

WORKDIR /

COPY --from=build /microsvc-dd /microsvc-dd

USER nonroot:nonroot

ENV PORT=8080
EXPOSE ${PORT}

ENTRYPOINT ["/microsvc-dd"]