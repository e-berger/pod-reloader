package registry

import (
	"fmt"

	"github.com/e-berger/pod-reloader/internal/imageref"
)

type IRegistry interface {
	GetType() RegistryType
	RetreiveImage(i *imageref.ImageRef) (string, error)
	String() string
}

type Registry struct {
	Type RegistryType
}

func CreateRegistryFromType(regType RegistryType) (IRegistry, error) {
	switch {
	case regType == ECR:
		return NewRegistryEcr()
	case regType == DOCKER:
		return NewRegistryDocker()
	}
	return nil, fmt.Errorf("registry type %d not found", regType)
}
