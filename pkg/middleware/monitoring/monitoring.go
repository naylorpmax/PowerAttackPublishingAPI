package monitoring

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/naylorpmax/homebrew-users-api/pkg/responsewriter"
)

func Monitoring(l zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &monitoringHandler{next: h, logger: l}
	}
}

type monitoringHandler struct {
	next   http.Handler
	logger zap.Logger
}

func (h *monitoringHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	h.logger.Info("new request",
		zap.String("endpoint", r.URL.Path),
	)

	sw := &responsewriter.StatusWriter{Writer: w}

	h.next.ServeHTTP(sw, r)

	h.logger.Info("request handled",
		zap.String("endpoint", r.URL.Path),
		zap.Int("status_code", sw.Status),
		zap.Int64("response_time", time.Since(start).Milliseconds()),
	)
}
