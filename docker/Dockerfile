FROM dhi.io/golang:1.26-debian13-dev AS builder
WORKDIR /build
COPY . /build
RUN make clean
RUN --mount=type=cache,sharing=locked,target=/go/pkg/mod \
    --mount=type=cache,sharing=locked,target=/root/.cache/go-build \
    make build

#FROM dhi.io/alpine-base:3.23 AS runner  # NOPE: sqlite-go requires CGO_ENABLED
#FROM dhi.io/debian-base:trixie AS runner  # NOPE: we need to install chromium
FROM dhi.io/debian-base:trixie-dev AS runner
WORKDIR /bot
# We need chromium to convert SVG to PNG (inkscape does not cut it)
RUN --mount=type=cache,sharing=locked,target=/var/lib/apt/lists \
    --mount=type=cache,sharing=locked,target=/var/cache/apt \
    ls -FlAh /etc/apt/apt.conf.d/ \
    && rm -f /etc/apt/apt.conf.d/docker-clean \
    && apt-get update \
    && apt-get --yes --no-install-recommends install \
        chromium
# We cannot copy chromium from that image, the build complains about the CPU lacking sse3
#COPY --from=linuxserver/chromium:latest /usr/bin/chromium /usr/bin/chromium
COPY --from=builder /build/.env .
COPY --from=builder /build/.env.local .
COPY --from=builder /build/mjbot .
CMD ["/bot/mjbot"]
