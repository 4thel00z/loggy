package loggy

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12/x/errors"
	"io"
	"log"
	"net/http"
	"os"
)

type Fetcher[T any] func(offset, limit int) (T, error)

func GenerateLogEntriesFetcher(
	client *http.Client,
	host string,
	port int,
	tls bool,
	retries int,
) Fetcher[LogEntries] {
	protocol := "http"
	if tls {
		protocol = "https"
	}
	url := fmt.Sprintf("%s://%s:%d/logs", protocol, host, port)
	logger := log.New(os.Stderr, "[Fetcher]:", log.LstdFlags)

	return func(offset, limit int) (LogEntries, error) {
		for i := 0; i < retries; i++ {
			resp, err := client.Get(url + fmt.Sprintf("?offset=%d&limit=%d", offset, limit))
			if err != nil {
				logger.Println("Error fetching logs:", err.Error())
				continue
			}

			if resp.StatusCode != http.StatusOK {
				logger.Println("Error fetching logs: unexpected status code", resp.StatusCode)
				continue
			}
			content, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Println("Error reading body:", err.Error())
				continue
			}
			var logEntries []LogEntry
			err = json.Unmarshal(content, &logEntries)
			if err != nil {
				logger.Println("Error deserializing logs:", err.Error())
				continue
			}
			return logEntries, nil
		}
		return nil, errors.New("exhausted retries")
	}
}
