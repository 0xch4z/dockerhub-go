package dockerhub

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestUserService_GetLoggedInUser(t *testing.T) {
	client, mux, teardown := makeMockClient()
	defer teardown()

	user := &User{}

	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		w.Write(mustJSONMarshal(user))
	})

	res, err := client.User.GetLoggedInUser(context.Background())
	if err != nil {
		t.Errorf("User.GetLoggedInUser returned error: %v", err)
	}

	if !reflect.DeepEqual(res, user) {
		t.Errorf("user is %v; want %v", res, user)
	}
}
