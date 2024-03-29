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
	} else if strings.TrimSpace(value) == "" {
		values := make(map[string]interface{})
		return values
	} else {
		// 特殊处理：10001:100,10002:200
		arr := strings.Split(value, ",")
		values := make(map[string]interface{})
		for _, str := range arr {
			v := strings.Split(str, ":")
			if len(v) < 2 || v[0] == "" || v[1] == "" {
				panic("瞎配数据，也不检查报错了吧:) " + value)
			}
			values[v[0]] = convert.Str2Int(v[1])
		}
		return values
	}
}

type ObjectStringTypeConverter struct{}

// Handle 对象转换
func (*ObjectStringTypeConverter) Handle(value string) interface{} {
	if gjson.Valid(value) {
		parse := gjson.Parse(value)
		if parse.IsObject() {
			return parse.Value()
		}
		panic("Invalid value for json object: " + value)
	} else if strings.TrimSpace(value) == "" {
		values := make(map[string]interface{})
		return values
	} else {
		// 特殊处理：10001:100,10002:200
		arr := strings.Split(value, ",")
		values := make(map[string]interface{})
		for _, str := range arr {
			v := strings.Split(str, ":")
			if len(v) < 2 {
				panic("Invalid value for object: " + value)
			}
			values[v[0]] = v[1]
		}
		return values
	}
}

type ArrayTypeConverter struct{}

// Handle 数组转换
func (*ArrayTypeConverter) Handle(value string) interface{} {
	// 以标准方式：[1001, 1002]
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") && gjson.Valid(value) {
		parse := gjson.Parse(value)
		if parse.IsArray() {
			return parse.Array()
		}
		panic("Invalid value for json array: " + value)
	} else if strings.TrimSpace(value) == "" {
		// 空字符串则直接返回空数组即可
		result := make([]interface{}, 0)
		return result
	} else if !strings.ContainsAny(value, ",") {
		// 单个方式: 10010
		result := make([]interface{}, 0)
		numCvt := new(NumberTypeConvert)
		if util.IsNumber(value) {
			result = append(result, numCvt.Handle(value))
		} else if len(value) > 0 {
			result = append(result, value)
		}
		return result
	} else {
		// 多个方式：10001, 10002, 1003
		// 如果不是合法的json array格式则采用 123,456,789这种字符串分隔符的方式进行处理
		arr := strings.Split(value, `,`)
		result := make([]interface{}, 0)
		numCvt := new(NumberTypeConvert)
		for _, str := range arr {
			if util.IsNumber(str) {
				result = append(result, numCvt.Handle(str))
			} else if len(str) > 0 {
				result = append(result, str)
			}
		}
		return result
	}
}

type PairTypeConverter struct{}

// Handle 键值转换
func (*PairTypeConverter) Handle(value string) interface{} {
	if gjson.Valid(value) {
		parse := gjson.Parse(value)
		if parse.IsObject() {
			return parse.Value()
		}
		panic("Invalid value for json object: " + value)
	} else if strings.TrimSpace(value) == "" {
		values := make(map[string]interface{})
		return values
	} else {
		// 特殊处理：10001:100,10002:200 -> Pair(10010, 100)
		arr := strings.Split(value, ",")
		values := make(map[string]interface{})
		for _, str := range arr {
			v := strings.Split(str, ":")
			if len(v) < 2 {
				panic("Invalid value for object: " + value)
			}
			values["x"] = convert.Str2Int(v[0])
			values["y"] = convert.Str2Int(v[1])
		}
		return values
	}
}

type TripleTypeConverter struct{}

// Handle 三键值转换
func (*TripleTypeConverter) Handle(value string) interface{} {
	if gjson.Valid(value) {
		parse := gjson.Parse(value)
		if parse.IsObject() {
			return parse.Value()
		}
		panic("Invalid value for json object: " + value)
	} else if strings.TrimSpace(value) == "" {
		values := make(map[string]interface{})
		return values
	} else {
		// 特殊处理：10001:100:1,10002:200:1 -> Triple(10010, 100, 1)
		arr := strings.Split(value, ",")
		values := make(map[string]interface{})
		for _, str := range arr {
			v := strings.Split(str, ":")
			if len(v) < 3 {
				panic("Invalid value for object: " + value)
			}
			values["x"] = convert.Str2Int(v[0])
			values["y"] = convert.Str2Int(v[1])
			values["z"] = convert.Str2Int(v[2])
		}
		return values
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
	case "pair":
		conv = new(PairTypeConverter)
	case "triple":
		conv = new(TripleTypeConverter)
	case "int[]":
		conv = new(ArrayTypeConverter)
	case "string[]":
		conv = new(ArrayTypeConverter)
	case "map<string>":
		conv = new(ObjectStringTypeConverter)
	default:
		conv = new(StringTypeConverter)
	}
	return
}
