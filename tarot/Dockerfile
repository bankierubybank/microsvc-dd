## Build
FROM golang:1.20-alpine3.19 AS build

WORKDIR /app

ARG CIRunNumber
# ENV CIRunNumber=$CIRunNumber

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . ./

RUN echo $CIRunNumber

# RUN swag init -g ./tarot/main.go -o ./tarot/docs

RUN go build -o /microsvc-dd --ldflags="-X 'main.CIRunNumber=$CIRunNumber'"

## Deploy
FROM gcr.io/distroless/base-debian10:debug

WORKDIR /

COPY --from=build /microsvc-dd /microsvc-dd

USER nonroot:nonroot

ENV PORT=8080
EXPOSE ${PORT}

ENTRYPOINT ["/microsvc-dd"]