package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"os"
	"time"
)

const (
	FileWriterFormatText       = "text"
	FileWriterFormatAppendJSON = "json"
)

type fileWriter struct {
	metrics   *Metrics
	filename  string
	format    string
	frequency time.Duration
}

func newFileWriter(metrics *Metrics, filename string, format string, frequency time.Duration, start bool) *fileWriter {
	fw := &fileWriter{
		metrics:   metrics,
		filename:  filename,
		format:    format,
		frequency: frequency,
	}
	if start {
		fw.Start()
	}
	return fw
}

func (fw *fileWriter) Start() {
	go func() {
		ticker := time.NewTicker(fw.frequency)
		defer ticker.Stop()
		for {
			<-ticker.C

			// wrap in a func to defer closing the file
			func() {
				perms := os.O_CREATE | os.O_WRONLY
				if fw.format == FileWriterFormatAppendJSON {
					perms |= os.O_APPEND
				} else {
					perms |= os.O_TRUNC
				}

				f, err := os.OpenFile(fw.filename, perms, 0666)
				assert.NoError(err, "failed to open file for writing")
				defer f.Close()

				stats := fw.metrics.GetAll()
				if len(stats) > 0 {
					switch fw.format {
					case FileWriterFormatText:
						fw.writeText(f, stats)
					case FileWriterFormatAppendJSON:
						fw.appendJSON(f, stats)
					default:
						assert.Assert(false, "unknown file writer format: %s", fw.format)
					}
				}
			}()
		}
	}()
}

func writeBytes(f *os.File, data []byte) {
	_, err := f.Write(data)
	assert.NoError(err, "failed to write to file")
}

func writeString(f *os.File, data string) {
	_, err := f.WriteString(data)
	assert.NoError(err, "failed to write to file")
}

func (fw *fileWriter) writeText(f *os.File, stats map[string]int) {
	writeString(f, fmt.Sprintf("timestamp: %s\n", time.Now().Format(time.RFC3339)))
	for key, value := range stats {
		writeString(f, fmt.Sprintf("%s: %d\n", key, value))
	}
}

func (fw *fileWriter) appendJSON(f *os.File, stats map[string]int) {
	stats["timestamp"] = int(time.Now().Unix())
	data, err := json.Marshal(stats)
	assert.NoError(err, "failed to marshal stats")

	writeBytes(f, data)
	writeBytes(f, []byte("\n"))
}
