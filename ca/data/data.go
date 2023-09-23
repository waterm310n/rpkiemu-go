package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var rirsMap = map[string]string{
	"RIPE":    "rsync://rpki.ripe.net/ta/ripe-ncc-ta.cer",
	"APNIC":   "rsync://rpki.apnic.net/repository/apnic-rpki-root-iana-origin.cer",
	"AFRINIC": "rsync://rpki.afrinic.net/repository/AfriNIC.cer",
	"LACNIC":  "rsync://repository.lacnic.net/rpki/lacnic/rta-lacnic-rpki.cer",
	"ARIN":    "rsync://rpki.arin.net/repository/arin-rpki-ta.cer",
}

var config *databaseConfig

// 字段首字母大写是因为viper模块需要调用
type databaseConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	Database   string
	Rirs       []string
	Tables     map[string]string
	Ases       []string
	LimitLayer int
}

// handle结构体
type handle struct {
	CertName     string
	PublishPoint string
	Ipv4         []string
	Ipv6         []string
	Asn          []string
}

// 返回拼接的mysql dsn url，并且返回数据库配置结构体，存在sql注入漏洞--！
func getDSN() string {
	//使用viper的反序列化，需要注意的是结构体的field名需要大写首字母
	err := viper.Sub("databaseConfig").Unmarshal(&config)
	if err != nil {
		slog.Error(err.Error())
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.Database)
	return url
}

