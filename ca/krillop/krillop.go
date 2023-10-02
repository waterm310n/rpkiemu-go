package krillop

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"encoding/json"
	"github.com/waterm310n/rpkiemu-go/k8sexec"
)


const CHILDREN_REQUEST_loCAtion_FILENAME = "children_request.xml"
const PARENT_RESPONSE_FILENAME = "parent_response.xml"
const PUBLISHER_REQUEST_FILENAME = "publisher_request.xml"
const REPOSITORY_RESPONSE_loCAtion_FILENAME = "repository_response.xml"

/*一个CA可以操作的所有接口*/
type CA interface{
	createHandle(handle string) //在当前CA中创建指定handle
	deleteHandle(handle string) //在当前CA中删除指定handle
	getHandles() []string //获取当前CA管理的所有handle名
	existHandle(handle string) bool //检查一个handle是否在当前CA
	getChildren(handle string) []string //获取当前handle的所有子handle名
	getParent(handle string) []string //获取当前handle的所有父亲handle名
	getRoaName(handle ,asn string) []string //获取当前handle指定asn的所有Roa证书名
	getCertName(handle,parent_handle string) string //获取当前handle从指定父handle获取的资源证书名
	
	//下面是发布点处理方法

	getRepoRequest(handle string) //在当前CA中为指定handle创建repo请求
	setRepoConfigure(handle string) //在当前发布点中为指定handle配置发布点信息
	setPubserver(handle string) //在当前CA中为指定handle配置它的发布点位置信息

	//下面是上下级CA资源处理方法

	getParentRequest(handle string) //在当前CA中为指定handle创建request请求
	setChild(handle,child_handle string,ipv4 ,ipv6 ,asn []string) //在当前CA从指定handle为下属handle分配资源
	setParent(handle,parent_handle string) //在当前CA为指定handle配置其上级handle

	//下面是roa发布处理方法

	addAsnIpPair(handle,ip,asn string) //在当前CA为指定handle添加一条ASN-IP对
	removeAsnIpPair(handle,ip,asn string) //在当前CA为指定handle删除一条ASN-IP对
	addDeltaRoa(handle ,file string) //在当前CA中
}

/*krill容器操作接口，符合该接口的所有类型皆可执行下面的所有方法*/
type KrillOp interface{
	Exec(cmd string) (map[string][]byte,error)
	Upload(srcFile string,dstFile string) error
	GetLog(limitLine int) ([]string,error)
}

/*实现了CA接口*/
type krillK8sCA struct{
	*k8sexec.ExecOptions
	token string
	isRIR bool
}

func NewKrillK8sCA(p *k8sexec.ExecOptions,isRIR bool) *krillK8sCA{
	logs, err := p.GetLog(1)
	if err != nil{
		slog.Error(err.Error())
	}
	parts := strings.Split(logs[0]," ")
	if isRIR {
		p.Upload("tmp/"+p.ContainerName+".conf","/var/krill/data/krill.conf")
	}
	return &krillK8sCA{
		ExecOptions:p,
		token :parts[len(parts)-1],
		isRIR:isRIR,
	}
}



func (kCA *krillK8sCA) createHandle(handle string) {
	if _,err :=kCA.Exec(fmt.Sprintf("krillc add --CA %s --token %s",handle,kCA.token));err != nil{
		slog.Error(err.Error())
	}
}

func (kCA *krillK8sCA) deleteHandle(handle string) {
	if _,err :=kCA.Exec(fmt.Sprintf("krillc delete --CA %s --token %s",handle,kCA.token));err != nil{
		slog.Error(err.Error())
	}
}

func (kCA *krillK8sCA) getHandles() []string {
	var handles []string
	if mp,err :=kCA.Exec(fmt.Sprintf("krillc list -f json --token %s",kCA.token));err != nil{
		slog.Error(err.Error())
	}else{
		var msg struct {
			CAs []struct {
				Handle string `json:"handle,omitempty"`
			} `json:"CAs,omitempty"`
		}
		err = json.Unmarshal(mp["stdout"], &msg)
		if err != nil{
			slog.Error(err.Error())
		}
		for _,v:= range msg.CAs{
			handles = append(handles, v.Handle)
		}
	}
	return handles
}

func (kCA *krillK8sCA) existHandle(handle string) bool {
	for _,handle2:= range kCA.getHandles(){
		if handle2 == handle{
			return true
		}
	}
	return false
}

func (kCA *krillK8sCA) getChildren(handle string) []string {
	panic("not implemented") // TODO: Implement
}

func (kCA *krillK8sCA) getParent(handle string) []string {
	panic("not implemented") // TODO: Implement
}

func (kCA *krillK8sCA) getRoaName(handle string, asn string) []string {
	panic("not implemented") // TODO: Implement
}

func (kCA *krillK8sCA) getCertName(handle string, parent_handle string) string {
	panic("not implemented") // TODO: Implement
}


// 下面是发布点处理方法
// 在当前CA中为指定handle创建repo请求
func (kCA *krillK8sCA) getRepoRequest(handle string) {
	if _, err:=kCA.Exec(fmt.Sprintf("krillc repo request --ca %s --token %s > %s",handle,kCA.token,PUBLISHER_REQUEST_FILENAME));err != nil{
		slog.Error(err.Error())
	}
}

func (kCA *krillK8sCA) setRepoConfigure(handle string) {
	panic("not implemented") // TODO: Implement
}

func (kCA *krillK8sCA) setPubserver(handle string) {
	panic("not implemented") // TODO: Implement
}

// 下面是上下级CA资源处理方法
// 在当前CA中为指定handle创建request请求
func (kCA *krillK8sCA) getParentRequest(handle string) {
	panic("not implemented") // TODO: Implement
}

func (kCA *krillK8sCA) setChild(handle string, child_handle string, ipv4 []string, ipv6 []string, asn []string) {
	panic("not implemented") // TODO: Implement
}

func (kCA *krillK8sCA) setParent(handle string, parent_handle string) {
	panic("not implemented") // TODO: Implement
}

// 下面是roa发布处理方法
// 在当前CA为指定handle添加一条ASN-IP对
func (kCA *krillK8sCA) addAsnIpPair(handle string, ip string, asn string) {
	panic("not implemented") // TODO: Implement
}

func (kCA *krillK8sCA) removeAsnIpPair(handle string, ip string, asn string) {
	panic("not implemented") // TODO: Implement
}

func (kCA *krillK8sCA) addDeltaRoa(handle string, file string) {
	panic("not implemented") // TODO: Implement
}



func Create(dataDir string){
	if err:=os.Mkdir("tmp",os.ModePerm);err!=nil && err.Error() != "mkdir tmp: file exists"{
		slog.Error(err.Error())
		os.Exit(1)		
	}
	createKrillConfig("tmp","ripe",true)
	op,err := k8sexec.NewExecOptions("bgp","r5","ripe")
	if err != nil{
		slog.Error(err.Error())
		return 
	}
	kCA := NewKrillK8sCA(op,true)
	fmt.Print(kCA.getHandles())
	// := NewKrillCA(op)
	// create_handle(op,"hello",token)
	
	// get_handles(op,token)
}

