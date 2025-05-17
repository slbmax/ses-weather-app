package static

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"time"

	"github.com/slbmax/ses-weather-app/assets"
)

type IndexData struct {
	BaseApiUrl string
}

func Serve(ctx context.Context, data IndexData, listener net.Listener) error {
	tmpl, err := template.ParseFS(assets.IndexHTML, assets.TemplateIndexHTML)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	preRendered := buf.Bytes()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(preRendered)
	})
	srv := &http.Server{Handler: mux}

	// graceful shutdown
	go func() {
		<-ctx.Done()

		shutdownDeadline, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_ = srv.Shutdown(shutdownDeadline)
	}()

	if err = srv.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
