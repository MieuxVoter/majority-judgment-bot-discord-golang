FROM dhi.io/golang:1.26-debian13-dev AS builder
WORKDIR /build
COPY . /build
RUN make clean
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build

#FROM dhi.io/alpine-base:3.23 AS runner  # NOPE: sqlite-go requires CGO_ENABLED
FROM dhi.io/debian-base:trixie AS runner
WORKDIR /bot
COPY --from=builder /build/.env .
COPY --from=builder /build/.env.local .
COPY --from=builder /build/mjbot .
CMD ["/bot/mjbot"]
