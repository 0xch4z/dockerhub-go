package dockerhub

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRepositoryPatchEmitsEmptyFields(t *testing.T) {
	assertMarshalledJSON(t, &RepositoryPatch{}, "{}")
	assertMarshalledJSON(t, &RepositoryPatch{
		FullDescription: "test",
		Description:     "test",
	}, `{"full_description":"test","description":"test"}`)
}

func TestRepositoriesService_EditRepository(t *testing.T) {
	client, mux, teardown := makeMockClient()
	defer teardown()

	namespace := "someone"
	reponame := "somerepo"

	repo := &Repository{
		User:           namespace,
		Name:           reponame,
		Namespace:      namespace,
		RepositoryType: "image",
		Status:         1,
		Affiliation:    String("owner"),
	}
	patch := &RepositoryPatch{
		FullDescription: "# Hello, World!",
		Description:     "My Repo (:",
	}

	uri := fmt.Sprintf("/repositories/%s/%s/", namespace, reponame)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPatch)
		assertBody(t, r, string(mustJSONMarshal(patch)))
		w.WriteHeader(http.StatusCreated)
		w.Write(mustJSONMarshal(repo))
	})

	res, err := client.Repositories.EditRepository(context.Background(), namespace, reponame, patch)
	if err != nil {
		t.Errorf("Repositories.EditRepository returned error: %v", err)
	}

	if !reflect.DeepEqual(res, repo) {
		t.Errorf("repository is %v; want %v", res, repo)
	}
}

func TestRepositoriesService_GetRepository(t *testing.T) {
	client, mux, teardown := makeMockClient()
	defer teardown()

	namespace := "library"
	reponame := "ubuntu"
	repo := &Repository{
		User:           namespace,
		Name:           reponame,
		Namespace:      namespace,
		RepositoryType: "image",
		Status:         1,
	}

	uri := fmt.Sprintf("/repositories/%s/%s/", namespace, reponame)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		w.Write(mustJSONMarshal(repo))
	})

	res, err := client.Repositories.GetRepository(context.Background(), namespace, reponame)
	if err != nil {
		t.Errorf("Repositories.GetRepository returned error: %v", err)
	}

	if !reflect.DeepEqual(res, repo) {
		t.Errorf("repository is %v; want %v", res, repo)
	}
}

func TestRepositoriesService_SetRepositoryPrivacy(t *testing.T) {
	for _, tc := range []struct {
		isPrivate bool
	}{{true}, {false}} {
		client, mux, teardown := makeMockClient()
		defer teardown()

		namespace := "pulumi"
		repo := "pulumi"

		uri := fmt.Sprintf("/repositories/%s/%s/privacy/", namespace, repo)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			assertMethod(t, r, http.MethodPost)
			assertBody(t, r, string(mustJSONMarshal(&RepositoryPrivacyPatch{tc.isPrivate})))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
		})

		if err := client.Repositories.SetRepositoryPrivacy(context.Background(), namespace, repo, tc.isPrivate); err != nil {
			t.Errorf("Repositories.SetRepositoryPrivacy returned error: %v", err)
		}
	}
}

func TestRepositoriesService_CreateRepository(t *testing.T) {
	client, mux, teardown := makeMockClient()
	defer teardown()

	namespace := "namespace"
	registry := "docker"
	name := "name"
	description := "description"
	privacy := "public"
	isPrivate := false
	list := &Repository{}
	mux.HandleFunc("/repositories/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertNoHeader(t, r, "Authorization")
		assertBody(t, r, string(mustJSONMarshal(&CreateRepositoryRequest{
			Namespace:   namespace,
			Registry:    registry,
			Image:       fmt.Sprintf("%s/%s", namespace, name),
			Name:        name,
			Description: description,
			Privacy:     privacy,
			IsPrivate:   isPrivate,
		})))
		w.WriteHeader(http.StatusCreated)
		w.Write(mustJSONMarshal(&Repository{}))
	})

	res, err := client.Repositories.CreateRepository(context.Background(), namespace, name, description, isPrivate)

	if err != nil {
		t.Errorf("Repositories.CreateRepository returned error: %v", err)
	}

	if !reflect.DeepEqual(res, list) {
		t.Errorf("repository list is %v; want %v", res, list)
	}
}

func TestRepositoriesService_GetRepositories(t *testing.T) {
	client, mux, teardown := makeMockClient()
	defer teardown()

	namespace := "pulumi"
	list := &RepositoryList{}

	uri := fmt.Sprintf("/repositories/%s/", namespace)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		w.Write(mustJSONMarshal(list))
	})

	res, err := client.Repositories.GetRepositories(context.Background(), namespace)
	if err != nil {
		t.Errorf("Repositories.GetRepositories returned error: %v", err)
	}

	if !reflect.DeepEqual(res, list) {
		t.Errorf("repository list is %v; want %v", res, list)
	}
}
