package sdk

import (
	"context"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type ctxKey string

const correlationKey ctxKey = "correlationId"

// InitCorrelationID set jika kosong, get jika ada
func InitCorrelationID(ctx context.Context, incoming string) context.Context {
	if incoming != "" {
		return context.WithValue(ctx, correlationKey, incoming)
	}
	return context.WithValue(ctx, correlationKey, uuid.New().String())
}

// GetCorrelationID dari context
func GetCorrelationID(ctx context.Context) string {
	if v := ctx.Value(correlationKey); v != nil {
		return v.(string)
	}
	return ""
}

// Log message dengan correlationId
func Log(ctx context.Context, msg string) {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	log.Printf("[CID:%s][Function:%s] %s", GetCorrelationID(ctx), fn.Name(), msg)
}

// HttpClient with Resty
func HttpRequest(ctx context.Context, method, url string, body interface{}) (*resty.Response, error) {
	client := resty.New()
	cid := GetCorrelationID(ctx)

	req := client.R().
		SetHeader("X-Correlation-ID", cid).
		SetBody(body)

	var resp *resty.Response
	var err error

	switch method {
	case http.MethodGet:
		resp, err = req.Get(url)
	case http.MethodPost:
		resp, err = req.Post(url)
		// extend PUT, DELETE, etc
	}

	return resp, err
}

// HttpResponse with correlationId in header
func HttpResponse(ctx *gin.Context, httpStatus int, resp any) {
	cid := GetCorrelationID(ctx.Request.Context())
	ctx.Header("X-Correlation-ID", cid)
	ctx.JSON(httpStatus, resp)
}

// Middleware Gin untuk handle correlationId
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := c.GetHeader("X-Correlation-ID")
		ctx := InitCorrelationID(c.Request.Context(), cid)

		// update request dengan ctx baru
		c.Request = c.Request.WithContext(ctx)

		// tetap bisa set ke gin.Context kalau mau akses cepat
		c.Set(string(correlationKey), GetCorrelationID(ctx))

		// propagate correlationId ke response header
		c.Writer.Header().Set("X-Correlation-ID", GetCorrelationID(ctx))
		c.Next()
	}
}
