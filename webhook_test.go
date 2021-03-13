package dockerhub

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestWebhookService_CreateWebhook(t *testing.T) {
	client, mux, teardown := makeMockClient()
	defer teardown()
	namespace := "namespace"
	registry := "registry-1.docker.io"
	repo := "repo"
	name := "name"
	url := "https://www.ngrok.com/re78rt/hook"

	hook := &Webhook{}

	uri := fmt.Sprintf("/repositories/%s/%s/", namespace, repo)

	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertNoHeader(t, r, "Authorization")
		assertBody(t, r, string(mustJSONMarshal(&WebhookRequest{
			Name:     name,
			Registry: registry,
			Webhooks: []Webhooks{
				{
					Name:    name,
					HookURL: url,
				},
			},
		},
		)))

		w.WriteHeader(http.StatusCreated)
		w.Write(mustJSONMarshal(&Webhook{}))
	})

	res, err := client.Webhook.CreateWebhook(context.Background(), namespace, repo, name, url)
	if err != nil {
		t.Errorf("Webhook.CreateWebhook returned error: %v", err)
	}

	if !reflect.DeepEqual(res, hook) {
		t.Errorf("webhook is %v; want %v", res, hook)
	}
}

func TestWebhookService_GetWebhooks(t *testing.T) {
	client, mux, teardown := makeMockClient()
	defer teardown()

	namespace := "namespace"
	repo := "repo"

	hook := &Webhook{}

	uri := fmt.Sprintf("/repositories/%s/%s/", namespace, repo)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		w.Write(mustJSONMarshal(hook))
	})

	res, err := client.Webhook.GetWebhooks(context.Background(), namespace, repo)
	if err != nil {
		t.Errorf("Webhook.GetWebhooks returned error: %v", err)
	}

	if !reflect.DeepEqual(res, hook) {
		t.Errorf("webhook is %v; want %v", res, hook)
	}
}
