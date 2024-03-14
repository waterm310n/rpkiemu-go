package setup

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

/*
为ca容器，创建krill.conf文件，完成一些前置准备
*/

const KRILL_TEMPLATE = `admin_token = "krillTestBed"
data_dir = "/var/krill/data/"
log_type = "stderr"
ip = "0.0.0.0"
service_uri = "https://%s:3000/"
bgp_risdumps_enabled = false
`

const KRILL_TESTBED_TEMPLATE = `admin_token = "krillTestBed"
data_dir = "/var/krill/data/"
log_type = "stderr"
ip = "0.0.0.0"
service_uri = "https://%s:3000/" 
bgp_risdumps_enabled = false

[testbed]
rrdp_base_uri = "https://%s:3000/rrdp/"
rsync_jail = "rsync://%s/repo/"
ta_uri = "https://%s:3000/ta/ta.cer"
ta_aia = "rsync://%s/ta/ta.cer"
`
const HOST2IP_TEMPLATE = `echo %s %s >> /etc/hosts`

const KRILL_IMAGE = "krill:local"
const RSYNCD_IMAGE = "rsyncd:local"
const RELYPARTY_IMAGE = "routinator:local"

// 用于制作配置
type RelyParty struct {
	Namespace     string `json:"namespace,omitempty"`
	PodName       string `json:"pod_name,omitempty"`
	ContainerName string `json:"container_name,omitempty"`
}

// 发布点与依赖方组合
type PointConfigure struct {
	PublishPoints map[string]*PublishPoint `json:"publish_points,omitempty"`
	RelyParties   map[string]*RelyParty    `json:"rely_parties,omitempty"`
}

