package notifier

import (
	"fmt"

	"github.com/nlopes/slack"
)

// AttachmentRegistry contains all Slack Attachments
type AttachmentRegistry struct {
	Attachments []slack.Attachment
}

// NewRegistry creates empty registry
func NewRegistry() *AttachmentRegistry {
	return &AttachmentRegistry{}
}

// AddStatus adds status for the single token into the AttachmentRegistry
func (r *AttachmentRegistry) AddStatus(token string, success bool, message string) {
	attachment := slack.Attachment{}

	// success
	color := "#008000"

	// fail
	if !success {
		color = "#FF0000"
	}

	attachment.Color = color
	attachment.Text = fmt.Sprintf("Token %s...: %s", token[0:8], message)

	r.Attachments = append(r.Attachments, attachment)
}

// NotifySlack sends notification to the slack webhook
func NotifySlack(registry *AttachmentRegistry, webhookURL string) error {
	msg := slack.WebhookMessage{Text: "Vault token renewal status", Attachments: registry.Attachments}
	err := slack.PostWebhook(webhookURL, &msg)
	if err != nil {
		return fmt.Errorf("could not send slack message: %s", err)
	}

	return nil
}
