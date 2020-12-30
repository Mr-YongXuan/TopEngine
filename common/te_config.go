package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)


//程序集通用变量
var Running bool = true


type Config struct {
	RootPath   string   `json:"RootPath"`
	IndexPage  []string `json:"IndexPage"`
	ServerName string   `json:"ServerName"`
	ServerList []struct {
		Host        string   `json:"Host"`
		Addr        string   `json:"Addr"`
		Port        int      `json:"Port"`
		Root        string   `json:"Root"`
		IndexPage   []string `json:"IndexPage"`
		RouterGroup []string `json:"RouterGroup"`
		Charset     string   `json:"Charset"`
		Cert        bool     `json:"Cert"`
		CertCRT     string   `json:"CertCRT"`
		CertKEY     string   `json:"CertKEY"`
	} `json:"ServerList"`
	ErrorPage string `json:"ErrorPage"`
}

func ReadConf() (config *Config) {
	var (
		readBuf []byte
		err     error
	)
	confPath, _ := os.Getwd()
	confPath = path.Join(confPath, "conf/config.json")
	if readBuf, err = ioutil.ReadFile(confPath); err != nil {
		fmt.Println("ERROR => can not decoding config json!")
		return //ERROR EXIT
	}
	config = &Config{}
	json.Unmarshal(readBuf, config)
	return
}
