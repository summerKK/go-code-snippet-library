package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/uber/jaeger-client-go"
)

func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx context.Context
		var span opentracing.Span
		// 尝试从header里面提取traceId
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span, ctx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(), global.Tracer, c.Request.URL.Path)
		} else {
			span, ctx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				global.Tracer,
				c.Request.URL.Path,
				opentracing.ChildOf(spanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		}
		defer span.Finish()

		var traceId string
		var spanId string
		var spanContext = span.Context()
		switch spanContext.(type) {
		case jaeger.SpanContext:
			traceId = spanContext.(jaeger.SpanContext).TraceID().String()
			spanId = spanContext.(jaeger.SpanContext).SpanID().String()
		}

		c.Set("X-Trace-ID", traceId)
		c.Set("X-Spa-ID", spanId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
