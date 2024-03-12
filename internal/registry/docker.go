package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/docker/docker/client"
	"github.com/e-berger/pod-reloader/internal/ghcr"
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

	var authent []byte
	var authentValue map[string]string
	if d.Auths != nil {
		authentValue = d.GetAuth(i)
		if authentValue != nil {
			authent, _ = json.Marshal(authentValue)
		}
	}

	digest := ""
	var err error

	if strings.Contains(i.Repository, "ghcr.io") {
		digest, err = ghcr.GetDigestFromGithub(i.Repository, i.Tag, authentValue)
		if err != nil {
			return "", err
		}
	} else {
		encodedRegistryAuth := ""
		if string(authent) != "" {
			encodedRegistryAuth = base64.StdEncoding.EncodeToString(authent)
		}
		imageName := i.Repository + ":" + i.Tag
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.30"))
		if err != nil {
			return "", err
		}
		out, err := cli.DistributionInspect(ctx, imageName, encodedRegistryAuth)
		if err != nil {
			return "", nil
		}
		digest = out.Descriptor.Digest.String()
	}

	slog.Info("Image", "digest", digest)
	return digest, nil
}

func (d *RegistryDocker) GetAuth(i *imageref.ImageRef) map[string]string {
	authent := make(map[string]string)
	for key, auth := range d.Auths.Auths {
		if strings.HasPrefix(i.Repository, key) {
			authBytes, err := base64.StdEncoding.DecodeString(auth.Auth)
			if err != nil {
				slog.Error("Error decoding auth", "error", err)
				continue
			}
			authParts := strings.Split(string(authBytes), ":")
			authent["username"] = authParts[0]
			authent["password"] = authParts[1]
			break
		}
	}
	return authent
}
