package imagemanager

import (
	"context"
	"errors"
	"io"
	"minik8s/logger"
	"minik8s/minik8sTypes"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var (
	K8sLogger = logger.K8sLogger
)

// https://blog.csdn.net/zhonglinzhang/article/details/80697614 image——api的增删改查
type ImageManager struct {
	// contains filtered or unexported fields
	dc *client.Client
}

func NewImageManager(dc *client.Client) *ImageManager {
	return &ImageManager{dc: dc}
}

func (im *ImageManager) PullImage(ctx context.Context, imagePullPolicy minik8sTypes.ImagePullPolicyType, imageName string) error {
	switch imagePullPolicy {
	case minik8sTypes.Always:
		// Always pull image
		// 拉取镜像
		image, err := im.dc.ImagePull(ctx, imageName, types.ImagePullOptions{})
		// println(imageRef)
		if err != nil {
			K8sLogger.Error("PullImage error: ", err)
			return err
		}
		io.Copy(os.Stdout, image)
		defer image.Close()
		return nil
	case minik8sTypes.IfNotPresent:
		// IfNotPresent pull image
		// 检查镜像是否存在，不存在则拉取
		images, err := im.dc.ImageList(ctx, types.ImageListOptions{
			Filters: filters.NewArgs(filters.Arg("reference", imageName)),
		})
		if err != nil {
			K8sLogger.Error("PullImage error: ", err)
			return err
		}
		if len(images) == 0 {
			image, err := im.dc.ImagePull(ctx, imageName, types.ImagePullOptions{})
			if err != nil {
				K8sLogger.Error("PullImage error: ", err)
				return err
			}
			io.Copy(os.Stdout, image)
			defer image.Close()
		}
		return nil
	case minik8sTypes.Never:
		// Never pull image
		// 检查镜像是否存在，不存在则报错
		images, err := im.dc.ImageList(ctx, types.ImageListOptions{
			Filters: filters.NewArgs(filters.Arg("reference", imageName)),
		})
		if err != nil {
			K8sLogger.Error("PullImage error: ", err)
			return err
		}
		if len(images) == 0 {
			K8sLogger.Error("PullImage error: ", errors.New("image not found"))
			return errors.New("image not found")
		}
		return nil
	}
	return nil
}

func (im *ImageManager) RemoveImage(ctx context.Context, imageName string) error {
	_, err := im.dc.ImageRemove(ctx, imageName, types.ImageRemoveOptions{})
	if err != nil {
		K8sLogger.Error("RemoveImage error: ", err)
		return err
	}
	return nil
}
