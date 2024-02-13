package routinatorop

import (
	"log/slog"
	"github.com/waterm310n/rpkiemu-go/k8sexec"
)


type RoutintorOp interface {
	configRelyParty() error
}

type RoutinatorK8s struct {
	*k8sexec.ExecOptions
}

func NewRoutinatorK8s(p *k8sexec.ExecOptions) *RoutinatorK8s {
	return &RoutinatorK8s{
		ExecOptions: p,
	}
}

//在routinator容器中执行 /opt/entryPoint.sh 
func (kRoutinator *RoutinatorK8s) configRelyParty(args... string) error {
	cmd := "/opt/entrypoint.sh "
	for _, arg := range args {
		cmd = cmd + " " + arg
	}
	if _,err := kRoutinator.Exec(cmd) ; err != nil{
		slog.Debug(err.Error())
		return err
	}
	return nil
}