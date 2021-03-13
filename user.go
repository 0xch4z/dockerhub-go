package dockerhub

import (
	"context"
	"net/http"
	"time"
)

type UserService service

type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	FullName      string    `json:"full_name"`
	Location      string    `json:"location"`
	Company       string    `json:"company"`
	GravatarEmail string    `json:"gravatar_email"`
	IsStaff       bool      `json:"is_staff"`
	IsAdmin       bool      `json:"is_admin"`
	ProfileURL    string    `json:"profile_url"`
	DateJoined    time.Time `json:"date_joined"`
	GravatarURL   string    `json:"gravatar_url"`
	Type          string    `json:"type"`
}

func (s *UserService) GetLoggedInUser(ctx context.Context) (*User, error) {
	url := "/user/"

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res := &User{}

	if _, err := s.client.Do(ctx, req, res); err != nil {
		return nil, err
	}

	return res, nil
}
