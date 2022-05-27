package connect

import (
	"bscan/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"strings"
)

//type SubDomain struct {
//	SubDomainName  string `db:"subDomainName"`
//	RootDomainName string `db:"rootDomainName"`
//}

var Sylas *sqlx.DB
var sqlInfo config.Sylas
var DbConnectStatus bool = false

func LoadConfig(options *config.Config) {
	sqlInfo = options.Sylas
	var dataSourceName = "%s:%s@tcp(%s:%s)/%s"
	dataSourceName = fmt.Sprintf(dataSourceName, sqlInfo.User, sqlInfo.Passwd, sqlInfo.Host, sqlInfo.Port, sqlInfo.Db)
	database, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		//fmt.Println("open mysql failed", err)
		return
	}
	DbConnectStatus = true
	Sylas = database
	initDB()
}

func GetRootDomain(table string) []string {
	var rootDomain []string
	err := Sylas.Select(&rootDomain, fmt.Sprintf("select rootDomainName from %s group by rootDomainName", table))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return rootDomain
}

//查询还未添加进扫描的域名
func GetSubDomain(rootDomain string, table string) []string {
	var subDomains []string
	var sql = fmt.Sprintf("select subDomainName from %s where scanned = 0 and rootDomainName = ?", table)
	err := Sylas.Select(&subDomains, sql, rootDomain)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return subDomains
}

//当添加完毕后，在数据库中将scanned更改为1
func UpdateScanned(domain string, table string) {
	var _, err = Sylas.Exec(fmt.Sprintf("update %s.%s set scanned = 1 where subDomainName = ?", sqlInfo.Db, table), domain)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func InsertAliveDomainInfo(url string, status int32, title string, rootDomain string, table string) {
	// 这个地方 记得修改成双表模式，使用switch或者if来进行判断
	if strings.EqualFold(table, "SubDomain") {
		var _, err = Sylas.Exec("insert ignore into subDomainBscanAlive (url,rootDomainName,status,title) values (?,?,?,?)", url, rootDomain, status, title)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	} else if strings.EqualFold(table, "SimilarSubDomain") {
		var _, err = Sylas.Exec("insert ignore into similarDomainBscanAlive (url,rootDomainName,status,title) values (?,?,?,?)", url, rootDomain, status, title)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}

func initDB() {
	var tables []string
	Sylas.Select(&tables, "show tables")
	var subDomainInit bool = false
	var similarDomainInit bool = false
	for _, i := range tables {
		if strings.EqualFold("SubDomainBscanAlive", i) {
			subDomainInit = true
		} else if strings.EqualFold("SimilarDomainBscanAlive", i) {
			similarDomainInit = true
		}
	}
	if !subDomainInit {
		var _, err = Sylas.Exec("create table SubDomainBscanAlive" +
			"(id int not null AUTO_INCREMENT,\n" +
			"url varchar(128) not null,\n" +
			//"subDomainName varchar(128) not null,\n" +
			"rootDomainName varchar(128) not null,\n" +
			"status varchar(128) not null,\n" +
			"title varchar(128) null,\n" +
			"primary key (id)," +
			"UNIQUE KEY `SubDomainBscanAlive_subDomainName_uindex` (`url`));")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		subDomainInit = true
	}
	if !similarDomainInit {
		var _, err = Sylas.Exec("create table SimilarDomainBscanAlive" +
			"(id int not null AUTO_INCREMENT,\n" +
			"url varchar(128) not null,\n" +
			//"subDomainName varchar(128) not null,\n" +
			"rootDomainName varchar(128) not null,\n" +
			"status varchar(128) not null,\n" +
			"title varchar(128) null,\n" +
			"primary key (id)," +
			"UNIQUE KEY `SimilarDomainBscanAlive_subDomainName_uindex` (`url`));")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		similarDomainInit = true
	}
}
