package main

import (
	"encoding/json"
	"excel-tools/export"
	tps "excel-tools/types"
	"excel-tools/util"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

// Conf of yaml struct
type Conf struct {
	Config struct {
		Input    string
		Format   string
		Excludes string
	}
	Output struct {
		Client string
		Server string
	}
}

// ReadConf 读取配置文件
func ReadConf() Conf {
	var conf Conf
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	// 将读取的yaml文件解析为struct
	err = yaml.Unmarshal(yamlFile, &conf)
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
	fmt.Printf("Input excel directory: %s, format: %s, excludes: %s\r\n", conf.Config.Input, conf.Config.Format, conf.Config.Excludes)
	// 导出工厂
	exportFactory := export.FileExportFactory{}
	// 类型工厂
	typeFactory := tps.TypeFactory{}
	// 根据导出格式获取导出实现对象
	exp := exportFactory.GetExport(conf.Config.Format)

	fmt.Printf("Use %x format to export file\n", exp)

	files, err := ReadFiles(conf.Config.Input)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Read files: %s\n", files)
	for _, file := range files {
		f, err := excelize.OpenFile(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		sheets := f.GetSheetList()
		for index := range sheets {
			sheet := sheets[index]

			// 如果sheet页以#号开头表示忽略该sheet
			if strings.HasPrefix(sheet, "#") {
				continue
			}

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
			types := rows[2]
			// 第四行输出端
			outs := rows[3]
			fmt.Println(notes, names, types, outs)

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

					name := names[colIndex]
					colType := types[colIndex]
					out := outs[colIndex]
					if value != "" {
						// 类型转换
						value := typeFactory.GetConvert(colType).Handle(value)

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
				fmt.Println()
			}

			if data, err := json.Marshal(clients); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("客户端结果：" + string(data))
			}
			fmt.Println()
			if data, err := json.Marshal(servers); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("服务器结果：" + string(data))
			}
		}
	}
}

func test() {
	timeString := util.FormatTimeString("2021-12-12 00:00:00")
	fmt.Println("日期转换：" + timeString)

	jsonValue := `{"itemId": 1001, "num": 100}`
	js := gjson.Parse(jsonValue)
	fmt.Println(js.Map())
	fmt.Println(gjson.Valid(jsonValue))

	fmt.Println(gjson.Parse(`[1001, 1002, 1003]`).Array())

	//1001:100,1002:300
	val := "10001:100,10002:200"
	arr := strings.Split(val, ",")
	values := make(map[string]interface{})
	for _, str := range arr {
		split := strings.Split(str, ":")
		values[split[0]] = split[1]
	}

	fmt.Println(values)
}