// 在tmp目录下创建name.conf
func createKrillConfig(dataDir, name string, isRIR bool) error {
	var err error
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	file, err := os.OpenFile(dataDir+"/"+name+".conf", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer file.Close()
	if isRIR {
		_, err = file.WriteString(fmt.Sprintf(KRILL_TESTBED_TEMPLATE, name, name, name, name, name))
	} else {
		_, err = file.WriteString(fmt.Sprintf(KRILL_TEMPLATE, name))
	}
	return err
}

func parsePublishPointJson(publishPointsJson string) (map[string]*PublishPoint, error) {
	var pointConfigure PointConfigure
	var publishPoints map[string]*PublishPoint
	if bytes, err := os.ReadFile(publishPointsJson); err != nil {
		return nil, err
	} else {
		json.Unmarshal(bytes, &pointConfigure)
		publishPoints = pointConfigure.PublishPoints
	}
	return publishPoints, nil
}

func parseTopoYaml(topoYaml string) (*Topology, error) {
	var topo *Topology
	if bytes, err := os.ReadFile(topoYaml); err != nil {
		return nil, err
	} else {
		if err := yaml.Unmarshal(bytes, &topo); err != nil {
			return nil, err
		}
	}
	return topo, nil
}

func configContainerVolumes(config *Config, publishPoint *PublishPoint) {
	RPKI_VOLUME_NAME := "rpki"
	KRILL_VOLUME_PATH := "/var/krill/data/repo/rsync"
	RSYNCD_VOLUME_PATH := "/share"
	config.ContainerVolumes[publishPoint.CAContainerName] = &PublicVolumes{Volumes: map[string]string{RPKI_VOLUME_NAME: KRILL_VOLUME_PATH}}
	config.ContainerVolumes[publishPoint.RSYNCDContainerName] = &PublicVolumes{Volumes: map[string]string{RPKI_VOLUME_NAME: RSYNCD_VOLUME_PATH}}
	config.ShareVolumes[RPKI_VOLUME_NAME] = &ShareVolume{VolumeType_EMPTY}
}

func configExtraImages(config *Config, publishPoint *PublishPoint) {
	config.ExtraImages[publishPoint.CAContainerName] = KRILL_IMAGE
	config.ExtraImages[publishPoint.RSYNCDContainerName] = RSYNCD_IMAGE
}

func getPodName(topo *Topology, count int) string {
	if topo.Nodes[count].Config.IsResilient {
		return topo.Nodes[count].Name + "-0"
	} else {
		return topo.Nodes[count].Name
	}
}

// 从CIDR
func extractIp(IpAddr string) string{
	const MASK_REGEX = `\/(3[0-2]|[1-2][0-9]|[1-9])`
	MASK_MATCH := regexp.MustCompile(MASK_REGEX)
	s1 := MASK_MATCH.FindString(IpAddr)
	IpAddr = strings.ReplaceAll(IpAddr,s1,"")
	return IpAddr
}

// 这部分如果用viper来写，可能可以简单？
// 将CA方和依赖方随机分布在bgp topo提供的网络中
func SequentialAdapt(topoYaml, publishPointsJson, topoWithRPKIYaml string, rpCount int) {
	publishPoints, err := parsePublishPointJson(publishPointsJson)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	topo, err := parseTopoYaml(topoYaml)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if len(publishPoints)+rpCount > len(topo.Nodes) {
		slog.Error("The number of pods is insufficient.", "publishPoints", len(publishPoints), "rpCount", rpCount, "Pods", len(topo.Nodes))
		return
	}
	count := 0
	relyParties := make(map[string]*RelyParty)
	//处理依赖方
	for i := 0; i < rpCount; i++ {
		config := topo.Nodes[count].Config
		podName := getPodName(topo, count)
		containerName := podName + "-routinator"
		relyParties[containerName] = &RelyParty{Namespace: topo.Name, PodName: podName, ContainerName: containerName}
		config.ExtraImages[containerName] = RELYPARTY_IMAGE
		config.Tasks = append(config.Tasks, &Task{Container: containerName, Cmds: make([]string, 0)})
		count++
	}
	//处理CA方
	for name, publishPoint := range publishPoints {
		config := topo.Nodes[count].Config
		podName := getPodName(topo, count)
		publishPoint.Namespace = topo.Name
		publishPoint.PodName = podName
		publishPoint.CAContainerName = podName + "-" + name
		publishPoint.RSYNCDContainerName = podName + "-rsyncd"
		if publishPoint.IsRIR {
			configContainerVolumes(config, publishPoint)
		}
		configExtraImages(config, publishPoint)
		config.Tasks = append(config.Tasks, &Task{Container: publishPoint.CAContainerName, Cmds: make([]string, 0)})
		config.Tasks = append(config.Tasks, &Task{Container: publishPoint.RSYNCDContainerName, Cmds: make([]string, 0)})
		cmdToCur := fmt.Sprintf(HOST2IP_TEMPLATE, extractIp(topo.Nodes[count].IpAddr["eth1"]), podName)
		//向各个依赖方插入CA方的HOST到IP地址的映射
		for i := 0; i < rpCount; i++ {
			tasks := topo.Nodes[i].Config.Tasks
			cmds := &tasks[len(tasks)-1].Cmds
			*cmds = append(*cmds, cmdToCur)
			cmdToOther := fmt.Sprintf(HOST2IP_TEMPLATE, extractIp(topo.Nodes[i].IpAddr["eth1"]), getPodName(topo, i))
			cmds = &config.Tasks[len(config.Tasks)-1].Cmds
			*cmds = append(*cmds, cmdToOther)
			cmds = &config.Tasks[len(config.Tasks)-2].Cmds
			*cmds = append(*cmds, cmdToOther)
		}
		//向各个CA方插入其他CA方的HOST到IP地址的映射
		for i := rpCount; i < count; i++ {
			cmdToOther := fmt.Sprintf(HOST2IP_TEMPLATE,  extractIp(topo.Nodes[i].IpAddr["eth1"]), getPodName(topo, i))
			cmds := &config.Tasks[len(config.Tasks)-1].Cmds
			*cmds = append(*cmds, cmdToOther)
			cmds = &config.Tasks[len(config.Tasks)-2].Cmds
			*cmds = append(*cmds, cmdToOther)
			tasks := topo.Nodes[i].Config.Tasks
			cmds = &tasks[len(tasks)-1].Cmds
			*cmds = append(*cmds, cmdToCur)
			cmds = &tasks[len(tasks)-2].Cmds
			*cmds = append(*cmds, cmdToCur)
		}
		count++
	}
	file, err := os.Create(topoWithRPKIYaml)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer file.Close()
	if content, err := yaml.Marshal(&topo); err == nil {
		file.Write(content)
	} else {
		slog.Error(err.Error())
	}
	file, err = os.OpenFile(publishPointsJson, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer file.Close()
	pointConfigure := PointConfigure{publishPoints, relyParties}
	if content, err := json.Marshal(&pointConfigure); err == nil {
		if _, err := file.Write(content); err != nil {
			slog.Error(err.Error())
			return
		}
	} else {
		slog.Error(err.Error())
		return
	}
}
