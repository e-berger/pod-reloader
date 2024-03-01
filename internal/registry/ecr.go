package registry

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/e-berger/pod-reloader/internal/imageref"
)

type RegistryEcr struct {
	Registry
	AwsCfg aws.Config
}

func NewRegistryEcr() (IRegistry, error) {
	return &RegistryEcr{
		Registry: Registry{
			Type: ECR,
		},
	}, nil
}

func (r *RegistryEcr) SetAwsConfig(cfg aws.Config) {
	r.AwsCfg = cfg
}

func (r *RegistryEcr) GetType() RegistryType {
	return r.Type
}

func (r *RegistryEcr) String() string {
	return fmt.Sprintf("Registry type: %s", r.Type)
}

func (r *RegistryEcr) RetreiveImage(i *imageref.ImageRef) (string, error) {
	svc := ecr.NewFromConfig(r.AwsCfg)
	input := &ecr.DescribeImagesInput{
		RepositoryName: &i.Repository,
		ImageIds: []types.ImageIdentifier{
			{
				ImageTag: &i.Tag,
			},
		},
	}

	result, err := svc.DescribeImages(context.TODO(), input)
	if err != nil {
		slog.Error("Error retrieving image", "error", err)
	}
	slog.Info("Image", "image", result)
	return "", nil
}
