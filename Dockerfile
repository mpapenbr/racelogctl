FROM alpine:3.18
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ["/racelogctl"]
COPY racelogctl /
