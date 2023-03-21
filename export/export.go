package export

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileExport 导出接口
type FileExport interface {
	Export(path string, pretty bool, allowSingle bool, data []map[string]interface{})
}

type JsonExport struct{}

// Export JSON格式导出
func (*JsonExport) Export(dst string, pretty bool, allowSingle bool, values []map[string]interface{}) {
	if len(values) == 0 {
		return
	}
	path, _ := filepath.Split(dst)
	if _, err := os.Stat(path); err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println(dst, err.Error())
			return
		}
	}
	var payload interface{}
	if allowSingle && len(values) == 1 {
		payload = values[0]
	} else {
		payload = values
	}

	var data []byte
	// 是否格式化输出
	if pretty {
		data, _ = json.MarshalIndent(payload, "", "\t")
	} else {
		data, _ = json.Marshal(payload)
	}

	err := ioutil.WriteFile(dst, data, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

type FileExportFactory struct {
}

// GetExport 使用 FileExportFactory获得导出对象
func (*FileExportFactory) GetExport(format string) (export FileExport) {
	switch format {
	case "json":
		export = new(JsonExport)
	default:
		panic(fmt.Sprintf("no such export for %s", format))
	}
	return
}
