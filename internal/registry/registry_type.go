package registry

import "fmt"

type RegistryType uint

const (
	UNKNOWNMETHOD = iota
	ECR
	DOCKER
)

const (
	ECRString    = "ecr"
	DOCKERString = "docker"
)

func (r RegistryType) String() string {
	switch r {
	case ECR:
		return ECRString
	case DOCKER:
		return DOCKERString
	default:
		panic("unhandled default case")
	}
}

func ParseRegistry(registry string) (RegistryType, error) {
	switch registry {
	case ECRString:
		return ECR, nil
	case DOCKERString:
		return DOCKER, nil
	}
	return UNKNOWNMETHOD, fmt.Errorf("unknown method: %s", registry)
}
