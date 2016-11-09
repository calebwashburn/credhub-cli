package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/credhub-cli/repositories"

	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-cf/credhub-cli/client/clientfakes"
	"github.com/pivotal-cf/credhub-cli/config"
	cmcli_errors "github.com/pivotal-cf/credhub-cli/errors"
	"github.com/pivotal-cf/credhub-cli/models"
)

var _ = Describe("CaRepository", func() {
	var (
		repository Repository
		httpClient clientfakes.FakeHttpClient
		cfg        config.Config
	)

	Describe("SendRequest", func() {
		BeforeEach(func() {
			repository = NewCaRepository(&httpClient)
			cfg = config.Config{
				ApiURL:  "http://example.com",
				AuthURL: "http://example.com",
			}
		})

		Context("when there is a response body", func() {
			It("sends a request to the server", func() {
				request, _ := http.NewRequest("PUT", "http://example.com/foo", nil)

				responseObj := http.Response{
					StatusCode: 200,
					Body: ioutil.NopCloser(
						bytes.NewReader([]byte(`{"data":[{"type":"root","value":{"certificate":"my-cert","private_key":"my-priv"},"updated_at":"2016-01-01T12:00:00Z"}]}`)),
					),
				}

				httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
					Expect(req).To(Equal(request))

					return &responseObj, nil
				}

				caParams := models.CaParameters{
					Certificate: "my-cert",
					PrivateKey:  "my-priv",
				}
				expectedCaBody := models.CaBody{
					ContentType: "root",
					Value:       &caParams,
					UpdatedAt:   "2016-01-01T12:00:00Z",
				}
				expectedCa := models.Ca{
					Name:   "foo",
					CaBody: expectedCaBody,
				}

				ca, err := repository.SendRequest(request, "foo")

				Expect(err).ToNot(HaveOccurred())
				Expect(ca).To(Equal(expectedCa))
			})
		})

		Describe("Errors", func() {
			It("returns a NewResponseError when the JSON response cannot be parsed", func() {
				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("adasdasdasasd"))),
				}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				_, error := repository.SendRequest(request, "foo")
				Expect(error).To(MatchError(cmcli_errors.NewResponseError()))
			})
		})
	})
})
