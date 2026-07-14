package client

//go:generate go tool github.com/Khan/genqlient genqlient.yaml
//go:generate go run ./internal/gql/exportgen.go -- ./internal/gql/generated.go ./internal/gql/exports.go . github.com/KyaniteHQ/linctl/internal/client/internal/gql
