package routinatorop

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	caSetUp "github.com/waterm310n/rpkiemu-go/ca/setup"
	"github.com/waterm310n/rpkiemu-go/k8sexec"
)

type RelyParty struct {
	Namespace     string `json:"namespace,omitempty"`
	PodName       string `json:"pod_name,omitempty"`
	ContainerName string `json:"container_name,omitempty"`
}

func Setup()  {
	var publishPoints map[string]caSetUp.PublishPoint
	viper.Sub("publishPoints").Unmarshal(&publishPoints) 
	const TALtemplate = "https://%s:3000/ta/ta.tal"
	tals := []string{}
	for _,v := range publishPoints {
		tals = append(tals, fmt.Sprintf(TALtemplate,v.PodName))
	}
	var relyParties map[string]RelyParty
	viper.Sub("relyParties").Unmarshal(&relyParties)
	for _,v:= range relyParties {
		if p,err := k8sexec.NewExecOptions(v.Namespace,v.PodName,v.ContainerName) ; err != nil{
			slog.Error(err.Error())
		}else{
			kRoutinator	:= NewRoutinatorK8s(p)
			kRoutinator.configRelyParty(tals...)
		}
	}
}