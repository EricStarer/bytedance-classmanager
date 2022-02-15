package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)
/**
获取参数,不管是什么方式传的参数都可以接收到
使用方法
入参 request的struct要规定好json名
出参 json 需要用json.Unmarshal(jsons,&入参结构对象)
例子
	var test types.Test
	jsons := utils.GetParams(c, test)
	json.Unmarshal(jsons,&test)
	此时test就可以直接用

return nil 表示转换失败
 */
func GetParams(context *gin.Context,requestStruct interface{})  []byte{
	marshal, _:= json.Marshal(requestStruct)
	var needParamsMap map[string]interface{}
	json.Unmarshal(marshal,&needParamsMap)
	var needParams []string
	for key,_ := range needParamsMap {
		needParams=append(needParams,key)
	}

	params:=make(map[string]interface{})
	if context.Request.Method == "POST"{
		for _, val := range needParams {
			if data, ok := context.GetPostForm(val); ok{
				params[val]= data
			}
		}
	}else if context.Request.Method == "GET"{
		for _,val := range needParams {
			if data:=context.Query(val); len(data)>0 {
				if reflect.TypeOf(data).Name() == "string" &&
					(val =="Offset" || val == "Limit" || val =="offset" || val =="limit"){
					num,_:=strconv.Atoi(data)
					params[val]=num
				}else{
					params[val]= data
				}
				continue
			}
			if data, ok := context.GetPostForm(val); ok{
				params[val]= data
			}
		}
	}else{
		fmt.Println("没接收到数据,有错误")
	}
	if len(params) != len(needParams){
		decoder :=json.NewDecoder(context.Request.Body)
		decoder.Decode(&params)
	}

	jsons, err := json.Marshal(&params)
	if err !=nil{
		fmt.Println("参数转换失败1")
		return nil
	}
	return jsons
}

func convertToString(value interface{}) string{
	switch value.(type) {
	case string:
		return value.(string)
	case bool:
		return strconv.FormatBool(value.(bool))
	case int,int64,int32,int16,int8:
		return strconv.FormatInt(value.(int64),10)
	case uint64,uint32,uint16,uint8:
		return strconv.FormatUint(value.(uint64),10)
	case float64,float32:
		return strconv.FormatFloat(value.(float64), 'f', 0, 64)
	case complex64,complex128:
		return strconv.FormatComplex(value.(complex128),'f',0,128)
	default:
		return ""
	}
}
