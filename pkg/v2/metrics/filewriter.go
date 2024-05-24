package metrics

import (
	"fmt"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"os"
	"time"
	"encoding/json"
)

const (
	FileWriterFormatText = "text"
	FileWriterFormatJSON = "json"
	FileWriterFormatHTML = "html"
)

type fileWriter struct {
	metrics   *Metrics
	filename  string
	format    string
	frequency time.Duration
}

func newFileWriter(metrics *Metrics, filename string, format string, frequency time.Duration) *fileWriter {
	return &fileWriter{
		metrics:   metrics,
		filename:  filename,
		format:    format,
		frequency: frequency,
	}
}

func (fw *fileWriter) Start() {
	go func() {
		var err error
		ticker := time.NewTicker(fw.frequency)
		defer ticker.Stop()
		for {
			<-ticker.C

			var f *os.File
			f, err = os.OpenFile(fw.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				assert.NoError(err, "unable to open file for writing")
			}

			stats := fw.metrics.GetAll()
			switch fw.format {
			case FileWriterFormatText:
				fw.writeText(f, stats)
			case FileWriterFormatJSON:
				fw.writeJSON(f, stats)
			case FileWriterFormatHTML:
				fw.writeHTML(f, stats)
			default:
				assert.Assert(false, "unknown file writer format: %s", fw.format)
			}

			_ = f.Close()
		}
	}()
}

func (fw *fileWriter) writeText(f *os.File, stats map[string]int) {
	for key, value := range stats {
		_, err := f.WriteString(fmt.Sprintf("%s: %d\n", key, value))
		if err != nil {
			assert.NoError(err, "unable to write to file")
		}
	}
}

func (fw *fileWriter) writeJSON(f *os.File, stats map[string]int) {
	data, err := json.Marshal(stats)
	if err != nil {
		assert.NoError(err, "unable to marshal stats")
	}
	_, err = f.Write(data)
}

func (fw *fileWriter) writeHTML(f *os.File, stats map[string]int) {
	_, err := f.WriteString("<html><body><table>")
	if err != nil {
		assert.NoError(err, "unable to write to file")
	}

	for key, value := range stats {
		_, err := f.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", key, value))
		if err != nil {
			assert.NoError(err, "unable to write to file")
		}
	}

	_, err = f.WriteString("</table></body></html>")
	if err != nil {
		assert.NoError(err, "unable to write to file")
	}
}
