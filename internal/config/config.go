package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/e-berger/pod-reloader/internal/kube"
	"github.com/e-berger/pod-reloader/internal/process"
	"github.com/e-berger/pod-reloader/internal/registry"
)

const SESSIONDURATION = 3600

func SetLogger() {
	lvl := new(slog.LevelVar)
	logLevel := os.Getenv("LOGLEVEL")
	lvl.Set(slog.LevelInfo)
	if logLevel != "" {
		slog.Info("Logger", "loglevel", logLevel)
		lvl.UnmarshalText([]byte(logLevel))
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)
}

func CreateIAMConfig() (*aws.Config, error) {
	var err error
	var cfg aws.Config

	awsProfile := os.Getenv("AWS_PROFILE")
	awsRegion := os.Getenv("AWS_REGION")
	if awsProfile != "" {
		slog.Info("Aws config", "profile", awsProfile)
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile(awsProfile), config.WithRegion(awsRegion))
	} else {
		slog.Info("Aws config default role")
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	}
	if err != nil {
		return nil, err
	}

	assumeRole := os.Getenv("AWS_ASSUME_ROLE")
	if assumeRole != "" {
		slog.Info("Aws config", "assume role", assumeRole)
		stsClient := sts.NewFromConfig(cfg)
		stsCreds := stscreds.NewAssumeRoleProvider(stsClient, assumeRole, func(o *stscreds.AssumeRoleOptions) {
			o.Duration = SESSIONDURATION
		})
		cfg.Credentials = aws.NewCredentialsCache(stsCreds)
	}

	return &cfg, nil
}

func GetRegistryConfig() (registry.IRegistry, error) {
	reg := os.Getenv("REGISTRY")
	if reg == "" {
		return nil, fmt.Errorf("no registry defined")
	}
	registryType, err := registry.ParseRegistry(reg)
	if err != nil {
		return nil, err
	}
	slog.Info("Registry", "registry", registryType)

	registry, err := registry.CreateRegistryFromType(registryType)
	if err != nil {
		return nil, err
	}
	return registry, nil
}

func GetRegistryAuth(p *process.Process) (*registry.Auths, error) {
	secretname := os.Getenv("REGISTRY_AUTH_SECRET")
	if secretname != "" {
		slog.Info("Registry auth", "secretname", secretname)

		namespace := os.Getenv("POD_NAMESPACE")
		if namespace == "" {
			return nil, fmt.Errorf("no namespace defined")
		}

		secretValue, err := kube.GetSecret(p.Client, namespace, secretname)
		if err != nil {
			return nil, err
		}

		secretData, ok := secretValue.Data[".dockerconfigjson"]
		if !ok {
			return nil, fmt.Errorf(".dockerconfigjson not found in secret")
		}

		var auths = &registry.Auths{}
		err = json.Unmarshal(secretData, auths)
		if err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}

		return auths, nil
	}
	return nil, nil
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
