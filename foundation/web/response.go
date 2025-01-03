package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
)

// Respond converts a Go value to JSON and sends it to the client
func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	setStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func GraphqlResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, data *handler.Server, statusCode int) error {
	// ServeHTTP will handle the request and write the response
	data.ServeHTTP(w, r)

	return nil
}
