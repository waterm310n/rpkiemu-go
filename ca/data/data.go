package data

import (
	"database/sql"
	"log/slog"	
	"time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var rirs_map = map[string]string{
	"RIPE":"rsync://rpki.ripe.net/ta/ripe-ncc-ta.cer",
	"APNIC": "rsync://rpki.apnic.net/repository/apnic-rpki-root-iana-origin.cer",
	"AFRINIC": "rsync://rpki.afrinic.net/repository/AfriNIC.cer",
	"LACNIC": "rsync://repository.lacnic.net/rpki/lacnic/rta-lacnic-rpki.cer",
	"ARIN": "rsync://rpki.arin.net/repository/arin-rpki-ta.cer",
}

//字段首字母大写是因为viper模块需要调用
type databaseConfig struct{
	Host string
	Port int
	User string
	Password string
	Database string
	Rirs []string 
	Tables map[string]string
	Ases []string
	LimitLayer int
}

type parentHandle struct{
	ipv4 []string
	ipv6 []string
	asn []string
	publish_point string
}

//handle结构体
type handle struct{
	publish_point string
	parents map[string]parentHandle
}

// 返回拼接的mysql dsn url，并且返回数据库配置结构体，存在sql注入漏洞--！
func getDSN() (string,*databaseConfig) {
	var config databaseConfig
	//使用viper的反序列化，需要注意的是结构体的field名需要大写首字母
	err := viper.Sub("databaseConfig").Unmarshal(&config) 
	if err != nil {
		slog.Error(err.Error())
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",config.User,config.Password,config.Host,config.Port,config.Database)
	return url,&config
}

// 连接mysql数据库并初始化相关参数，返回数据库配置结构体
func connect() (*sql.DB,*databaseConfig){
	dsn,config := getDSN()
	slog.Info("connecting mysql server","DSN",dsn)
	db,err := sql.Open("mysql",dsn)
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
	return db,config
}

func preProcessing(ipResources ,asResources,uri,aia string){
	
}

func dfs(stmt *sql.Stmt,depth int){
	rirs := viper.GetStringSlice("databaseConfig.rirs")
	for _,v := range rirs{
		rows,err := stmt.Query(rirs_map[v])
		if err != nil {
			slog.Error(err.Error())
		}	
		for rows.Next() {
			var IPResources ,ASResources,URI string 
			rows.Scan(&IPResources,&ASResources,&URI) 

			fmt.Println(URI)
		}
	}
}

func GenerateData(){
	db,config := connect()
	defer db.Close()
	prepare_stmt :=fmt.Sprintf("select IPResources, ASResources,URI from %s.%s where aia = ? and isvalid = 1",config.Database,config.Tables["cas"])
	stmt,err := db.Prepare(prepare_stmt)
	defer stmt.Close()
	if err != nil{
		slog.Error(err.Error())
	}
	dfs(stmt,config.LimitLayer)
}