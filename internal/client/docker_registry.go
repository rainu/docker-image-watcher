package client

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type registryClient struct {
	httpClient *http.Client
}

func NewDockerRegistryClient(httpClient *http.Client) DockerRegistryClient {
	return &registryClient{
		httpClient: httpClient,
	}
}

func (r *registryClient) GetTokenFor(imageName string) (string, error) {
	targetUrl := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", imageName)

	response, err := r.httpClient.Get(targetUrl)
	if err != nil {
		return "", errors.Wrap(err, "Error while sending request")
	}
	if response.Body != nil {
		defer response.Body.Close()
	}

	parsedBody := struct {
		Token string `json:"token"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&parsedBody)
	if err != nil {
		return "", errors.Wrap(err, "Could not decode body")
	}

	return parsedBody.Token, nil
}

func (r *registryClient) GetManifest(imageName, tag string) (*ImageManifest, error) {
	token, err := r.GetTokenFor(imageName)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get auth-token")
	}

	targetUrl := fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", imageName, tag)
	request, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating request")
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	response, err := r.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "Error while sending request")
	}
	if response.Body != nil {
		defer response.Body.Close()
	}

	parsedBody := &ImageManifest{}
	err = json.NewDecoder(response.Body).Decode(parsedBody)
	if err != nil {
		return nil, errors.Wrap(err, "Could not decode body")
	}

	return parsedBody, nil
}
