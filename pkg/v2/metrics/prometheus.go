package metrics

import (
	"fmt"
	"io"
)

func (m *Metrics) WritePrometheus(w io.Writer) {
	m.RLock()
	defer m.RUnlock()

	for key, value := range m.metrics {
		typ, ok := metricTypes[key]
		if !ok {
			typ = "counter"
		}

		_, _ = w.Write([]byte(fmt.Sprintf(
			"# help %s\n"+
				"# type %s %s\n"+
				"%s %d\n",
			key, key, typ, key, value)))
	}
}
