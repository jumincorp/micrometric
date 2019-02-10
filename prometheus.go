package micrometric

import (
	"bytes"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

var (
	mutex            = &sync.Mutex{}
	formattedMetrics = make([][]byte, 0)
)

type prometheusExporter struct {
	Exporter
	address string
}

func NewPrometheusExporter(address string) Exporter {
	p := new(prometheusExporter)
	p.address = address
	return p
}

func httpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		for _, metric := range formattedMetrics {
			w.Write([]byte(metric))
		}
	}
}

func (p *prometheusExporter) Setup() {
	http.Handle("/metrics", httpHandler())
	log.Fatal(http.ListenAndServe(p.address, nil))
}

func formatMetric(m Metric) []byte {
	bb := bytes.Buffer{}

	bb.WriteString(m.Name)

	sortedLabels := make([]string, len(m.Labels))
	i := 0
	for k := range m.Labels {
		sortedLabels[i] = k
		i++
	}
	sort.Strings(sortedLabels)

	bb.WriteRune('{')
	for i := range sortedLabels {
		if i != 0 {
			bb.WriteRune(',')
		}
		bb.WriteString(sortedLabels[i])
		bb.WriteString("=\"")
		bb.WriteString(m.Labels[sortedLabels[i]])
		bb.WriteRune('"')
	}
	bb.WriteString("} ")

	bb.WriteString(strconv.FormatFloat(m.Value, 'f', -1, 64))

	bb.WriteRune('\n')

	return bb.Bytes()
}

func (p *prometheusExporter) Export(metrics []Metric) error {
	var err error

	mutex.Lock()
	defer mutex.Unlock()

	formattedMetrics = make([][]byte, len(metrics))

	for i, m := range metrics {
		//formattedLabels := formatLabels(m.labels)
		formattedMetrics[i] = formatMetric(m)
	}
	return err
}
