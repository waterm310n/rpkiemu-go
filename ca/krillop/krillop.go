package krillop

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/waterm310n/rpkiemu-go/k8sexec"
)

const CHILDREN_REQUEST_LOCATION_FILENAME = "children_request.xml"
const PARENT_RESPONSE_FILENAME = "parent_response.xml"
const PUBLISHER_REQUEST_FILENAME = "publisher_request.xml"
const REPOSITORY_RESPONSE_LOCATION_FILENAME = "repository_response.xml"

/*一个CA可以操作的所有接口*/
type CA interface {
	createHandle(handle string)                     //在当前CA中创建指定handle
	deleteHandle(handle string)                     //在当前CA中删除指定handle
	getHandles() []string                           //获取当前CA管理的所有handle名
	existHandle(handle string) bool                 //检查一个handle是否在当前CA
	getChildren(handle string) []string             //获取当前handle的所有子handle名
	getParent(handle string) []string               //获取当前handle的所有父亲handle名
	getRoaName(handle, asn string) []string         //获取当前handle指定asn的所有Roa证书名
	getCertName(handle, parentHandle string) string //获取当前handle从指定父handle获取的资源证书名

	//下面是发布点处理方法

	getRepoRequest(handle string)   //在当前CA中为指定handle创建repo请求
	setRepoConfigure(handle string) //在当前发布点中为指定handle配置发布点信息
	setPubserver(handle string)     //在当前CA中为指定handle配置它的发布点位置信息

	//下面是上下级CA资源处理方法

	getParentRequest(handle string)                                //在当前CA中为指定handle创建request请求
	setChild(handle, childHandle string, ipv4, ipv6, asn []string) //在当前CA从指定handle为下属handle分配资源
	setParent(handle, parentHandle string)                         //在当前CA为指定handle配置其上级handle

	//下面是roa发布处理方法

	addAsnIpPair(handle, ip, asn string)    //在当前CA为指定handle添加一条ASN-IP对
	removeAsnIpPair(handle, ip, asn string) //在当前CA为指定handle删除一条ASN-IP对
	addDeltaRoa(handle, file string)        //在当前CA中
}

/*krill容器操作接口，符合该接口的所有类型皆可执行下面的所有方法*/
type KrillOp interface {
	Exec(cmd string) ([]byte, error)
	Upload(srcFile string, dstFile string) error
	GetLog(limitLine int) ([]string, error)
}

/*实现了CA接口*/
type krillK8sCA struct {
	*k8sexec.ExecOptions
	isRIR bool
}

func NewKrillK8sCA(p *k8sexec.ExecOptions, isRIR bool) *krillK8sCA {
	return &krillK8sCA{
		ExecOptions: p,
		isRIR:       isRIR,
	}
}

// 更新容器中的krill配置
func (kCA *krillK8sCA) Configure() error {
	if err := kCA.Upload("tmp/"+kCA.ContainerName+".conf", "/var/krill/data/krill.conf"); err != nil {
		return err
	}
	if _, err := kCA.Exec("krill -c /var/krill/data/krill.conf"); err != nil {
		return err
	}
	return nil
}

func (kCA *krillK8sCA) createHandle(handle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc add --ca %s", handle)); err != nil && err.Error() != fmt.Sprintf("Error: CA '%s' was already initialized\n", handle) {
		slog.Error(err.Error())
	}
}

func (kCA *krillK8sCA) deleteHandle(handle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc delete --ca %s", handle)); err != nil {
		slog.Error(err.Error())
	}
}

func (kCA *krillK8sCA) existHandle(handle string) bool {
	for _, handle2 := range kCA.getHandles() {
		if handle2 == handle {
			return true
		}
	}
	return false
}

