package memory

import (
	"github.com/sergeysynergy/metricser/pkg/metrics"
)

// Restore Массово загружает переданные значения метрик в хранилища Storage.
func (r *Repo) Restore(prm *metrics.ProxyMetrics) error {
	if prm.Gauges != nil {
		r.gaugesMu.Lock()
		r.gauges = prm.Gauges
		r.gaugesMu.Unlock()
	}

	if prm.Counters != nil {
		r.countersMu.Lock()
		r.counters = prm.Counters
		r.countersMu.Unlock()
	}

	return nil
}
