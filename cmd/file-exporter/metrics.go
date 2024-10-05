package main

import (
	"errors"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

type CachedLabels struct {
	Labels []string
	Valid  bool
}

var (
	global_labels  = map[string]CachedLabels{}
	global_metrics = map[string]*prometheus.MetricVec{}
)

func toMetric(raw string) error {
	s := strings.Split(raw, "=")
	key := s[0]
	raw_val := s[1]
	key_parts := strings.Split(key, ",")
	key_name := key_parts[0]
	key_namespace := key_parts[1]
	labels := []string{}
	vals := []string{}

	for i := 2; i < len(key_parts); i += 2 {
		labels = append(labels, key_parts[i])
		vals = append(vals, key_parts[i+1])
	}

	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      key_name,
			Namespace: key_namespace,
		},
		labels,
	)

	val := util.Atof(raw_val)
	are := &prometheus.AlreadyRegisteredError{}
	err := prometheus.Register(metric)
	if errors.As(err, are) {
		metric = are.ExistingCollector.(*prometheus.GaugeVec)
	} else if err != nil {
		return err
	}
	// prometheus.MustRegister(metric)
	global_labels[s[0]] = CachedLabels{
		Valid:  true,
		Labels: vals,
	}
	global_metrics[s[0]] = metric.MetricVec
	metric.WithLabelValues(vals...).Set(val)

	return nil
}
