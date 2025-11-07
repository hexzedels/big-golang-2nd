FROM golang:1.24.2-alpine AS build
WORKDIR /app

# Pre-cache dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod apk add --no-cache git && go mod download

# Copy source
COPY . .

# Build scheduler binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/scheduler ./scheduler/cmd

# Build worker binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/worker ./worker/cmd

FROM gcr.io/distroless/static:nonroot AS scheduler
WORKDIR /
COPY --from=build /out/scheduler /scheduler
EXPOSE 8090
USER nonroot:nonroot
ENTRYPOINT ["/scheduler"]

FROM gcr.io/distroless/static:nonroot AS worker
WORKDIR /
COPY --from=build /out/worker /worker
USER nonroot:nonroot
ENTRYPOINT ["/worker"]

