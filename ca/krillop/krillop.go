package krillop

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/waterm310n/rpkiemu-go/ca/data"
	"github.com/waterm310n/rpkiemu-go/k8sexec"
)

const CHILDREN_REQUEST_LOCATION_FILENAME = "children_request.xml"
const PARENT_RESPONSE_FILENAME = "parent_response.xml"
const PUBLISHER_REQUEST_FILENAME = "publisher_request.xml"
const REPOSITORY_RESPONSE_LOCATION_FILENAME = "repository_response.xml"

/*krill容器操作接口，符合该接口的所有类型皆可执行下面的所有方法*/
type KrillOp interface {
	Exec(cmd string) ([]byte, error)
	Upload(srcFile string, dstFile string) error
	GetLog(limitLine int) ([]string, error)
}

/*实现了CA接口*/
type KrillK8sCA struct {
	*k8sexec.ExecOptions
	isRIR bool
}

func NewKrillK8sCA(p *k8sexec.ExecOptions, isRIR bool) *KrillK8sCA {
	return &KrillK8sCA{
		ExecOptions: p,
		isRIR:       isRIR,
	}
}

// 更新容器中的krill配置
func (kCA *KrillK8sCA) Configure(dataDir string) error {
	if err := kCA.Upload(filepath.Join(dataDir, kCA.PodName+".conf"), "/var/krill/data/krill.conf"); err != nil {
		return err
	}
	//这条命令存在问题
	if _, err := kCA.Exec("krill -c /var/krill/data/krill.conf"); err != nil {
		return err
	}
	return nil
}

func (kCA *KrillK8sCA) createHandle(handle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc add --ca %s", handle)); err != nil && err.Error() != fmt.Sprintf("Error: CA '%s' was already initialized\n", handle) {
		slog.Error(err.Error())
	}
	path := filepath.Join("/tmp", handle)
	if _, err := kCA.Exec(fmt.Sprintf("mkdir %s", path)); err != nil && err.Error() != fmt.Sprintf("mkdir: can't create directory '%s': File exists\n", path) {
		slog.Error(err.Error())
	}
}

func (kCA *KrillK8sCA) deleteHandle(handle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc delete --ca %s", handle)); err != nil {
		slog.Error(err.Error())
	}
}

func (kCA *KrillK8sCA) existHandle(handle string) bool {
	for _, handle2 := range kCA.getHandles() {
		if handle2 == handle {
			return true
		}
	}
	return false
}

