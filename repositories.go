package dockerhub

import (
	"context"
	"fmt"
	"net/http"
)

// RepositoriesService handles communication with the repository
// related methods of the Dockerhub API.
type RepositoriesService service

// RepositoryPatch represents payload to patch a Repository.
type RepositoryPatch struct {
	FullDescription string `json:"full_description,omitempty"`
	Description     string `json:"description,omitempty"`
}

// RepositoryPrivacyPatch represents payload to patch a Repository's
// privacy mode.
type RepositoryPrivacyPatch struct {
	IsPrivate bool `json:"is_private"`
}

// RepositoryPermissions specifies the permissions of the requesting user
// to the given Repository.
type RepositoryPermissions struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
	Admin bool `json:"admin"`
}

// RepositoryList represents a list of repositories with pagination details.
type RepositoryList struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`

	Results []Repository `json:"results"`
}

// Repository represents a Dockerhub repository.
type Repository struct {
	User            string  `json:"user"`
	Name            string  `json:"name"`
	Namespace       string  `json:"namespace"`
	RepositoryType  string  `json:"repository_type"`
	Status          int     `json:"status"`
	Description     string  `json:"description"`
	IsPrivate       bool    `json:"is_private"`
	IsAutomated     bool    `json:"is_automated"`
	CanEdit         bool    `json:"can_edit"`
	StarCount       int     `json:"star_count"`
	PullCount       int     `json:"pull_count"`
	LastUpdated     string  `json:"last_updated"`
	IsMigrated      bool    `json:"is_migrated"`
	HasStarred      bool    `json:"has_starred"`
	FullDescription string  `json:"full_description"`
	Affiliation     *string `json:"affiliation"`

	Permissions RepositoryPermissions `json:"repository_permissions"`
}

type CreateRepositoryRequest struct {
	Namespace     string        `json:"namespace"`
	Registry      string        `json:"registry"`
	Image         string        `json:"image"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Privacy       string        `json:"privacy"`
	BuildSettings []interface{} `json:"build_settings"`
	IsPrivate     bool          `json:"is_private"`
}

func (s RepositoriesService) buildRepoSlug(namespace, repo string) string {
	return fmt.Sprintf("/repositories/%s/%s/", namespace, repo)
}

// CreateRepository create a repository.
func (s *RepositoriesService) CreateRepository(ctx context.Context, namespace, name, description string, isPrivate bool) (*Repository, error) {
	url := "/repositories/"
	repo := &CreateRepositoryRequest{
		Namespace:   namespace,
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
		Registry:    "docker",
		Image:       fmt.Sprintf("%s/%s", namespace, name),
	}

	if repo.IsPrivate {
		repo.Privacy = "private"
	} else {
		repo.Privacy = "public"
	}

	req, err := s.client.NewRequest(http.MethodPost, url, repo)

	if err != nil {
		return nil, err
	}

	res := &Repository{}

	if _, err := s.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// EditRepository updates a repository.
func (s *RepositoriesService) EditRepository(ctx context.Context, namespace, repo string, patch *RepositoryPatch) (*Repository, error) {
	slug := s.buildRepoSlug(namespace, repo)
	req, err := s.client.NewRequest(http.MethodPatch, slug, patch)
	if err != nil {
		return nil, err
	}

	res := &Repository{}
	if _, err := s.client.Do(ctx, req, res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetRepository gets details for a given repository.
func (s *RepositoriesService) GetRepository(ctx context.Context, namespace, repo string) (*Repository, error) {
	slug := s.buildRepoSlug(namespace, repo)
	req, err := s.client.NewRequest(http.MethodGet, slug, nil)
	if err != nil {
		return nil, err
	}

	res := &Repository{}
	if _, err := s.client.Do(ctx, req, res); err != nil {
		return nil, err
	}

	return res, nil
}

// SetRepositoryPrivacy sets the privacy status of a repository.
func (s *RepositoriesService) SetRepositoryPrivacy(ctx context.Context, namespace, repo string, isPrivate bool) error {
	slug := s.buildRepoSlug(namespace, repo) + "privacy/"
	req, err := s.client.NewRequest(http.MethodPost, slug, &RepositoryPrivacyPatch{
		IsPrivate: isPrivate,
	})
	if err != nil {
		return err
	}

	if _, err := s.client.Do(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

// GetRepositories gets all repositories from a given Dockerhub namespace.
func (s *RepositoriesService) GetRepositories(ctx context.Context, namespace string) (*RepositoryList, error) {
	slug := fmt.Sprintf("/repositories/%s/", namespace)
	req, err := s.client.NewRequest(http.MethodGet, slug, nil)
	if err != nil {
		return nil, err
	}

	res := &RepositoryList{}
	if _, err := s.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}
