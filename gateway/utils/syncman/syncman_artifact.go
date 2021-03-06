package syncman

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
)

func (s *Manager) HandleArtifactRequests(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := admin.IsTokenValid(utils.GetTokenFromHeader(r)); err != nil {
			logrus.Errorf("error handling forwarding artifact request failed to validate token -%v", err)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		r.Host = strings.Split(s.artifactAddr, ":")[0]
		r.URL.Host = s.artifactAddr

		// http: Request.RequestURI can't be set in client requests.
		// http://golang.org/src/pkg/net/http/client.go
		r.RequestURI = ""

		r.URL.Scheme = "http"

		r.URL.Path = "/v1/api/space_cloud/files"

		token, err := admin.GetInternalAccessToken()
		if err != nil {
			logrus.Errorf("error handling forwarding artifact request failed to generate internal access token -%v", err)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// TODO: Use http2 client if that was the incoming request protocol
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		defer utils.CloseTheCloser(response.Body)

		// Copy headers and status code
		w.WriteHeader(response.StatusCode)
		for k, v := range response.Header {
			w.Header().Set(k, v[0])
		}

		// Copy the body
		n, err := io.Copy(w, response.Body)
		if err != nil {
			logrus.Errorf("Failed to copy upstream (%s) response to downstream - %s", r.URL.String(), err.Error())
		}

		logrus.Debugf("Successfully copied %d bytes from upstream server (%s)", n, r.URL.String())
	}

}