func (kCA *KrillK8sCA) getHandles() []string {
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

/**/
func (kCA *KrillK8sCA) getFile(srcFile, dstFile string) error {
	//TODO
	var err error
	kCA.Download(srcFile, dstFile)
	return err
}

// 获取一个handle下的所有下级handle
func (kCA *KrillK8sCA) getChildren(handle string) []string {
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
func (kCA *KrillK8sCA) getParent(handle string) []string {
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
func (kCA *KrillK8sCA) getCertName(handle string, parentHandle string) (map[string]interface{}, error) {
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
		if err = json.Unmarshal(out, &msg); err != nil {
			return nil, err
		}
		certNames := map[string]interface{}{}
		for _, v := range msg.ResourceClasses {
			if v.ParentHandle == parentHandle {
				certNames[v.Keys.Active.ActiveKey.IncomingCert.Name] = struct{}{}
			}
		}
		return certNames, nil
	} else {
		//可能原因是handle不存在
		return nil, err
	}
}

func (kCA *KrillK8sCA) getRoaName(handle string, asn int) (map[string]interface{}, error) {
	if content, err := kCA.Exec(fmt.Sprintf("krillc roas list --ca %s --format json", handle)); err == nil {
		roas := make([]struct {
			Asn        int `json:"asn,omitempty"`
			RoaObjects []struct {
				Uri string `json:"uri,omitempty"`
			} `json:"roa_objects,omitempty"`
		}, 0)
		if err := json.Unmarshal(content, &roas); err == nil {
			res := map[string]interface{}{}
			for _, roa := range roas {
				if roa.Asn == asn {
					uri := strings.Split(roa.RoaObjects[0].Uri, "/")
					res[uri[len(uri)-1]] = struct{}{}
				}
			}
			return res, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// 下面是发布点处理方法

// 在当前CA中为指定handle创建repo请求
func (kCA *KrillK8sCA) getRepoRequest(handle string) error {
	publishRequestFileName := filepath.Join("/tmp", fmt.Sprintf("%s_%s", handle, PUBLISHER_REQUEST_FILENAME))
	if _, err := kCA.Exec(fmt.Sprintf("krillc repo request --ca %s > %s", handle, publishRequestFileName)); err != nil {
		return err
	}
	return nil
}

// 为发布点中添加publisher
func (kCA *KrillK8sCA) setPubserver(handle string) error {
	publishRequestFileName := filepath.Join("/tmp", fmt.Sprintf("%s_%s", handle, PUBLISHER_REQUEST_FILENAME))
	repositoryResponseLocationFileName := filepath.Join("/tmp", fmt.Sprintf("%s_%s", handle, REPOSITORY_RESPONSE_LOCATION_FILENAME))
	if _, err := kCA.Exec(fmt.Sprintf("krillc pubserver publishers add --publisher %s --request %s > %s", handle, publishRequestFileName, repositoryResponseLocationFileName)); err != nil && err.Error() != fmt.Sprintf("Error: Duplicate publisher '%s'\n", handle) {
		return err
	}
	return nil
}

// 为handle配置发布点
func (kCA *KrillK8sCA) setRepoConfigure(handle string) error {
	repositoryResponseLocationFileName := filepath.Join("/tmp", fmt.Sprintf("%s_%s", handle, REPOSITORY_RESPONSE_LOCATION_FILENAME))
	if _, err := kCA.Exec(fmt.Sprintf("krillc repo configure --ca %s --response %s", handle, repositoryResponseLocationFileName)); err != nil && err.Error() != "Invalid RFC 8183 XML: malformed XML\n" {
		return err
	}
	return nil
}

// 下面是上下级CA资源处理方法
// 在当前CA中为指定handle创建request请求
func (kCA *KrillK8sCA) getParentRequest(handle string) (string, error) {
	childrenRequestLocationFileName := filepath.Join("/tmp", fmt.Sprintf("%s_%s", handle, CHILDREN_REQUEST_LOCATION_FILENAME))
	if _, err := kCA.Exec(fmt.Sprintf("krillc parents request --ca %s > %s", handle, childrenRequestLocationFileName)); err != nil {
		return "", err
	}
	return childrenRequestLocationFileName, nil
}

func (kCA *KrillK8sCA) setChild(handle string, childHandle string, ipv4 []string, ipv6 []string, asn []string) (string, error) {
	childrenRequestLocationFileName := filepath.Join("/tmp", fmt.Sprintf("%s_%s", childHandle, CHILDREN_REQUEST_LOCATION_FILENAME))
	parentResponseFileName := filepath.Join("/tmp", fmt.Sprintf("%s_%s", childHandle, PARENT_RESPONSE_FILENAME))
	cmd := []string{"krillc children add --ca", handle,
		"--child", childHandle,
		"--request", childrenRequestLocationFileName}
	if len(ipv4) != 0 {
		cmd = append(cmd, "--ipv4")
		cmd = append(cmd, strconv.Quote(strings.Join(ipv4, ",")))
	}
	if len(ipv6) != 0 {
		cmd = append(cmd, "--ipv6")
		cmd = append(cmd, strconv.Quote(strings.Join(ipv6, ",")))
	}
	if len(asn) != 0 {
		cmd = append(cmd, "--asn")
		cmd = append(cmd, strconv.Quote(strings.Join(asn, ",")))
	}
	cmd = append(cmd, ">")
	cmd = append(cmd, parentResponseFileName)
	cmdStr := strings.Join(cmd, " ")
	if _, err := kCA.Exec(cmdStr); err != nil &&
		err.Error() != fmt.Sprintf("Error: CA '%s' already has a child named '%s'\n", handle, childHandle) &&
		err.Error() != fmt.Sprintf("Error: Child '%s' cannot have resources not held by CA %s'\n", childHandle, handle) {
		slog.Error(err.Error(), "cmd", cmdStr)
		return "", err
	}
	return parentResponseFileName, nil
}

func (kCA *KrillK8sCA) updateChild(handle string, childHandle string, ipv4 []string, ipv6 []string, asn []string) error {
	cmd := []string{"krillc children update --ca", handle,
		"--child", childHandle}
	if len(ipv4) != 0 {
		cmd = append(cmd, "--ipv4")
		cmd = append(cmd, strconv.Quote(strings.Join(ipv4, ",")))
	}
	if len(ipv6) != 0 {
		cmd = append(cmd, "--ipv6")
		cmd = append(cmd, strconv.Quote(strings.Join(ipv6, ",")))
	}
	if len(asn) != 0 {
		cmd = append(cmd, "--asn")
		cmd = append(cmd, strconv.Quote(strings.Join(asn, ",")))
	}
	cmdStr := strings.Join(cmd, " ")
	if _, err := kCA.Exec(cmdStr); err != nil {
		slog.Error(err.Error(), "cmd", cmdStr)
		return err
	}
	return nil
}

func (kCA *KrillK8sCA) setParent(handle string, parentHandle string) error {
	parentResponseFileName := filepath.Join("/tmp", fmt.Sprintf("%s_%s", handle, PARENT_RESPONSE_FILENAME))
	if _, err := kCA.Exec(fmt.Sprintf("krillc parents add --parent %s --ca %s --response %s", parentHandle, handle, parentResponseFileName)); err != nil && err.Error() != "Invalid RFC 8183 XML: malformed XML\n" {
		return err
	}
	return nil
}

// 下面是roa发布处理方法
// 在当前CA为指定handle添加一条ASN-IP对
func (kCA *KrillK8sCA) AddAsnIpPair(handle string, ip string, asn string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc roas update --add '%s => %s' --ca %s", ip, asn, handle)); err != nil {
		slog.Error(err.Error())
	}
}

func (kCA *KrillK8sCA) RemoveAsnIpPair(handle string, ip string, asn string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc roas update --remove '%s => %s' --ca %s", ip, asn, handle)); err != nil {
		slog.Error(err.Error())
	}
}

func (kCA *KrillK8sCA) addDeltaRoa(handle string, file string) {
	dstFile := filepath.Join("/tmp", filepath.Base(file))
	if err := kCA.Upload(file, dstFile); err != nil {
		slog.Error(err.Error())
	}
	if _, err := kCA.Exec(fmt.Sprintf("krillc roas update --delta %s --ca %s", dstFile, handle)); err != nil {
		slog.Error(err.Error())
	}
}

func (kCA *KrillK8sCA) DeleteRoa(handle string, asn int) {
	roasName, err := kCA.getRoaName(handle, asn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	for roaName := range roasName {
		kCA.Exec(fmt.Sprintf("rm /var/krill/data/repo/rsync/current/%s/0/%s", handle, roaName))
	}
}

func (kCA *KrillK8sCA) DeleteCert(handle string, parentHandle string, parentCa *KrillK8sCA) {
	certNames, err := kCA.getCertName(handle, parentHandle)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	for certName := range certNames {
		parentCa.Exec(fmt.Sprintf("rm /var/krill/data/repo/rsync/current/%s/0/%s", parentHandle, certName))
	}
}

func (kCA *KrillK8sCA) CorruptCert(handle string, parentHandle string, parentCa *KrillK8sCA) {
	certNames, err := kCA.getCertName(handle, parentHandle)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	for certName := range certNames {
		parentCa.Exec(fmt.Sprintf("echo 'corrupted corrupted' >> /var/krill/data/repo/rsync/current/%s/0/%s", parentHandle, certName))
	}
}

func (kCA *KrillK8sCA) CorruptRoa(handle string, asn int) {
	roasName, err := kCA.getRoaName(handle, asn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	for roaName := range roasName {
		kCA.Exec(fmt.Sprintf("echo 'corrupted corrupted' >> /var/krill/data/repo/rsync/current/%s/0/%s", handle, roaName))
	}
}

func (kCA *KrillK8sCA) Revocate(handle string, parentHandle string) {
	if _, err := kCA.Exec(fmt.Sprintf("krillc children remove --child %s --ca %s", handle, parentHandle)); err != nil {
		slog.Error(err.Error())
	}
}

func (kCA *KrillK8sCA) Modificate(handle string, parentHandle string, ipv4, ipv6, asn []string) {
	if err := kCA.updateChild(parentHandle, handle, ipv4, ipv6, asn); err != nil {
		slog.Error("Modificate attack error")
	}
	/*
		需要注意，执行后并不能够马上更新，krill每十分钟才会让所有children去向parent请求新的资源
		要想马上生效，还需执行krill bulk refresh或其他相关指令
	*/
	if _, err := kCA.Exec("krillc bulk refresh"); err != nil {
		slog.Error("failed to excute 'krillc bulk refresh'")
	}
}

func (kCA *KrillK8sCA) Inject(certName string, parentHandle, publishPoint, parentPublishPoint string, ipv4, ipv6, asn []string, caOps map[string]*KrillK8sCA) {
	handle := data.Handle{
		CertName:     certName,
		Ipv4:         ipv4,
		Ipv6:         ipv6,
		Asn:          asn,
		PublishPoint: publishPoint,
	}
	if publishPoint == parentPublishPoint {
		caOps[publishPoint].createHandle(certName)
		if err := setRepo(publishPoint, certName, caOps); err != nil {
			slog.Error(err.Error())
		}
		if err := setParentChildrenRel(publishPoint, publishPoint, certName, parentHandle, handle, caOps); err != nil {
			slog.Error(err.Error())
		}
	} else {
		slog.Debug(publishPoint)
		caOps[publishPoint].createHandle(certName)
		if err := setRepo(publishPoint, certName, caOps); err != nil {
			slog.Error(err.Error())
		}
	}
}
