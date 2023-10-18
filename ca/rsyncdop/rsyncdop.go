package rsyncdop

import (
	"log/slog"
	"github.com/waterm310n/rpkiemu-go/k8sexec"
)

type RsyncdOp interface {
	configureCerAndTal() error
}

type RsyncK8s struct {
	*k8sexec.ExecOptions
}

func NewRsyncK8s(p *k8sexec.ExecOptions) *RsyncK8s {
	return &RsyncK8s{
		ExecOptions: p,
	}
}
func (kRsync *RsyncK8s) ConfigureCerAndTal(args ...string) error{
	cmd := "/opt/entrypoint.sh "
	for _, arg := range args {
		cmd = cmd + " " + arg
	}
	if _, err := kRsync.Exec(cmd); err != nil {
		slog.Debug(err.Error(),"cmd ",cmd,"container",kRsync.ContainerName)
		return err
	}
	return nil
}
