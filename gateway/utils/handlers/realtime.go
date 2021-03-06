package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/realtime"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleRealtimeEvent handles the request coming from the eventing module
func HandleRealtimeEvent(auth *auth.Module, realtime *realtime.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Load the params from the body
		eventDoc := model.CloudEventPayload{}
		json.NewDecoder(r.Body).Decode(&eventDoc)
		defer r.Body.Close()

		// Get the token
		token := utils.GetTokenFromHeader(r)

		// Check if the token is valid
		if err := auth.IsTokenInternal(token); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := realtime.HandleRealtimeEvent(r.Context(), &eventDoc); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

// HandleRealtimeProcessRequest handles the request received from the realtime module. This is a request sent to every gateway
// instance in the cluster to propagate realtime changes
func HandleRealtimeProcessRequest(auth *auth.Module, realtime *realtime.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load the params from the body
		eventDoc := model.CloudEventPayload{}
		json.NewDecoder(r.Body).Decode(&eventDoc)
		defer r.Body.Close()

		// Get the token
		token := utils.GetTokenFromHeader(r)

		// Check if the token is valid
		if err := auth.IsTokenInternal(token); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := realtime.ProcessRealtimeRequests(&eventDoc); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}
