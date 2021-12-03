package main

import (
	"excel-tools/export"
	"fmt"
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
		if !fi.IsDir() {
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
	// 根据导出格式获取导出实现对象
	export := exportFactory.GetExport(conf.Config.Format)

	fmt.Printf("Use %x format to export file\n", export)

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
			// 第四行位置
			sides := rows[3]
			fmt.Println(notes, names, types, sides)

			for rowIndex, row := range rows {
				if rowIndex < 4 {
					continue
				}
				for colIndex, colCell := range row {
					note := notes[colIndex]
					// 如果注释标记有#号表示忽略该字段
					if strings.HasPrefix(note, "#") {
						continue
					}

					name := names[colIndex]
					colType := types[colIndex]
					side := sides[colIndex]
					fmt.Print(colCell, note, name, colType, side, "\t")
				}
				fmt.Println()
			}
		}
	}
}
