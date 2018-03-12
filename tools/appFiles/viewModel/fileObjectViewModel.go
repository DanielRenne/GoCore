package viewModel

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type FileObjectViewModel struct {
	imageResize []byte
	FileObject  model.FileObject `json:"FileObject"`
	Width       int              `json:"Width"`
	Height      int              `json:"Height"`
}

func (this *FileObjectViewModel) LoadDefaultState() {

}

func (self *FileObjectViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}

func (self *FileObjectViewModel) Write(content []byte) (n int, err error) {
	n = len(content)
	self.imageResize = append(self.imageResize, content...)
	return
}

func (self *FileObjectViewModel) SaveResize() {

	self.FileObject.Size = len(self.imageResize)
	value := base64.StdEncoding.EncodeToString(self.imageResize)
	self.FileObject.Content = value
	self.FileObject.MD5 = fmt.Sprintf("%x", md5.Sum([]byte(self.FileObject.Content)))
}

type FileObjectResizeViewModel struct {
	ImageResize []byte
}

func (self *FileObjectResizeViewModel) Write(content []byte) (n int, err error) {
	n = len(content)
	self.ImageResize = append(self.ImageResize, content...)
	return
}
