package handler

import (
	"net/http"
	"strconv"
)

type MemStorage struct {
	counters map[string]int64
	gauges   map[string]float64
}

func (ms *MemStorage) AddCounter(name string, value string) {
	curr, _ := strconv.ParseInt(value, 10, 64)

	ms.counters[name] = ms.counters[name] + curr
}

func (ms *MemStorage) SetGauge(name string, value string) {
	val, _ := strconv.ParseFloat(value, 64)

	ms.gauges[name] = val
}

func (ms *MemStorage) Processing(res http.ResponseWriter, req *http.Request) {
	metricType := req.PathValue("type")
	metricName := req.PathValue("name")
	metricValue := req.PathValue("value")

	switch metricType {
	case "gauge":
		ms.SetGauge(metricName, metricValue)
		return
	case "counter":
		ms.AddCounter(metricName, metricValue)
		return
	default:
		http.Error(res, "Unsupported metric", http.StatusUnprocessableEntity)
		return
	}
}

var ms = &MemStorage{
	counters: make(map[string]int64),
	gauges:   make(map[string]float64),
}

type Validate struct{}

func (v Validate) Req(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Allow only POST method", http.StatusMethodNotAllowed)

		return
	}

	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(res, "Content-Type must be `text/plain`", http.StatusMethodNotAllowed)

		return
	}
}

func (v Validate) Metric(res http.ResponseWriter, req *http.Request) {
	metricType := req.PathValue("type")
	metricName := req.PathValue("name")
	metricValue := req.PathValue("value")

	if metricName == "" {
		http.Error(res, "Metric name is empty", http.StatusNotFound)

		return
	}

	if metricType != "gauge" && metricType != "counter" {
		http.Error(res, "Allow gauge or counter", http.StatusBadRequest)

		return
	}

	if _, err := strconv.ParseFloat(metricValue, 64); metricType == "gauge" && err != nil {
		http.Error(res, "Gauge value must be float", http.StatusBadRequest)

		return
	}

	if _, err := strconv.ParseInt(metricValue, 10, 64); metricType == "counter" && err != nil {
		http.Error(res, "Counter value must be int", http.StatusBadRequest)

		return
	}
}

func (v Validate) Update(res http.ResponseWriter, req *http.Request) {
	v.Req(res, req)
	v.Metric(res, req)
}

func Update(res http.ResponseWriter, req *http.Request) {
	(Validate{}).Update(res, req)

	ms.Processing(res, req)

	res.WriteHeader(http.StatusOK)
}
