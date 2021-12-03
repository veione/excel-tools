package types

import "strconv"

type TypeConverter interface {
	Handle(string) interface{}
}

type IntTypeConverter struct{}

// Handle 整型数据转换
func (*IntTypeConverter) Handle(value string) interface{} {
	val, _ := strconv.Atoi(value)
	return val
}

type StringTypeConverter struct{}

// Handle 字符串转换
func (*StringTypeConverter) Handle(value string) interface{} {
	return value
}

type BoolTypeConverter struct{}

// Handle 布尔转换
func (*BoolTypeConverter) Handle(value string) interface{} {
	val, _ := strconv.ParseBool(value)
	return val
}

type FloatTypeConverter struct{}

// Handle 布尔转换
func (*FloatTypeConverter) Handle(value string) interface{} {
	val, _ := strconv.ParseFloat(value, 2)
	return val
}
