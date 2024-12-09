package rabbitmq

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-rmq-rnd/consumer/config"
)

type Data struct {
	Arguments                     Arguments                     `json:"arguments"`
	AutoDelete                    bool                          `json:"auto_delete"`
	ConsumerCapacity              float64                       `json:"consumer_capacity"`
	ConsumerUtilisation           float64                       `json:"consumer_utilisation"`
	Consumers                     int                           `json:"consumers"`
	Durable                       bool                          `json:"durable"`
	EffectivePolicyDefinition     EffectivePolicyDefinition     `json:"effective_policy_definition"`
	Exclusive                     bool                          `json:"exclusive"`
	Memory                        int                           `json:"memory"`
	MessageBytes                  int                           `json:"message_bytes"`
	MessageBytesPagedOut          int                           `json:"message_bytes_paged_out"`
	MessageBytesPersistent        int                           `json:"message_bytes_persistent"`
	MessageBytesRAM               int                           `json:"message_bytes_ram"`
	MessageBytesReady             int                           `json:"message_bytes_ready"`
	MessageBytesUnacknowledged    int                           `json:"message_bytes_unacknowledged"`
	Messages                      int                           `json:"messages"`
	MessagesDetails               MessagesDetails               `json:"messages_details"`
	MessagesPagedOut              int                           `json:"messages_paged_out"`
	MessagesPersistent            int                           `json:"messages_persistent"`
	MessagesRAM                   int                           `json:"messages_ram"`
	MessagesReady                 int                           `json:"messages_ready"`
	MessagesReadyDetails          MessagesReadyDetails          `json:"messages_ready_details"`
	MessagesReadyRAM              int                           `json:"messages_ready_ram"`
	MessagesUnacknowledged        int                           `json:"messages_unacknowledged"`
	MessagesUnacknowledgedDetails MessagesUnacknowledgedDetails `json:"messages_unacknowledged_details"`
	MessagesUnacknowledgedRAM     int                           `json:"messages_unacknowledged_ram"`
	Name                          string                        `json:"name"`
	Node                          string                        `json:"node"`
	Reductions                    int                           `json:"reductions"`
	ReductionsDetails             ReductionsDetails             `json:"reductions_details"`
	State                         string                        `json:"state"`
	StorageVersion                int                           `json:"storage_version"`
	Type                          string                        `json:"type"`
	Vhost                         string                        `json:"vhost"`
}
type Arguments struct {
	XMessageTTL int `json:"x-message-ttl"`
}
type EffectivePolicyDefinition struct {
}
type MessagesDetails struct {
	Rate float64 `json:"rate"`
}
type MessagesReadyDetails struct {
	Rate float64 `json:"rate"`
}
type MessagesUnacknowledgedDetails struct {
	Rate float64 `json:"rate"`
}
type ReductionsDetails struct {
	Rate float64 `json:"rate"`
}

func GetDataQueue(username, password, queueName string) ([]Data, error) {

	// Create an HTTP client
	client := &http.Client{}

	url := "http://localhost:15672"
	if config.Config.Environment == "production" {
		url = "rmq.jubelio.com"
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/queues", url), nil)
	if err != nil {
		log.Printf("Failed to create request: %v\n", err)
		return nil, err
	}

	// Add Basic Auth
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch queues: %s\n", resp.Status)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		return nil, err
	}

	var queues []Data
	err = json.Unmarshal(body, &queues)
	if err != nil {
		log.Printf("Failed to parse JSON: %v\n", err)
		return nil, err
	}

	newQueue := make([]Data, 0)
	for _, eachQueue := range queues {
		if eachQueue.Name == queueName {
			newQueue = append(newQueue, eachQueue)
		}
	}
	return newQueue, nil
}
