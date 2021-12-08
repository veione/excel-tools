package types

import (
	"excel-tools/convert"
	"excel-tools/util"
	"github.com/tidwall/gjson"
	"strings"
)

type TypeConverter interface {
	Handle(string) interface{}
}

type NumberTypeConvert struct{}

// Handle 数字类型转换
func (*NumberTypeConvert) Handle(value string) interface{} {
	if !util.IsNumber(value) {
		panic("invalid number type for value: " + value)
	}
	if util.IsInt(value) {
		return convert.Str2Int(value)
	} else if util.IsDec(value) {
		return convert.Str2Float32(value)
	}
	return value
}

type IntTypeConvert struct{}

// Handle 数字类型转换
func (*IntTypeConvert) Handle(value string) interface{} {
	return convert.Str2Int(value)
}

type FloatTypeConvert struct{}

// Handle 数字类型转换
func (*FloatTypeConvert) Handle(value string) interface{} {
	return convert.Str2Float32(value)
}

type LongTypeConvert struct{}

// Handle 数字类型转换
func (*LongTypeConvert) Handle(value string) interface{} {
	return convert.Str2Int64(value)
}

type BoolTypeConverter struct{}

// Handle 整型数据转换
func (*BoolTypeConverter) Handle(value string) interface{} {
	return convert.Str2Bool(value)
}

type StringTypeConverter struct{}

// Handle 字符串转换
func (*StringTypeConverter) Handle(value string) interface{} {
	return value
}

type DateTypeConverter struct{}

// Handle 布尔转换
func (*DateTypeConverter) Handle(value string) interface{} {
	return util.FormatTimeString(value)
}

type ObjectTypeConverter struct{}

// Handle 对象转换
func (*ObjectTypeConverter) Handle(value string) interface{} {
	if gjson.Valid(value) {
		parse := gjson.Parse(value)
		if parse.IsObject() {
			return parse.Value()
		}
		panic("Invalid value for json object: " + value)
	} else {
		// 特殊处理：10001:100,10002:200
		arr := strings.Split(value, ",")
		values := make(map[string]interface{})
		for _, str := range arr {
			v := strings.Split(str, ":")
			values[v[0]] = convert.Str2Int(v[1])
		}
		return values
	}
}

type ArrayTypeConverter struct{}

// Handle 数组转换
func (*ArrayTypeConverter) Handle(value string) interface{} {
	if gjson.Valid(value) {
		parse := gjson.Parse(value)
		if parse.IsArray() {
			return parse.Array()
		}
		panic("Invalid value for json array: " + value)
	} else {
		// 如果不是合法的json array格式则采用 123,456,789这种字符串分隔符的方式进行处理
		arr := strings.Split(value, `,`)
		result := make([]interface{}, 0)
		numCvt := new(NumberTypeConvert)
		for _, str := range arr {
			if util.IsNumber(str) {
				result = append(result, numCvt.Handle(str))
			} else {
				result = append(result, str)
			}
		}
		return result
	}
}

type TypeFactory struct {
}

// GetConvert 根据类型获取对应的类型转换器
func (*TypeFactory) GetConvert(types string) (conv TypeConverter) {
	switch types {
	case "number":
		conv = new(NumberTypeConvert)
	case "int":
		conv = new(IntTypeConvert)
	case "float":
		conv = new(FloatTypeConvert)
	case "long":
		conv = new(LongTypeConvert)
	case "bool":
		conv = new(BoolTypeConverter)
	case "date":
		conv = new(DateTypeConverter)
	case "object":
		conv = new(ObjectTypeConverter)
	case "array":
		conv = new(ArrayTypeConverter)
	case "string":
		conv = new(StringTypeConverter)
	default:
		conv = new(StringTypeConverter)
	}
	return
}
