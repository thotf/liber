package parser
//指令对象，用于解析指令

type order uint32

func (o order) OpCode () uint8 {
	return uint8(o >> 24)
}

//获取单操作数
func (o order) R1 () uint32 {
	return uint32(o) & 0x00ffffff
}

//func (o order) R2 () (){}