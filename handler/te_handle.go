package handler

import (
	"TopEngine/common"
	"os"
)

var config = common.ReadConf()

type Config struct {
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
}

type HandleMethods struct {
	Params Config
}

func (hm *HandleMethods) InitMethods(params Config) {
	hm.Params = params
}

func isFile(f string) int64 {
	fi, e := os.Stat(f)
	if e != nil {
		return 0
	}
	if fi.IsDir() {
		return 1
	} else {
		return 2
	}
}

func valInArray(val string, array []string) bool {
	for _, value := range array {
		if val == value {
			return true
		}
	}
	return false
}

func getLastMod(filePath string) int64 {
	f, _ := os.Open(filePath)
	fi, _ := f.Stat()
	return fi.ModTime().Unix()
}

func getFileSize(filePath string) int64 {
	f, _ := os.Open(filePath)
	fi, _ := f.Stat()
	return fi.Size()
}