func (kCA *krillK8sCA) getHandles() []string {
	var handles []string
	if out, err := kCA.Exec(fmt.Sprintf("krillc list -f json")); err != nil {
		slog.Error(err.Error())
	} else {
		var msg struct {
			CAs []struct {
				Handle string `json:"handle,omitempty"`
			} `json:"CAs,omitempty"`
		}
		err = json.Unmarshal(out, &msg)
		if err != nil {
			slog.Error(err.Error())
		}
		for _, v := range msg.CAs {
			handles = append(handles, v.Handle)
		}
	}
	return handles
}

// 获取一个handle下的所有下级handle
func (kCA *krillK8sCA) getChildren(handle string) []string {
	if out, err := kCA.Exec(fmt.Sprintf("krillc show --ca %s -f json", handle)); err == nil {
		var msg struct {
			Children []string `json:"children,omitempty"`
		}
		err = json.Unmarshal(out, &msg)
		if err != nil {
			slog.Error(err.Error())
		}
		return msg.Children
	} else {
		//可能原因是handle不存在
		slog.Error(err.Error())
		return nil
	}
}

// //获取一个handle下的所有上级handle
func (kCA *krillK8sCA) getParent(handle string) []string {
	if out, err := kCA.Exec(fmt.Sprintf("krillc show --ca %s -f json", handle)); err == nil {
		var msg struct {
			Parents []struct {
				Handle string `json:"handle,omitempty"`
			} `json:"parents,omitempty"`
		}
		err = json.Unmarshal(out, &msg)
		if err != nil {
			slog.Error(err.Error())
		}
		parents := []string{}
		for _, v := range msg.Parents {
			parents = append(parents, v.Handle)
		}
		return parents
	} else {
		//可能原因是handle不存在
		slog.Error(err.Error())
		return nil
	}
}

// 将某个上级handle给他的下级handle的资源证书找出来
func (kCA *krillK8sCA) getCertName(handle string, parentHandle string) string {
	if out, err := kCA.Exec(fmt.Sprintf("krillc show --ca %s -f json", handle)); err == nil {
		var msg struct {
			ResourceClasses map[int]struct {
				ParentHandle string `json:"parent_handle,omitempty"`
				Keys         struct {
					Active struct {
						ActiveKey struct {
							IncomingCert struct {
								Name string `json:"name,omitempty"`
							} `json:"incoming_cert,omitempty"`
						} `json:"active_key,omitempty"`
					} `json:"active,omitempty"`
				} `json:"keys,omitempty"`
			} `json:"resource_classes,omitempty"`
		}
		err = json.Unmarshal(out, &msg)
		if err != nil {
			slog.Error(err.Error())
		}
		for _, v := range msg.ResourceClasses {
			if v.ParentHandle == parentHandle {
				return v.Keys.Active.ActiveKey.IncomingCert.Name
			}
		}
	} else {
		//可能原因是handle不存在
		slog.Error(err.Error())
	}
	return ""
}

func (kCA *krillK8sCA) getRoaName(handle string, asn string) []string {
	panic("not implemented") // TODO: Implement
}

// 下面是发布点处理方法

// 在当前CA中为指定handle创建repo请求
func (kCA *krillK8sCA) getRepoRequest(handle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc repo request --ca %s > %s", handle, PUBLISHER_REQUEST_FILENAME)); err != nil {
		slog.Error(err.Error())
	}
	// TODO: Implement 如果有需求的话在做修改，目前默认就是运行CA服务的都运行自己的仓库，自己的CA handle自己管理仓库
}

// 为发布点中添加publisher
func (kCA *krillK8sCA) setPubserver(handle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc pubserver publishers add --publisher %s --request %s > %s", handle, PUBLISHER_REQUEST_FILENAME, REPOSITORY_RESPONSE_LOCATION_FILENAME)); err != nil && err.Error() != fmt.Sprintf("Error: Duplicate publisher '%s'\n", handle) {
		slog.Error(err.Error())
	}
	// TODO: Implement 如果有需求的话在做修改，目前默认就是运行CA服务的都运行自己的仓库，自己的CA handle自己管理仓库
}

