package ghcr

import (
	b64 "encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"log/slog"
)

type Platform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

// Define the struct for the manifest, including optional annotations
type Manifest struct {
	MediaType   string            `json:"mediaType"`
	Digest      string            `json:"digest"`
	Size        int               `json:"size"`
	Platform    Platform          `json:"platform,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Define the top-level struct for the whole JSON object
type ImageIndex struct {
	SchemaVersion int        `json:"schemaVersion"`
	MediaType     string     `json:"mediaType"`
	Manifests     []Manifest `json:"manifests"`
}

const (
	GHCRIO = "https://ghcr.io/v2/"
)

func GetDigestFromGithub(repository string, tag string, auth map[string]string) (string, error) {
	httpURL := strings.Replace(repository, "ghcr.io/", GHCRIO, 1) + "/manifests/" + tag
	slog.Info("Github registry", "url", httpURL)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", httpURL, nil)
	if auth != nil && auth["password"] != "" {
		auth := b64.StdEncoding.EncodeToString([]byte(auth["password"]))
		req.Header.Add("Authorization", "Bearer "+auth)
	}
	req.Header.Add("Accept", "application/vnd.oci.image.index.v1+json")
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	imageIndex, err := parseManifest(string(body))
	if err != nil {
		return "", err
	}

	for _, manifest := range imageIndex.Manifests {
		if manifest.Platform.Architecture == "amd64" && manifest.Platform.OS == "linux" {
			return manifest.Digest, nil
		}
	}
	return "", nil
}

func parseManifest(data string) (*ImageIndex, error) {
	var imageIndex ImageIndex
	err := json.Unmarshal([]byte(data), &imageIndex)
	if err != nil {
		return nil, err
	}

	return &imageIndex, nil
}
