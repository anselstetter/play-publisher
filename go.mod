module github.com/anselstetter/play-publisher

go 1.24.0

require (
	github.com/shogo82148/androidbinary v1.0.5
	github.com/xmxu/aab-parser v0.1.0
	google.golang.org/api v0.218.0
)

require (
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	github.com/client9/misspell v0.3.4 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/fzipp/gocyclo v0.6.0 // indirect
	github.com/gordonklaus/ineffassign v0.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kisielk/errcheck v1.8.0 // indirect
	github.com/nishanths/exhaustive v0.12.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/tools v0.6.0 // indirect
)

require (
	cloud.google.com/go/auth v0.14.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.7 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/spf13/cobra v1.8.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel v1.31.0 // indirect
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.31.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/image v0.7.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/oauth2 v0.25.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc v1.69.4 // indirect
	google.golang.org/protobuf v1.36.3 // indirect
)

// Accidents happen
retract v1.0.0

tool (
	github.com/client9/misspell/cmd/misspell
	github.com/fzipp/gocyclo/cmd/gocyclo
	github.com/gordonklaus/ineffassign
	github.com/kisielk/errcheck
	github.com/nishanths/exhaustive/cmd/exhaustive
	honnef.co/go/tools/cmd/staticcheck
)
