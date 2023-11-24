package sinks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"test.kubebuilder.io/project/api/v1alpha1"
)

var _ ISink = (*SlackSink)(nil)

type SlackSink struct {
	Endpoint string
	Kubegpt  string
	Client   Client
}

type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Color string `json:"color"`
	Title string `json:"title"`
}

func formatEvent(event v1alpha1.Event) string {
	// Format the event into a string.
	// Modify this according to how you want to display each event.
	return fmt.Sprintf("Event: %s - Count: %v - Reason: %s - Message: %s", event.Type, event.Count, event.Reason, event.Message)
}
func buildSlackMessage(result v1alpha1.ResultSpec, k8sgptCR string) SlackMessage {
	var detailsText string
	for _, event := range result.Event {
		// Use the formatEvent function to get the string representation of each event
		detailsText += formatEvent(event) + "\n" // Add a newline after each event
	}

	labelsText := fmt.Sprintf("%v", result.Labels)
	imagesText := fmt.Sprintf("%v", result.Images)

	return SlackMessage{
		Text: fmt.Sprintf(">*[%s] K8sGPT analysis of the %s %s %s %s*", k8sgptCR, result.Kind, result.Name, labelsText, imagesText),
		Attachments: []Attachment{
			{
				Type:  "mrkdwn",
				Text:  detailsText,
				Color: "danger",
				Title: "Detailed Report",
			},
		},
	}
}

func (s *SlackSink) Configure(config v1alpha1.Kubegpt, c Client) {
	log.Printf("Configured Slack sink for %v", s)
	s.Endpoint = config.Spec.Sink.Endpoint
	s.Client = c
	// take the name of the K8sGPT Custom Resource
	//s.Kubegpt = config
	log.Printf("Configured Slack sink for %v", s)
}

func (s *SlackSink) Emit(results v1alpha1.ResultSpec) error {
	message := buildSlackMessage(results, "Kubegpt")
	payload, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPost, s.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.hclient.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send report: %s", resp.Status)
		return fmt.Errorf("failed to send report: %s", resp.Status)
	}
	log.Printf("Successfully sent report to Slack for %s", results.Name)
	return nil
}
