package checkapi

import (
	"encoding/json"
	"net/http"
)

func Liveness(w http.ResponseWriter, r *http.Request) {
	lineness := struct {
		Status string
	}{
		Status: "ready and alive",
	}

	json.NewEncoder(w).Encode(lineness)

}

func Readiness(w http.ResponseWriter, r *http.Request) {
	readiness := struct {
		Status string
	}{
		Status: "status is ready",
	}

	json.NewEncoder(w).Encode(readiness)

}
