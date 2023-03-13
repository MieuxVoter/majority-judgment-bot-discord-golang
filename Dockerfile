FROM andersfylling/disgord:latest as builder
MAINTAINER https://github.com/andersfylling
WORKDIR /build
COPY . /build
#RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o discordmjbot src/main.go
RUN GOOS=linux go build -a -o discordmjbot src/main.go

FROM gcr.io/distroless/base
WORKDIR /bot
COPY --from=builder /build/.env .
COPY --from=builder /build/.env.local .
COPY --from=builder /build/discordmjbot .
CMD ["/bot/discordmjbot"]
