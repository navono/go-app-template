package logging

import (
	"fmt"
	"net/http"
	"time"

	"go-app-template/internal/dependency"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//// NewResponseLogger returns you an instance of a *ResponseLogger
//func NewResponseLogger(w http.ResponseWriter) *ResponseLogger {
//	return &ResponseLogger{
//		ResponseWriter: w,
//		Status:         http.StatusOK,
//	}
//}

// ResponseLogger is a ResponseWriter that is able to log the
// status code of the response
type ResponseLogger struct {
	http.ResponseWriter
	Status int
}

// WriteHeader intercepts the call to the base ResponseWriter and logs the
// status code sent
func (l *ResponseLogger) WriteHeader(statusCode int) {
	l.Status = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

// NewMiddleware returns you a new instance of the Logger middleware
func NewMidlleware(logger *zap.Logger, config dependency.ConfigGetter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			fields := []zapcore.Field{
				zap.Int("status", res.Status),
				zap.String("latency", time.Since(start).String()),
				zap.String("id", id),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("host", req.Host),
				zap.String("remote_ip", c.RealIP()),
			}

			fields = append(fields, requestHeaders(req, config.GetStringSlice("excluded-headers"))...)
			fields = append(fields, queryParams(req)...)

			n := res.Status
			switch {
			case n >= 500:
				logger.Error("Server error", fields...)
			case n >= 400:
				logger.Warn("Client error", fields...)
			case n >= 300:
				logger.Info("Redirection", fields...)
			default:
				logger.Info("Success", fields...)
			}

			return nil
		}
	}
	//return func(handler http.Handler) http.Handler {
	//	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	//		fields := []zap.Field{
	//			zap.String("method", r.Method),
	//			zap.String("host", r.Host),
	//			zap.String("path", path(r)),
	//			zap.String("protocol", r.Proto),
	//			zap.Int64("request.content-length", r.ContentLength),
	//		}
	//		fields = append(fields, requestHeaders(r, config.GetStringSlice("excluded-headers"))...)
	//		fields = append(fields, queryParams(r)...)
	//		responseLogger := NewResponseLogger(rw)
	//		handler.ServeHTTP(responseLogger, r)
	//		fields = append(fields, zap.Int("status-code", responseLogger.Status))
	//		fields = append(fields, responseHeaders(responseLogger.Header())...)
	//		logger.Info("request log", fields...)
	//	})
	//}
}

//
//func path(r *http.Request) string {
//	path := r.URL.Path
//	muxPath, err := mux.
//		CurrentRoute(r).
//		GetPathTemplate()
//	if err != nil {
//		return path
//	}
//	return muxPath
//}
//
//func responseHeaders(headers http.Header) []zap.Field {
//	fields := make([]zap.Field, 0, len(headers))
//	for header := range headers {
//		headerName := fmt.Sprintf("response.header.%s", header)
//		field := zap.String(headerName, headers.Get(header))
//		fields = append(fields, field)
//	}
//	return fields
//}

func requestHeaders(r *http.Request, excludedHeaders []string) []zap.Field {
	headers := map[string]struct{}{}
	for _, header := range excludedHeaders {
		headers[header] = struct{}{}
	}
	fields := make([]zap.Field, 0, len(r.Header))
	for header := range r.Header {
		_, ok := headers[header]
		if ok {
			continue
		}
		headerName := fmt.Sprintf("request.header.%s", header)
		field := zap.String(headerName, r.Header.Get(header))
		fields = append(fields, field)
	}
	return fields
}

func queryParams(r *http.Request) []zap.Field {
	fields := make([]zap.Field, 0, len(r.Header))
	for param, values := range r.URL.Query() {
		queryName := fmt.Sprintf("query.%s", param)
		field := zap.Strings(queryName, values)
		fields = append(fields, field)
	}
	return fields
}
