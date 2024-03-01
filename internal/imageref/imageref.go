package imageref

import (
	"fmt"
	"log/slog"
	"strings"

	v1 "k8s.io/api/core/v1"
)

type ImageRef struct {
	Repository string
	Digest     string
	Tag        string
}

func ExtractImages(pod v1.Pod) []*ImageRef {
	var listImage []*ImageRef
	for _, c := range pod.Status.ContainerStatuses {
		imageRef, err := NewImageRef(c)
		if err != nil {
			slog.Error("Invalid container", "error", err)
		} else {
			listImage = append(listImage, imageRef)
		}
	}
	return listImage
}

func NewImageRef(c v1.ContainerStatus) (*ImageRef, error) {
	var image ImageRef

	if strings.Contains(c.ImageID, "@") {
		// Handle SHA1 digest
		parts := strings.SplitN(c.ImageID, "@", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid image reference: %s", c.ImageID)
		}
		image.Digest = parts[1]
		image.Repository = parts[0]
	} else {
		return nil, fmt.Errorf("invalid image reference: %s", c.ImageID)
	}

	if strings.Contains(c.Image, ":") {
		parts := strings.SplitN(c.Image, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid image reference: %s", c.Image)
		}
		image.Tag = parts[1]
	} else {
		return nil, fmt.Errorf("invalid image reference: %s", c.Image)
	}

	return &image, nil
}
