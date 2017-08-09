package credhub

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

// Retrieves a list of stored credential names which contain the search.
func (ch *CredHub) FindByPartialName(nameLike string) ([]credentials.Base, error) {
	panic("Not implemented")
}

// Retrieves a list of stored credential names which are within the specified path.
func (ch *CredHub) FindByPath(path string) ([]credentials.Base, error) {
	var creds map[string][]credentials.Base

	resp, err := ch.Request(http.MethodGet, "/api/v1/data?path="+path, nil)

	if err != nil {
		return []credentials.Base{}, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &creds)

	if err != nil {
		return []credentials.Base{}, err
	}

	return creds["credentials"], nil
}

// Retrieves a list of all paths which contain credentials.
func (ch *CredHub) ShowAllPaths() ([]credentials.Path, error) {
	panic("Not implemented")
}
