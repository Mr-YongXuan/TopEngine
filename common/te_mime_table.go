package common

import (
	"strings"
)

type MimeTable struct {
	mime map[string]string
}

func (mt *MimeTable) InitMime() {
	mt.mime = make(map[string]string)
	mt.mime["html"] = "text/html"
	mt.mime["xml"] = "text/xml"
	mt.mime["xhtml"] = "application/xhtml+xml"
	mt.mime["txt"] = "text/plain"
	mt.mime["rtf"] = "application/rtf"
	mt.mime["pdf"] = "application/pdf"
	mt.mime["png"] = "image/png"
	mt.mime["gif"] = "image/gif"
	mt.mime["jpg"] = "image/jpeg"
	mt.mime["jpeg"] = "image/jpeg"
	mt.mime["au"] = "audio/basic"
	mt.mime["binary"] = "application/octet-stream"
	mt.mime["css"] = "text/css"
	mt.mime["js"] = "text/javascript"
	mt.mime["json"] = "application/json"
	mt.mime["ico"] = "image/vnd.microsoft.icon"
	mt.mime["doc"] = "application/msword"
	mt.mime["docx"] = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	mt.mime["ttf"] = "font/ttf"
}

func (mt *MimeTable) GetMime(name string) (mime string) {
	var (
		fileType []string
		ok       bool
	)
	fileType = strings.Split(name, ".")
	if mime, ok = mt.mime[fileType[len(fileType)-1]]; !ok {
		return "application/octet-stream"
	}
	return
}

func (mt *MimeTable) AddMime(fileType string, mime string) {
	mt.mime[fileType] = mime
}
