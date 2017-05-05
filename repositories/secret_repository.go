package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	cm_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type secretRepository struct {
	httpClient client.HttpClient
}

func NewSecretRepository(httpClient client.HttpClient) Repository {
	return secretRepository{httpClient: httpClient}
}

func (r secretRepository) SendRequest(request *http.Request, identifier string) (models.Printable, error) {
	secret := models.Secret{}
	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return secret, err
	}

	if request.Method == "DELETE" {
		return secret, nil
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&secret.SecretBody)

	if err != nil {
		return secret, cm_errors.NewResponseError()
	}

	if data, ok := secret.SecretBody["data"].([]interface{}); ok {
		secret.SecretBody = data[0].(map[string]interface{})
	}

	return secret, nil
}
