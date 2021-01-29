package tracer

import (
	"context"
	"fmt"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
)

type TracerConfig struct {
	Service   string
	Collector string
	AgentHost string
	AgentPort string
	Metadata  metadata.Metadata
}

func NewTracer(cfg *TracerConfig) (opentracing.Tracer, io.Closer, error) {
	jcfg := &jconfig.Configuration{
		ServiceName: cfg.Service,
		Sampler: &jconfig.SamplerConfig{
			Type:  "remote",
			Param: 1,
		},
		Reporter: &jconfig.ReporterConfig{
			LogSpans:  false,
			QueueSize: 100,
		},
	}

	if len(cfg.AgentHost) > 0 {
		jcfg.Sampler.SamplingServerURL = fmt.Sprintf("http://%s:5778/sampling", cfg.AgentHost)
	}

	if len(cfg.AgentHost) > 0 {
		agentPort := "6831"
		if len(cfg.AgentPort) > 0 {
			agentPort = cfg.AgentPort
		}
		jcfg.Reporter.LocalAgentHostPort = fmt.Sprintf("%s:%s", cfg.AgentHost, agentPort)
	} else if len(cfg.Collector) > 0 {
		jcfg.Reporter.CollectorEndpoint = cfg.Collector
	}

	for k, v := range cfg.Metadata {
		jcfg.Tags = append(jcfg.Tags, opentracing.Tag{Key: k, Value: v})
	}
	tracer, closer, err := jcfg.NewTracer(jconfig.Logger(&log{}))
	if err != nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}

type log struct{}

func (l *log) Error(msg string) {
	logger.Error(context.Background(), msg)
}

func (l *log) Infof(msg string, args ...interface{}) {
	logger.Errorf(context.Background(), msg, args...)
}

func (l *log) Debugf(msg string, args ...interface{}) {
	logger.Debugf(context.Background(), msg, args...)
}
