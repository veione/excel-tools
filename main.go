package main

import (
	"excel-tools/export"
	"excel-tools/types"
	"excel-tools/util"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// Conf of yaml struct
type Conf struct {
	Config struct {
		// 输入文件目录
		Input string
		// 过滤文件
		Excludes string
		Output   struct {
			// 输出格式
			Format string
			// 输出是否格式化
			Pretty bool
			// 是否开启单个文件为对象
			Single bool
			// 客户端输出目录
			Client string
			// 服务端输出目录
			Server string
		}
	}
}

// MergeCell 合并单元格
type MergeCell struct {
	col   int
	row   int
	value string
}

// ReadConf 读取配置文件
func ReadConf() Conf {
	var conf Conf
	f, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	// 将读取的yaml文件解析为struct
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	return conf
}

// ReadFiles from specified path
func ReadFiles(path string) ([]string, error) {
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return nil, err
	}

	var s []string
	for _, fi := range rd {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), "xlsx") || strings.HasSuffix(fi.Name(), "xls") {
			fullName := path + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func main() {
	conf := ReadConf()
	fmt.Printf("input excel directory: %s, excludes: %s\r\n", conf.Config.Input, conf.Config.Excludes)
	// 导出工厂
	exportFactory := export.FileExportFactory{}
	// 类型工厂
	typeFactory := types.TypeFactory{}
	// 根据导出格式获取导出实现对象
	exp := exportFactory.GetExport(conf.Config.Output.Format)

	files, err := ReadFiles(conf.Config.Input)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("read files: %s\n", files)

	// 过滤文件
	var excludes []string
	if len(conf.Config.Excludes) > 0 {
		excludes = strings.Split(conf.Config.Excludes, ",")
	}

	// 统计数据
	var (
		start   = time.Now()
		succeed = 0
		total   = 0
	)

	for _, file := range files {
		if len(excludes) > 0 && util.ArrayContainMember(file, excludes) {
			fmt.Printf("skip file %s\r\n", file)
			continue
		}

		f, err := excelize.OpenFile(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		sheets := f.GetSheetList()
		for _, sheet := range sheets {

			// 如果sheet页以#号开头表示忽略该sheet
			if strings.HasPrefix(sheet, "#") {
				continue
			}
			fmt.Printf("Parse sheet: %s, file: %s\r\n", sheet, file)
			total++
			rows, err := f.GetRows(sheet)
			if err != nil {
				fmt.Println(err)
				return
			}
			// 第一行注释
			notes := rows[0]
			// 第二行字段名
			names := rows[1]
			// 第三行类型
			forms := rows[2]
			// 第四行输出端
			outs := rows[3]
			// 合并单元格
			cells, _ := f.GetMergeCells(sheet)
			// 合并单元格值
			var mergeValues []MergeCell

			if len(cells) > 0 {
				for _, cell := range cells {
					startCol, startRow, _ := excelize.CellNameToCoordinates(cell.GetStartAxis())
					endCol, endRow, _ := excelize.CellNameToCoordinates(cell.GetEndAxis())
					for j := startRow - 1; j <= endRow-1; j++ {
						for i := startCol - 1; i <= endCol-1; i++ {
							mergeValues = append(mergeValues, MergeCell{
								i, j, cell.GetCellValue(),
							})
						}
					}
				}
			}

			// 存储客户端/服务器列表
			var clients []map[string]interface{}
			var servers []map[string]interface{}

			for rowIndex, row := range rows {
				if rowIndex < 4 {
					continue
				}
				// 列数据
				client := make(map[string]interface{})
				server := make(map[string]interface{})

				for colIndex, value := range row {
					note := notes[colIndex]
					// 如果注释标记有#号表示忽略该字段
					if strings.HasPrefix(note, "#") {
						continue
					}

					if value == "" {
						for _, cell := range mergeValues {
							if cell.col == colIndex && cell.row == rowIndex {
								value = cell.value
								break
							}
						}
					}

					name := names[colIndex]
					form := forms[colIndex]
					out := outs[colIndex]
					if value != "" {
						// 类型转换
						value := typeFactory.GetConvert(form).Handle(value)

						// 如果列输出类型为空或者包含cs或者sc表示会全部输出
						if out == "" || strings.Contains(out, "cs") || strings.Contains(out, "sc") {
							client[name] = value
							server[name] = value
						} else if strings.Contains(out, "c") {
							client[name] = value
						} else if strings.Contains(out, "s") {
							server[name] = value
						}
					}
				}

				if len(client) > 0 {
					clients = append(clients, client)
				}
				if len(server) > 0 {
					servers = append(servers, server)
				}
			}

			// 写出到文件
			clientDst := fmt.Sprintf("%s%s%s%s", conf.Config.Output.Client, string(os.PathSeparator), sheet, ".json")
			serverDst := fmt.Sprintf("%s%s%s%s", conf.Config.Output.Server, string(os.PathSeparator), sheet, ".json")
			exp.Export(clientDst, conf.Config.Output.Pretty, conf.Config.Output.Single, clients)
			exp.Export(serverDst, conf.Config.Output.Pretty, conf.Config.Output.Single, servers)
			succeed++
		}
	}
	fmt.Println("export finished, enjoy it! :)")
	fmt.Printf("total: %d, succeed: %d, fail: %d, time consuming: %d(ms)", total, succeed, total-succeed, time.Now().Sub(start).Milliseconds())
}
