package spider

import (
	"net/url"
	"os"
	"path"
)

import (
	"minispider/internal/util"
)

type Persistencer interface {
	// SaveToFile saves the data to outputDir with filename urlStr.
	SaveToFile(urlStr string, outputDir string, data []byte) error
}

func (s *Spider) SaveToFile(urlStr string, outputDir string, data []byte) error {
	err := util.BuildDir(outputDir)
	if err != nil {
		return err
	}
	//domain := util.RemoveProtocol(urlStr)
	// Escape the special characters of the domain name to make it meet the requirements of the file name
	filePath := url.QueryEscape(urlStr)
	filePath = path.Join(outputDir, filePath)
	err = os.WriteFile(filePath, data, 0644)
	return err
}
