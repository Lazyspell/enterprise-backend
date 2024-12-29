package checkapi

import (
	"context"
	"encoding/json"
	"net/http"
)

func Liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	lineness := struct {
		Status string
	}{
		Status: "ready and alive",
	}

	return json.NewEncoder(w).Encode(lineness)

}

func Readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	readiness := struct {
		Status string
	}{
		Status: "status is ready",
	}

	return json.NewEncoder(w).Encode(readiness)

}
