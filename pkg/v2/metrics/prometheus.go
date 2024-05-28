package metrics

import (
	"fmt"
	"io"
	"strings"
)

func (m *Metrics) WritePrometheus(w io.Writer) {
	m.RLock()
	defer m.RUnlock()

	for key, value := range m.metrics {
		typ, ok := metricTypes[key]
		if !ok {
			typ = "counter"
		}

		name := strings.SplitN(key, "{", 2)[0]

		_, _ = w.Write([]byte(fmt.Sprintf(
			"# help %s\n"+
				"# type %s %s\n"+
				"%s %d\n",
			name, name, typ, key, value)))
	}
}
