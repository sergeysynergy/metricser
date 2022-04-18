package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sergeysynergy/gopracticum/pkg/metrics"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != applicationJSON {
		h.errorJSONUnsupportedMediaType(w)
		return
	}

	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.errorJSONReadBodyFailed(w, err)
		return
	}
	defer r.Body.Close()

	m := metrics.Metrics{}
	err = json.Unmarshal(respBody, &m)
	if err != nil {
		h.errorJSONUnmarshalFailed(w, err)
		return
	}

	switch m.MType {
	case "gauge":
		err = h.hashCheck(&m)
		if err != nil {
			h.errorJSON(w, err.Error(), http.StatusBadRequest)
		}
		h.storage.PutGauge(m.ID, metrics.Gauge(*m.Value))
	case "counter":
		err = h.hashCheck(&m)
		if err != nil {
			h.errorJSON(w, err.Error(), http.StatusBadRequest)
		}
		h.storage.PostCounter(m.ID, metrics.Counter(*m.Delta))
	default:
		err = fmt.Errorf("not implemented")
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)

	h.storeMetrics(w)
}

// записываем метрики в файл если storeInterval равен 0
func (h *Handler) storeMetrics(w http.ResponseWriter) {
	if h.fileStore == nil {
		return
	}
	if h.fileStore.GetStoreInterval() == 0 {
		err := h.fileStore.WriteMetrics()
		if err != nil {
			h.errorJSON(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
