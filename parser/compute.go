package parser
//计算
import (
	"bytes"
	"fmt"
	"strconv"
)

//根据operator来判断进行什么运算。
//在进行运算时
//       ①双操作数的指令会将操作数r2强制类型转换
//		 ②根据operator类型进行相应运算
func Calculation (r1 Value,r2 Value,operator uint8) Value {
	result := Value{}

	switch r1.T {
	case VT_STRING:
		result.T = VT_STRING
		result.V = stringOperation(r1.V.(string),r2,operator)
		return result

	case VT_FLOAT:
		result.T = VT_FLOAT
		result.V = floatOperation(r1.V.(float64),r2,operator)
		return result

	case VT_INT:
		result.T = VT_INT
		result.V = intOperation(r1.V.(int),r2,operator)
		return result
	}
	fmt.Println(r1,operator)
	panic("123")
}

func BoolCalculation (r1 Value, r2 Value,operator uint8) Value {
	result := Value{}
	var ok bool

	switch operator {
	case IN_LAND:
		ok = arbitrarilyToBool(r1) && arbitrarilyToBool(r2)
	case IN_LOR:
		ok = arbitrarilyToBool(r1) || arbitrarilyToBool(r2)
	case IN_EQL:
		ok = arbitrarilyToBool(r1) == arbitrarilyToBool(r2)
	case IN_LSS:
		fallthrough
	case IN_GTR:
		ok = compare(&r1,&r2,operator)
	}
	if ok {
		result.T = VT_TRUE
	} else {
		result.T = VT_FALSE
	}
	return result
}

func compare(r1 *Value,r2 *Value,op uint8) bool {
	switch r1.T {
	case VT_INT:
		return intCompare(r1.V.(int),r2.V.(int),op)
	case VT_FLOAT:
		return floatCompare(r1.V.(float64),r2.V.(float64),op)
	case VT_STRING:
		return stringCompare(r1.V.(string),r2.V.(string),op)
	}
	panic("其他类型不能比较")
}

func intCompare(a int,b int,op uint8) bool {
	switch op {
	case IN_LSS:
		return a < b
	case IN_GTR:
		return a > b
	}
	panic("其他类型不能比较")
}

func floatCompare(a float64,b float64,op uint8) bool {
	switch op {
	case IN_LSS:
		return a < b
	case IN_GTR:
		return a > b
	}
	panic("其他类型不能比较")
}

func stringCompare(a string,b string,op uint8) bool {
	switch op {
	case IN_LSS:
		return a < b
	case IN_GTR:
		return a > b
	}
	panic("其他类型不能比较")
}

func intToString (r1 Value) string {
	return fmt.Sprint(r1.V.(int))
}

func intToFloat (r1 Value) float64 {
	return r1.V.(float64)
}

func intToBool (r1 Value) bool {
	if r1.V.(int) == 0 {
		return false
	}
	return true
}

func stringToFloat (r1 Value) (float64,error) {
	return strconv.ParseFloat(r1.V.(string),64)
}

func stringToINT (r1 Value) (int,error) {
	return strconv.Atoi(r1.V.(string))
}

func stringToBool (r1 Value) bool {
	if r1.V.(string) != "" {
		return true
	}
	return false
}

func floatToString (r1 Value) string {
	return fmt.Sprint(r1.V.(float64))
}

func floatToInt (r1 Value) int {
	return int(r1.V.(float64))
}

func floatToBool (r1 Value) bool {
	if r1.V.(float64) == 0 {
		return false
	}
	return true
}

func boolToString (r1 Value) string {
	if r1.T == VT_FALSE {
		return "Flase"
	}
	return "True"
}

func boolToNumber (r1 Value) int {
	if r1.T == VT_FALSE {
		return 0
	}
	return 1
}

//任意类型转为string类型
func arbitrarilyToString (r2 Value) string {
	switch r2.T {
	case VT_STRING:
		return r2.V.(string)

	case VT_FALSE:
		return boolToString(r2)

	case VT_INT:
		return intToString(r2)

	case VT_FLOAT:
		return floatToString(r2)

	case VT_TRUE:
		return boolToString(r2)
	}
	return ""
}

