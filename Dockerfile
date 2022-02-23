FROM golang:1.17.1-buster as build

ADD ./github-action /github-action
WORKDIR /github-action
RUN go get && go build -o /ga.bin main.go

FROM debian:buster-slim as run

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
COPY --from=build /ga.bin /usr/bin/ga
ENTRYPOINT ["ga"]
