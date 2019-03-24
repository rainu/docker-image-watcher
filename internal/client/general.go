package client

type DockerRegistryClient interface {
	GetManifest(imageName, tag string) (*ImageManifest, error)
}

type ImageManifest struct {
	MediaType string               `json:"mediaType"`
	Config    ImageManifestConfig  `json:"config"`
	Layer     []ImageManifestLayer `json:"layers"`
}

type ImageManifestConfig struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

type ImageManifestLayer struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}
