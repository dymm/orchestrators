package tracer

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

//NewTracer create a new tracer
func NewTracer() (opentracing.Tracer, io.Closer, error) {

	cfg, err := jaegercfg.FromEnv()
	cfg.Sampler = &jaegercfg.SamplerConfig{
		Type:  jaeger.SamplerTypeConst,
		Param: 1,
	}
	var tracer opentracing.Tracer
	var closer io.Closer
	if err == nil {
		tracer, closer, err = cfg.NewTracer()
	}
	if err == nil {
		opentracing.SetGlobalTracer(tracer)
	} else {
		err = errors.Wrap(err, "Could not initialize jaeger tracer")
	}
	return tracer, closer, nil

}
