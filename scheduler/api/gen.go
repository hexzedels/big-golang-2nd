package api

//go:generate go tool oapi-codegen -o ../internal/input/http/gen/types.go -generate types -package gen openapi.yaml
//go:generate go tool oapi-codegen  -o ../internal/input/http/gen/server.go -generate chi-server,strict-server -package gen openapi.yaml
//go:generate go tool oapi-codegen  -o ../internal/input/http/gen/spec.go -generate spec -package gen openapi.yaml

//go:generate go tool oapi-codegen -o ../pkg/client/http/types.go -generate types -package client openapi.yaml
//go:generate go tool oapi-codegen -o ../pkg/client/http/http_client.go -generate client -package client openapi.yaml
