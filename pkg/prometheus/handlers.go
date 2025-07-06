package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func Metrics(f func() ([]string, error)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		for key := range global_labels {
			tmp := global_labels[key]
			tmp.Valid = false
			global_labels[key] = tmp
		}

		raw, err := f()
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}

		for _, metric := range raw {
			err := toMetric(metric)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = writer.Write([]byte(err.Error()))
				return
			}
		}

		for key, val := range global_labels {
			if !val.Valid {
				ok := global_metrics[key].DeleteLabelValues(val.Labels...)
				if ok {
					delete(global_metrics, key)
					delete(global_labels, key)
				}
			}
		}

		promhttp.Handler().ServeHTTP(writer, request)
	}
}
