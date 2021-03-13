package dockerhub

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type WebhookService service

type Webhooks struct {
	Name        string    `json:"name"`
	HookURL     string    `json:"hook_url"`
	Created     time.Time `json:"created"`
	LastUpdated time.Time `json:"last_updated"`
}

type WebhookRequest struct {
	Name                string     `json:"name"`
	ExpectFinalCallback bool       `json:"expect_final_callback"`
	Webhooks            []Webhooks `json:"webhooks"`
	Registry            string     `json:"registry"`
}

type Webhook struct {
	Name                string     `json:"name"`
	Slug                string     `json:"slug"`
	ExpectFinalCallback bool       `json:"expect_final_callback"`
	Webhooks            []Webhooks `json:"webhooks"`
	Created             time.Time  `json:"created"`
	LastUpdated         time.Time  `json:"last_updated"`
}

func (s WebhookService) buildWebhookSlug(namespace, repo string) string {
	return fmt.Sprintf("/repositories/%s/%s/webhook_pipeline/", namespace, repo)
}
func (s *WebhookService) CreateWebhook(ctx context.Context, namespace, repo, name, url string) (*Webhook, error) {
	slug := s.buildWebhookSlug(namespace, repo)

	hook := &WebhookRequest{
		Name:                name,
		ExpectFinalCallback: false,
		Registry:            "registry-1.docker.io",
	}

	hook.Webhooks = append(hook.Webhooks, Webhooks{
		Name:    name,
		HookURL: url,
	})

	req, err := s.client.NewRequest(http.MethodPost, slug, hook)

	if err != nil {
		return nil, err
	}

	res := &Webhook{}

	if _, err := s.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}
