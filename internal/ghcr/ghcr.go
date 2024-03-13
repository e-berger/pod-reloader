package ghcr

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"log/slog"
)

type Package struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	PackageHTMLURL string `json:"package_html_url"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	HTMLURL        string `json:"html_url"`
	Metadata       struct {
		PackageType string `json:"package_type"`
		Container   struct {
			Tags []string `json:"tags"`
		} `json:"container"`
	} `json:"metadata"`
}

const (
	GHCRIO = "https://api.github.com/orgs/"
)

func GetDigestFromGithub(repository string, tag string, auth map[string]string) (string, error) {

	repositoryParts := strings.Split(repository, "/")
	httpURL := GHCRIO + repositoryParts[1] + "/packages/container/" + repositoryParts[2] + "/versions"
	slog.Info("Github registry", "url", httpURL)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", httpURL, nil)
	if auth != nil && auth["password"] != "" {
		req.Header.Add("Authorization", "Bearer "+auth["password"])
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	jsonData, _ := io.ReadAll(resp.Body)

	// Unmarshal the JSON data into a slice of Package structs
	var packages []Package

	if err := json.Unmarshal([]byte(jsonData), &packages); err != nil {
		return "", err
	}

	// Filter packages with non-empty tags
	for _, p := range packages {
		if len(p.Metadata.Container.Tags) > 0 && contains(p.Metadata.Container.Tags, tag) {
			return p.Name, nil
		}
	}

	return "", nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
