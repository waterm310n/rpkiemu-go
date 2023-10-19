package setup

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/waterm310n/rpkiemu-go/ca/krillop"
	"github.com/waterm310n/rpkiemu-go/ca/rsyncdop"
	"github.com/waterm310n/rpkiemu-go/k8sexec"
)

type PublishPoint struct {
	Namespace     string `json:"namespace,omitempty"`
	PodName       string `json:"pod_name,omitempty"`
	ContainerName string `json:"container_name,omitempty"`
	IsRIR         bool   `json:"is_rir,omitempty"`
}

func configureKrill(dataDir string, publishPoints map[string]PublishPoint) {
	for _, v := range publishPoints {
		if err := createKrillConfig(dataDir, v.PodName, v.IsRIR); err != nil {
			slog.Error(err.Error())
		}
		if execOptions, err := k8sexec.NewExecOptions(v.Namespace, v.PodName, v.ContainerName); err == nil {
			kCA := krillop.NewKrillK8sCA(execOptions, v.IsRIR)
			if err := kCA.Configure(dataDir); err != nil {
				slog.Debug(err.Error())
			}
		}
	}
}

func configureRsyncd(publishPoints map[string]PublishPoint) {

	for _, v := range publishPoints {
		containerName := v.PodName + "-rsyncd"
		if execOptions, err := k8sexec.NewExecOptions(v.Namespace, v.PodName, containerName); err == nil {
			kRsync := rsyncdop.NewRsyncK8s(execOptions)
			const cerTemplate = "https://%s:3000/ta/ta.cer"
			const talTemplate = "https://%s:3000/ta/ta.tal"
			//TODO 这里有点问题，容器里shell写的有点问题，导致正常信息也输出到错误信息当中了，所以这块的错误处理暂时先留着，以后再改。
			kRsync.ConfigureCerAndTal(fmt.Sprintf(cerTemplate, v.PodName), fmt.Sprintf(talTemplate, v.PodName))
		}
	}
}

func SetUp(dataDir string) {
	var publishPoints map[string]PublishPoint
	viper.Sub("publishPoints").Unmarshal(&publishPoints)
	configureKrill(dataDir, publishPoints)
	configureRsyncd(publishPoints)
}
