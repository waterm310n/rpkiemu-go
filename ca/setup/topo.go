// 下面的代码是bgpemu项目下topo.yaml的部分使用到的结构
// 主要用于本项目反序列化topo.yaml，以向其中插入本项目所需要其支持的内容

package setup

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Link struct {
	AInt  string `yaml:"a_int,omitempty"`
	ANode string `yaml:"a_node,omitempty"`
	ZInt  string `yaml:"z_int,omitempty"`
	ZNode string `yaml:"z_node,omitempty"`
}

type VolumeType int32

const (
	VolumeType_DEFAULT  VolumeType = 0
	VolumeType_EMPTY    VolumeType = 1
	VolumeType_HOSTPATH VolumeType = 2
)

// 将VolumeType的值与有意义的字符串进行互相转换的map
// Enum value maps for VolumeType.
var (
	VolumeType_name = map[int32]string{
		0: "DEFAULT",
		1: "EMPTY",
		2: "HOSTPATH",
	}
	VolumeType_value = map[string]int32{
		"DEFAULT":  0,
		"EMPTY":    1,
		"HOSTPATH": 2,
	}
)

// 这里与MarshalJson不同，要使用非指针接收者，可以参考https://github.com/go-yaml/yaml/issues/714
func (t VolumeType) MarshalYAML() (interface{}, error) {
	if v, ok := VolumeType_name[int32(t)]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("can not encode %v while marshal VolumeType in topo.go func (VolumeType) MarshalYAML()", t)
	}
}

// 由于map映射的关系，默认情况下VolumeType的值为字符串，因此需要进行最小转化处理
func (t *VolumeType) UnmarshalYAML(node *yaml.Node) error {
	if v, ok := VolumeType_value[node.Value]; ok {
		*t = VolumeType(v)
	} else {
		return fmt.Errorf("can not decode %v while unmarshal VolumeType in topo.go func (*VolumeType) UnmarshalYAML(yaml.Node)", node.Value)
	}
	return nil
}

type Type int32

const (
	Type_UNKNOWN Type = 0
	Type_HOST    Type = 1
	Type_BGP     Type = 2
	Type_SUBTOPO Type = 3
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "UNKNOWN",
		1: "HOST",
		2: "BGP",
		3: "SUBTOPO",
	}
	Type_value = map[string]int32{
		"UNKNOWN": 0,
		"HOST":    1,
		"BGP":     2,
		"SUBTOPO": 3,
	}
)

// 这里与MarshalJson不同，要使用非指针接收者，可以参考https://github.com/go-yaml/yaml/issues/714
func (t Type) MarshalYAML() (interface{}, error) {
	if v, ok := Type_name[int32(t)]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("can not encode %v while marshal Type in topo.go func (t Type) MarshalYAML()", t)
	}
}

// 由于map映射的关系，默认情况下Type的值为字符串，因此需要进行最小转化处理
func (t *Type) UnmarshalYAML(node *yaml.Node) error {
	if v, ok := Type_value[node.Value]; ok {
		*t = Type(v)
	} else {
		return fmt.Errorf("can not decode %v while unmarshal Type in topo.go func (*Type) UnmarshalYAML(yaml.Node)", node.Value)
	}
	return nil
}

type ShareVolume struct {
	Type VolumeType `yaml:"type,omitempty"`
}

type PublicVolumes struct {
	Volumes map[string]string `yaml:"volumes,omitempty"`
}

type Task struct {
	/*
		下面这个如果标注omitempty的话，序列化后如果Cmds为空，那么将直接省略Cmds，
		因此可能不符合bgpemu的要求，故没有忽略
	*/
	Cmds []string `yaml:"cmds"`
	Container string `yaml:"container,omitempty"`
}

type Config struct {
	IsResilient      bool                      `yaml:"is_resilient,omitempty"`
	ContainerVolumes map[string]*PublicVolumes `yaml:"container_volumes,omitempty"`
	ExtraImages      map[string]string         `yaml:"extra_images,omitempty"`
	ShareVolumes     map[string]*ShareVolume   `yaml:"share_volumes,omitempty"`
	Tasks            []*Task                   `yaml:"tasks,omitempty"`
}

type Service struct {
	Inside    uint32 `yaml:"inside,omitempty"`
	Outside   uint32 `yaml:"outside,omitempty"`
	InsideIp  string `yaml:"inside_ip,omitempty"`
	OutsideIp string `yaml:"outside_ip,omitempty"`
	NodePort  uint32 `yaml:"node_port,omitempty"`
	Name      string `yaml:"name,omitempty"`
}

type Node struct {
	Config   *Config             `yaml:"config,omitempty"`
	IpAddr   map[string]string   `yaml:"ip_addr,omitempty"`
	Name     string              `yaml:"name,omitempty"`
	Services map[uint32]*Service `yaml:"services,omitempty"`
	Type     Type                `yaml:"type,omitempty"`
}

type Topology struct {
	Links []*Link `yaml:"links,omitempty"`
	Name  string  `yaml:"name,omitempty"`
	Nodes []*Node `yaml:"nodes,omitempty"` // List of nodes in the topology
}
