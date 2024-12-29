package checkapi

import "net/http"

func Routes(mux *http.ServeMux) {

	mux.HandleFunc("/liveness", Liveness)
	mux.HandleFunc("/readiness", Readiness)

}
