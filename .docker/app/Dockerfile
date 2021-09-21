# build stage
FROM golang:alpine AS build-env
WORKDIR /build
ADD . ./
RUN go build -mod vendor -o app ./cmd/app

# final stage
FROM alpine
WORKDIR /opt/app
COPY --from=build-env /build/app ./
RUN chmod +x ./app
ENTRYPOINT ./app