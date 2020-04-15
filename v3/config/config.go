package config

import (
    "strings"
    "os"
    "log"
    "bufio"
)

func ReadConfig(path string) map[string]string {
	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

    out := map[string]string{}

    for _, s := range txtlines {
        result := strings.Split(s, " = ")

        if len(result) != 2 {
            continue
        }

        out[result[0]] = result[1]
    }

    return out
}

