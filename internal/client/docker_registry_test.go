package client

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDockerRegistry_GetManifest(t *testing.T) {
	//given
	toTest := NewDockerRegistryClient(http.DefaultClient)

	//when
	manifest, err := toTest.GetManifest("library/alpine", "latest")

	//then
	assert.NoError(t, err)
	assert.NotNil(t, manifest)
}
