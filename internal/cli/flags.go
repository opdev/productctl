package cli

// FlagID is the string representation of a CLI flag. That is, the long flag
// value bound to the CLI. E.g. "log-level" for --log-level
type FlagID = string

// These flagIDs are used in productctl's core functionality of manipulating
// product-listings
const (
	FlagIDEndpoint                FlagID = "env"                             // For choosing GraphQL endpoints based on env labels.
	FlagIDLogLevel                FlagID = "log-level"                       // For specifying log verbosity.
	FlagIDVersionAsJSON           FlagID = "json"                            // For printing version output as JSON.
	FlagIDCustomEndpoint          FlagID = "custom-endpoint"                 // For defining a GraphQL endpoint that isn't predefined.
	FlagIDCreateBackupOnOverwrite FlagID = "backup-declaration-on-overwrite" // For creating declaration backups before overwriting
)

// These flagIDs are mostly used in cert-tool-runner applications.
const (
	FlagIDUserfilesDir    FlagID = "userfiles-dir"
	FlagIDLogsDir         FlagID = "logs-dir"
	FlagIDCatalogAPIToken FlagID = "catalog-api-token"
	FlagIDRuntimeImage    FlagID = "runtime-image"
	FlagIDKeepTempDir     FlagID = "keep-temp-dir"
)