// 连接mysql数据库并初始化相关参数，返回数据库配置结构体
func connect() *sql.DB {
	dsn := getDSN()
	slog.Info("connecting mysql server", "DSN", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.Error(err.Error())
	}
	err = db.Ping()
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("connect successfully")
	db.SetConnMaxLifetime(time.Minute * 3) // 时间建议小于5mins
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

// 将ipv4和ipv6从字符串中分离
func splitIpType(ipResources string) (ipv4 []string, ipv6 []string) {
	ipv4 = make([]string, 0)
	ipv6 = make([]string, 0)
	for _, ip := range strings.Split(ipResources, ",") {
		if strings.Index(ip, ".") != -1 { //有'.'
			ipv4 = append(ipv4, ip)
		} else {
			ipv6 = append(ipv6, ip)
		}
	}
	return ipv4, ipv6
}

// 检查给定的Asn是否在某个证书划分的as资源中
func checkAsn(parts []string) bool {
	var intervals []struct {
		min int
		max int
	}
	for _, part := range parts {
		part = strings.Trim(part, " ")
		if part == "" {
			continue
		}
		part := strings.Split(part, "-")
		min, err := strconv.Atoi(part[0][strings.Index(part[0], "AS")+2:])
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		var max int
		if len(part) == 1 {
			max = min
		} else {
			max, err = strconv.Atoi(part[1][strings.Index(part[1], "AS")+2:])
			if err != nil {
				slog.Error(err.Error())
				continue
			}
		}
		intervals = append(intervals, struct {
			min int
			max int
		}{
			min,
			max,
		})
	}
	for _, asn := range config.Ases {
		asn, err := strconv.Atoi(asn[2:])
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		for _, interval := range intervals {
			if asn >= interval.min && asn <= interval.max {
				return true
			}
		}
	}
	return false
}
// 从ipResources，asResources，uri中提取信息构造handle
func preProcessing(ipResources, asResources, uri string) *handle {
	asResources = ASN_MATCH.ReplaceAllString(asResources, "AS$1")
	asResources = ASN_MINMAX_MATCH.ReplaceAllString(asResources, "$1-$2")
	asn := []string{}
	if  asResources != "[ ]"{
		asn = strings.Split(asResources[1:len(asResources)-3], ",")
	}
	if ipResources != "[ ]"{
		ipResources = IPV4_MINMAX_MATCH.ReplaceAllString(ipResources[1:len(ipResources)-3], "$1-$2")
		ipResources = IPV6_MINMAX_MATCH.ReplaceAllString(ipResources, "$1-$2")
	}	
	ipv4, ipv6 := splitIpType(ipResources)
	parts := URI_MATCH.FindStringSubmatch(uri)
	publishPoint := strings.Split(parts[2], ".")[len(strings.Split(parts[2], "."))-2]
	return &handle{
		Ipv4:         ipv4,
		Ipv6:         ipv6,
		Asn:          asn,
		PublishPoint: publishPoint,
		CertName:     "m" + parts[3],
	}
}

// 对给定的aia执行查找语句并将结果输出到以_roas结尾的文件中
func writeRoas(roasStmt *sql.Stmt, aia string,path string){
	rows,err:=roasStmt.Query(aia)
	if err != nil {
		slog.Error(err.Error())
	}
	helper := func(Ipaddrblocks,AsId string) []string {
		var f interface{}
		res := []string{}
		//TODO用接口写好麻烦，以后用结构体反序列化
		json.Unmarshal([]byte(Ipaddrblocks),&f)
		if ipAddrBlocks,ok := f.(map[string]interface{})["ipAddrBlocks"];ok{
			if ipAddrBlocks,ok:=ipAddrBlocks.([]interface{});ok{
				for _,ipAddrBlock:= range ipAddrBlocks{
					if addresses,ok:=ipAddrBlock.(map[string]interface{})["addresses"];ok{
						if addresses,ok:=addresses.([]interface{});ok{
							for _,address := range addresses{
								if address,ok := address.(map[string]interface{})["address"];ok{
									if address,ok := address.(string);ok{
										res = append(res, "A: "+address+" => "+AsId)
									}
								}
							}
						}
					}
				}
			}
		}
		return res
	}
	res := []string{}
	flag := true
	var file *os.File
	for rows.Next() {
		if flag{
			file,err = os.Create(path)
			if err != nil{
				slog.Error(err.Error())
				return
			}
			flag=false
		}
		var Ipaddrblocks,AsId string
		rows.Scan(&Ipaddrblocks,&AsId)
		res = append(res,helper(Ipaddrblocks,AsId)...)
	}
	file.WriteString(strings.Join(res,"\n"))
}

// 深度优先搜索层次结构生成roas数据文件和层次结构数据文件
func dfsHierarchy(hierarchyStmt *sql.Stmt, roasStmt *sql.Stmt, dataDir string, depth int) {
	var helper func(string, string, int)
	helper = func(dataDir string, aia string, depth int) {
		if depth == 0 {
			return
		}
		rows, err := hierarchyStmt.Query(aia)
		if err != nil {
			slog.Error(err.Error())
		}
		flag := true
		for rows.Next() {
			if flag {
				if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
					slog.Error(err.Error())
				}
				flag = false
			}
			var IPResources, ASResources, URI string
			rows.Scan(&IPResources, &ASResources, &URI)
			handle := preProcessing(IPResources, ASResources, URI)
			if checkAsn(handle.Asn) {
				file, err := os.Create(dataDir + "/" + handle.CertName)
				if err != nil {
					slog.Error(err.Error())
				}
				if content, err := json.Marshal(&handle); err == nil {
					file.Write(content)
					helper(dataDir+"/"+handle.CertName+"_children", URI, depth-1)
				}
				writeRoas(roasStmt,URI,dataDir+"/"+handle.CertName+"_roas")
			}
		}
	}
	rirs := viper.GetStringSlice("databaseConfig.rirs")
	for _, v := range rirs {
		parts := URI_MATCH.FindStringSubmatch(rirsMap[v])
		certName := "m" + parts[3]
		file, err := os.Create(dataDir + "/" + certName)
		if err != nil {
			slog.Error(err.Error())
		}
		if content, err := json.Marshal(&handle{
			Ipv4:         []string{"0.0.0.0/0"},
			Ipv6:         []string{"::/0"},
			Asn:          []string{"AS0-AS4294967295"},
			CertName:     certName,
			PublishPoint: strings.Split(parts[2], ".")[len(strings.Split(parts[2], "."))-2],
		}); err == nil {
			file.Write(content)
			helper(dataDir+"/"+certName+"_children", rirsMap[v], depth-1)
		}
	}
}

func GenerateData(dataDir string) {
	db := connect()
	defer db.Close()
	prepareHierarchyStmt := fmt.Sprintf("select IPResources, ASResources,URI from %s.%s where aia = ? and isvalid = 1", config.Database, config.Tables["cas"])
	hierarchyStmt, err := db.Prepare(prepareHierarchyStmt)
	defer hierarchyStmt.Close()
	if err != nil {
		slog.Error(err.Error())
	}
	prepareRoasStmt := fmt.Sprintf("select Ipaddrblocks,AsId from %s.%s where aia= ? and isvalid=1",config.Database, config.Tables["roas"])
	roasStmt, err := db.Prepare(prepareRoasStmt)
	defer roasStmt.Close()
	if err != nil {
		slog.Error(err.Error())
	}
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		slog.Error(err.Error())
	}
	dfsHierarchy(hierarchyStmt, roasStmt,dataDir, config.LimitLayer)
	slog.Info("generate data compelete")
}
