module github.com/temporalio/samples-go

go 1.16

replace go.temporal.io/sdk => /Users/rob/Code/github.com/robholland/sdk-go

replace go.temporal.io/server => /Users/rob/Code/github.com/robholland/temporal

require (
	github.com/HdrHistogram/hdrhistogram-go v0.9.0 // indirect
	github.com/golang/mock v1.5.0
	github.com/hashicorp/go-plugin v1.4.0
	github.com/m3db/prometheus_client_golang v0.8.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/uber-go/tally v3.3.17+incompatible
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	go.temporal.io/api v1.4.1-0.20210319015452-3dc250bb642a
	go.temporal.io/sdk v1.5.1-0.20210318225734-a39bbe82b2ba
	go.temporal.io/server v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