// 为handle配置发布点
func (kCA *krillK8sCA) setRepoConfigure(handle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc repo configure --ca %s --response %s", handle, REPOSITORY_RESPONSE_LOCATION_FILENAME)); err != nil && err.Error() != fmt.Sprint("Invalid RFC 8183 XML: malformed XML\n") {
		slog.Error(err.Error())
	}
	// TODO: Implement 如果有需求的话在做修改，目前默认就是运行CA服务的都运行自己的仓库，自己的CA handle自己管理仓库
}

// 下面是上下级CA资源处理方法
// 在当前CA中为指定handle创建request请求
func (kCA *krillK8sCA) getParentRequest(handle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc parents request --ca %s > %s", handle, CHILDREN_REQUEST_LOCATION_FILENAME)); err != nil {
		slog.Error(err.Error())
	}
}

func (kCA *krillK8sCA) setChild(handle string, childHandle string, ipv4 []string, ipv6 []string, asn []string) {
	cmd := []string{"krillc children add --ca", handle,
		"--child", childHandle,
		"--request", CHILDREN_REQUEST_LOCATION_FILENAME}
	fmt.Print(len(ipv4))
	if ipv4 != nil && len(ipv4) != 0 {
		cmd = append(cmd, "--ipv4")
		cmd = append(cmd, strings.Join(ipv4, ","))
	}
	if ipv6 != nil && len(ipv6) != 0 {
		cmd = append(cmd, "--ipv6")
		cmd = append(cmd, strings.Join(ipv6, ","))
	}
	if asn != nil && len(asn) != 0 {
		cmd = append(cmd, "--asn")
		cmd = append(cmd, strings.Join(asn, ","))
	}
	cmd = append(cmd, ">")
	cmd = append(cmd, PARENT_RESPONSE_FILENAME)
	if _, err := kCA.Exec(strings.Join(cmd, " ")); err != nil && err.Error() != fmt.Sprintf("Error: CA '%s' already has a child named '%s'\n", handle, childHandle) {
		slog.Error(err.Error())
	}
}

func (kCA *krillK8sCA) setParent(handle string, parentHandle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc parents add --parent %s --ca %s --response %s", parentHandle, handle, PARENT_RESPONSE_FILENAME)); err != nil && err.Error() != fmt.Sprint("ERROR Invalid RFC 8183 XML: malformed XML\n") {
		slog.Error(err.Error())
	}
}

// 下面是roa发布处理方法
// 在当前CA为指定handle添加一条ASN-IP对
func (kCA *krillK8sCA) addAsnIpPair(handle string, ip string, asn string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc roas update --add '%s => %s' --ca %s", ip, asn, handle)); err != nil {
		slog.Error(err.Error())
	}
}

func (kCA *krillK8sCA) removeAsnIpPair(handle string, ip string, asn string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc roas update --remove '%s => %s' --ca %s", ip, asn, handle)); err != nil {
		slog.Error(err.Error())
	}
}

func (kCA *krillK8sCA) addDeltaRoa(handle string, file string) {
	panic("not implemented") // TODO: Implement
}

func Create(dataDir string) {
	createKrillConfig("tmp", "ripe", true)
	op, err := k8sexec.NewExecOptions("bgp", "r5", "ripe")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	kCA := NewKrillK8sCA(op, true)
	kCA.Configure() 
	kCA.createHandle("child1")
	kCA.getRepoRequest("child1")
	kCA.setPubserver("child1")
	kCA.setRepoConfigure("child1")
	kCA.getParentRequest("child1")
	ipv6 := []string{}
	kCA.setChild("testbed", "child1", []string{"192.168.0.0/16"}, ipv6, []string{"AS5454"})
	kCA.setParent("child1", "testbed")
	kCA.addAsnIpPair("child1","192.168.0.0/16","5454")
	kCA.removeAsnIpPair("child1","192.168.0.0/16","5454")
}
