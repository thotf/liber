package parser

type Frame struct {
	curFn 		*closure //当前活动记录的任务
	ip			uint		 //当前命令位置
	preFrame    *Frame
	StackTop    int
	stack       []Value
}

//getAnOpCode

//获取一条指令
func (fe *Frame) GetAnInstruction (ip uint) uint32 {
	if uint(len(fe.curFn.f.instRuction)) == ip {
		return uint32(uint32(IN_EXIT) << 24)
	}
	return fe.curFn.f.instRuction[ip]
}

//stackBotm 当前活动记录栈底
func InitNewFrame(c *closure,ip uint,StackTop int,preFrame *Frame) *Frame {
	return	&Frame{
		curFn:      c,
		ip:         ip,
		StackTop:StackTop,
		preFrame:preFrame,
		stack:nil,
	}
}

