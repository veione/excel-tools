package export

import "fmt"

// FileExport 导出接口
type FileExport interface {
	Export(path string, dest string)
}

type JsonExport struct{}

// Export JSON格式导出
func (*JsonExport) Export(path string, dest string) {
	fmt.Println("Json export")
}

type CsvExport struct{}

// Export CSV格式导出
func (*CsvExport) Export(path string, dest string) {
	fmt.Println("csv export")
}

type SqlExport struct{}

// Export Sql格式导出
func (*SqlExport) Export(path string, dest string) {
	fmt.Println("sql export")
}

type FileExportFactory struct {
}

// GetExport 使用 FileExportFactory获得导出对象
func (*FileExportFactory) GetExport(format string) (export FileExport) {
	switch format {
	case "json":
		export = new(JsonExport)
	case "csv":
		export = new(CsvExport)
	case "sql":
		export = new(SqlExport)
	default:
		panic(fmt.Sprintf("no such export for %s", format))
	}
	return
}
