package oc

import (
	// stdlib
	"io"

	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	// external
	zipkin "github.com/openzipkin/zipkin-go"
	reporter "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/shared/network"
)

const (
	zipkinURL = "http://localhost:9411/api/v2/spans"
)

func setupZipkin(serviceName string) (trace.Exporter, io.Closer) {
	var (
		rep     = reporter.NewReporter(zipkinURL)
		addr, _ = network.HostIP()
	)
	localEndpoint, _ := zipkin.NewEndpoint(serviceName, addr)

	return oczipkin.NewExporter(rep, localEndpoint), rep
}