// string 类型运算
// 加法、乘法
func stringOperation (r1str string,r2 Value,op uint8) string {
	switch op {
	case IN_ADD:
		return stringAdd(r1str,arbitrarilyToString(r2))

	case IN_MUL:        //为乘法时，只能是字符串和int类型相乘
		a,ok := r2.V.(int)
		if ok {
			return stringAndIntMul(r1str,a)
		}
		panic("r2必须为int")
	}
	panic("string类型只有乘法和加法")
}

func stringAdd (r1str,r2str string) string {
	return r1str + r2str
}

func stringAndIntMul (r1str string,n int) string {
	buf := new(bytes.Buffer)
	for i := 0 ; i < n ; i ++ {
		buf.WriteString(r1str)
	}
	return buf.String()
}

func arbitrarilyToFloat (r2 Value) float64 {
	switch r2.T {
	case VT_STRING:
		if f,err := strconv.ParseFloat(r2.V.(string),64) ; err == nil {
			return f
		}
		panic("string 不能转为float")

	case VT_FALSE:
		return float64(0)

	case VT_TRUE:
		return float64(1)

	case VT_INT:
		return float64(r2.V.(int))

	case VT_FLOAT:
		return r2.V.(float64)
	}
	panic("不支持其他类型转为float")
}

// float 类型有 +、-、*、/、
func floatOperation (r1f float64,r2f Value,op uint8) float64 {
	switch op {
	case IN_ADD:
		return r1f + arbitrarilyToFloat(r2f)

	case IN_SUB:
		return r1f - arbitrarilyToFloat(r2f)

	case IN_MUL:
		return r1f * arbitrarilyToFloat(r2f)

	case IN_QUO:
		return r1f / arbitrarilyToFloat(r2f)
	}
	panic("不支持其他运算符")
}

func arbitrarilyToInt (r2 Value) int {
	switch r2.T {
	case VT_STRING:
		number ,err := strconv.Atoi(r2.V.(string))
		if err != nil {
			panic("字符串不能转为INT")
		}
		return number

	case VT_FLOAT:
		return int(r2.V.(float64))

	case VT_INT:
		return int(r2.V.(int))

	case VT_TRUE:
		return int(1)

	case VT_FALSE:
		return int(0)

	}
	panic("不支持其他类型转为INT")
}

func arbitrarilyToBool (v Value) bool {
	switch v.T {
	case VT_TRUE:
		return true
	case VT_FALSE:
		return false
	case VT_STRING:
		return stringToBool(v)
	case VT_FLOAT:
		return floatToBool(v)
	case VT_INT:
		return intToBool(v)
	}
	panic("此类型不能转为bool")
}

func intOperation (r1f int,r2f Value,op uint8) int {
	switch op {
	case IN_ADD:
		return r1f + arbitrarilyToInt(r2f)

	case IN_SUB:
		return r1f - arbitrarilyToInt(r2f)

	case IN_MUL:
		return r1f * arbitrarilyToInt(r2f)

	case IN_QUO:
		return r1f / arbitrarilyToInt(r2f)

	case IN_REM:
		return r1f % arbitrarilyToInt(r2f)

	case IN_SHL:
		return r1f << arbitrarilyToInt(r2f)

	case IN_SHR:
		return r1f >> arbitrarilyToInt(r2f)
	}
	panic("不支持其他类型的运算符")

}

func ReverseBool (v Value) Value {
	value := Value{}
	if arbitrarilyToBool(v) {
		value.T = VT_FALSE
	} else  {
		value.T = VT_TRUE
	}
	return value
}

func TakeNegative (v Value) Value {
	switch v.T {
	case VT_FLOAT:
		v.V = - v.V.(float64)
	case VT_INT:
		v.V = - v.V.(int)
	}
	return v
}
