package sinks

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"test.kubebuilder.io/project/api/v1alpha1"
)

type SlackSink struct {
	Endpoint string
	Client   *Client
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
func (s *SlackSink) Configure(config *v1alpha1.Kubegpt) {
	s.Endpoint = config.Spec.Sink.Endpoint
	s.Client = NewClient()
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
		Text: fmt.Sprintf(">*[%s] 파드에 에러가 발생했습니다 : %s %s label: %s image: %s*", k8sgptCR, result.Kind, result.Name, labelsText, imagesText),
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

func (c *Client) SendHTTPRequest(method, url string, body []byte) (*http.Response, error) {
	// HTTP 요청 생성
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.WithError(err).WithField("component", "HTTPClient").Error("Failed to create HTTP request")
		return nil, err
	}

	// 필요한 헤더 설정
	req.Header.Set("Content-Type", "application/json")

	// HTTP 요청 전송
	resp, err := c.hclient.Do(req)
	if err != nil {
		log.WithError(err).WithField("component", "HTTPClient").Error("Failed to send HTTP request")
		return nil, err
	}

	// HTTP 응답 상태 코드 확인 (필요에 따라)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.WithField("status", resp.Status).WithField("component", "HTTPClient").Error("HTTP request returned non-success status")
		return resp, fmt.Errorf("HTTP request returned non-success status: %s", resp.Status)
	}

	return resp, nil
}

func (s *SlackSink) Emit(results v1alpha1.ResultSpec) error {
	message := buildSlackMessage(results, "Kubegpt")
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.WithError(err).WithField("component", "SlackSink").Error("Failed to marshal message")
		return err
	}
	req, err := http.NewRequest(http.MethodPost, s.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.WithError(err).WithField("component", "SlackSink").Error("Failed to create HTTP request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization")

	resp, err := s.Client.hclient.Do(req)
	if err != nil {
		log.WithError(err).WithField("component", "SlackSink").Error("Error sending HTTP request")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithField("status", resp.Status).Error("Failed to send report")
		return fmt.Errorf("failed to send report: %s", resp.Status)
	}
	log.Printf("Successfully sent report to Slack for %s", results.Name)
	return nil
}
