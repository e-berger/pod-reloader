package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/docker/docker/client"
	"github.com/e-berger/pod-reloader/internal/imageref"
)

type Auths struct {
	Auths map[string]Auth `json:"auths"`
}

type Auth struct {
	Auth string `json:"auth"`
}

type RegistryDocker struct {
	Registry
	Auths *Auths
}

func NewRegistryDocker() (IRegistry, error) {
	return &RegistryDocker{
		Registry: Registry{
			Type: DOCKER,
		},
	}, nil
}

func (d *RegistryDocker) SetAuths(auths *Auths) {
	d.Auths = auths
}

func (d *RegistryDocker) GetType() RegistryType {
	return d.Type
}

func (d *RegistryDocker) String() string {
	return "Registry type: " + d.Type.String()
}

func (d *RegistryDocker) RetreiveImage(i *imageref.ImageRef) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	var authent []byte
	if d.Auths != nil {
		authent = d.GetAuth(i)
	}

	encodedRegistryAuth := ""
	if string(authent) != "" {
		encodedRegistryAuth = base64.StdEncoding.EncodeToString(authent)
	}

	imageName := i.Repository + ":" + i.Tag
	out, err := cli.DistributionInspect(ctx, imageName, encodedRegistryAuth)
	if err != nil {
		return "", err
	}
	slog.Info("Image", "digest", out.Descriptor.Digest)
	return out.Descriptor.Digest.String(), nil
}

func (d *RegistryDocker) GetAuth(i *imageref.ImageRef) []byte {
	var authent []byte
	for key, auth := range d.Auths.Auths {
		if strings.HasPrefix(i.Repository, key) {
			authBytes, err := base64.StdEncoding.DecodeString(auth.Auth)
			if err != nil {
				slog.Error("Error decoding auth", "error", err)
				continue
			}
			authParts := strings.Split(string(authBytes), ":")
			username := authParts[0]
			password := authParts[1]
			authent, _ = json.Marshal(map[string]string{
				"username": username,
				"password": password,
			})
			break
		}
	}
	return authent
}
