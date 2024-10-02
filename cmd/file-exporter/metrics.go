package main

import (
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	global_labels  = map[string]bool{}
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
	prometheus.Unregister(metric)
	prometheus.MustRegister(metric)
	global_labels[key_name+"_"+key_namespace] = true
	global_metrics[key_name+"_"+key_namespace] = metric.MetricVec
	metric.WithLabelValues(vals...).Set(val)

	return nil
}
