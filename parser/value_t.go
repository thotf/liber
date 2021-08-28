package parser

//栈中的value，value 类型
type Value_Type int

const (
	VT_NIL  = iota
	VT_STRING
	VT_FLOAT
	VT_INT
	VT_TRUE
	VT_FALSE
	VT_FUN  //闭包
	VT_lIST
)

type Value struct {
	T	uint
	V 	interface{}
}


func TokenChangeValue(token *Token) Value {
	var nValue Value
	switch token.T {
	case STRING:
		nValue.T = VT_STRING
	case FLOAT:
		nValue.T = VT_FLOAT
	case TRUE:
		nValue.T = VT_TRUE
	case FALSE:
		nValue.T = VT_FALSE
	case INT:
		nValue.T = VT_INT
	}
	nValue.V = token.Value
	return nValue
}

type commaVar struct {
	varType    LocalVarType
	index      uint8
}