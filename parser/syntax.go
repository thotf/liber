package parser

import "fmt"

//语法分析核心

//语法分析将token进行分析
//
//以空间换时间
//指令类型
//无参数	   8位指令   +  24 位 空位
//有一个参数   8位指令   +  24 位 操作数
//有两个参数   8位指令   +  12位操作数1  + 12位操作数2



//遇到变量 要有个写入到常量表的函
func init() {
	Instruction = map[string]OpCode{
		"NIL":{0,nil,nil,BP_NONE},
		"ID":{0,nil,idNud,BP_NONE},
		"+":{IN_ADD,TwoOperandsLed,nil,BP_TERM},
		"CONST":{IN_LOAD,nil,constNud,BP_NONE},    //所有的字面量都用这条指令
		"-":{IN_SUB,TwoOperandsLed,OneOperandsNegateNud,BP_TERM},
		"*":{IN_MUL,TwoOperandsLed,nil,BP_FACTOR},
		"%":{IN_REM,TwoOperandsLed,nil,BP_FACTOR},
		"/":{IN_QUO,TwoOperandsLed,nil,BP_FACTOR},
		"|":{IN_OR,TwoOperandsLed,nil,BP_BIT_OR},
		"&":{IN_AND,TwoOperandsLed,nil,BP_BIT_AND},
		"^":{IN_XOR,TwoOperandsLed,nil,BP_BIT_OR},
		"<<":{IN_SHL,TwoOperandsLed,nil,BP_BIT_SHIFT},
		">>":{IN_SHR,TwoOperandsLed,nil,BP_BIT_SHIFT},
		"&^":{IN_AND_NOT,TwoOperandsLed,nil,BP_BIT_AND},
		"(":{0,LparenLed,LparenNud,BP_CALL},
		")":{0,nil,nil,BP_NONE},
		"!":{IN_NOT,nil,OneOperandsReverseNud,BP_UNARY},
		"<":{IN_LSS,TwoOperandsLed,nil,BP_CMP},
		">":{IN_GTR,TwoOperandsLed,nil,BP_CMP},
		"||":{IN_LOR,TwoOperandsLed,nil,BP_LOGIC_OR},
		"&&":{IN_LAND,TwoOperandsLed,nil,BP_LOGIC_AND},
		",":{0,nil,nil,BP_NONE},
		"==":{IN_EQL,TwoOperandsLed,nil,BP_EQUAL},
		"[":{0,listLed,listGenerateNud,BP_CALL},
		"]":{0,nil,nil,BP_NONE},
	}
}
var Instruction map[string]OpCode

//指令集
const (
	IN_STORE	= uint8(iota)   //从栈顶移动数据到符号表，弹出栈顶
	IN_MOVE              //从栈中复制数据到栈顶
	IN_OUT_MOVE          //将栈顶数据复制到栈中，弹出
	IN_LOAD              //从符号表移动数据到栈顶
	//IN_OUT_LOAD          //将栈顶数据复制到符号表，弹出
	IN_OVER_LOAD		 //从全局变量表移动数据到栈顶
	IN_PUSH              //新的局部变量入栈
	IN_POP              //栈顶出栈，一般用于弹出局部变量、返回值
	IN_OUT_OVER_LOAD     //将栈顶数据复制到全局变量，弹出
	IN_UPVALUE_LOAD		 //从upvalue移动数据到栈顶  *(*(*(*(*
	IN_IF
	IN_JUMP				//跳转
	IN_WHILE
	IN_JUMP_B

	IN_NEW_LIST			//创建一个list结构放在栈顶
	IN_LIST_APPEND      //像栈顶的list中添加一个数据
	IN_LIST_ADD			//像栈顶的list中 key，value  ,栈顶为value，top-1为key
	IN_LIST_GET         //
	IN_LIST_OUT_OVER_LOAD   //将栈顶的数据添加到全局列表中的list
	IN_LIST_OUT_MOVE        //将栈顶的数据添加到当前局部列表中的list


	IN_ADD               //栈顶两个数相加结果放在栈顶
	IN_SUB               // -
	IN_MUL               //栈顶两个数相乘结果放在栈顶
	IN_QUO               //栈顶两个数相除结果放在栈顶
	IN_REM               //栈顶两个数做模运算
	IN_NEGATE            //栈顶取负数

	IN_AND     // &
	IN_OR      // |
	IN_XOR     // ^
	IN_SHL     // <<
	IN_SHR     // >>
	IN_AND_NOT // &^
	IN_LAND  // &&
	IN_LOR   // ||
	IN_INC   // ++
	IN_DEC   // --

	IN_EQL    // ==
	IN_LSS    // <
	IN_GTR    // >
	IN_ASSIGN // =
	IN_NOT    // !

	IN_CALL
	IN_STACK_CALL  //栈顶调用函数
	IN_RETURN
	IN_EXIT


)
//led函数，渴求左操作数的符号，比如+ - * /
type  LED func (parser *Parser)

//nud方法，不需要左操作数的的符号，包括各种字面量加载到符号表、！、取负等
type  NUD func (parser *Parser)


type level int

const (
	BP_NONE		level = iota  //无绑定能力
	BP_LOWEST	  			 //最低绑定能力 )
	// BP_ASSIGN                //  =
	BP_LOGIC_OR              // ||
	BP_LOGIC_AND			 // &&
	BP_EQUAL				 // ==  !=
	BP_CMP					 // < > <= >=
	BP_BIT_OR				 // | ^ XOR
	BP_BIT_AND				 // &   &^
	BP_BIT_SHIFT			 // << >>
	BP_RANGE                 // ..
	BP_TERM					 // + -
	BP_FACTOR                // * / %
	BP_UNARY				 // - ! ~
	BP_CALL					 // . ( []
	BP_HIGHEST               // 最高
)

type OpCode struct {
	b     uint8
	led   LED
	nud   NUD
	level level
}

//A = abc()
//只能返回一个参数
func LparenLed (p *Parser) {
	param := uint32(0)
	for p.curToken.T != RPAREN {
		param ++
		Expression(p,0)
		if p.curToken.T == COMMA {
			p.GetToken()
		}
	}
	p.GetToken()  //读入)下一个符号

	p.StoreOneInstruction(tableVarDisplace(IN_STACK_CALL, param))
}

func LparenNud (p *Parser) {
	Expression(p,BP_LOWEST)

	if p.curToken.T == RPAREN {
		p.GetToken()
	}else{
		panic("缺少)")
	}
}

func listGenerateNud (p *Parser) {
	//curToken 为常量1  [ 1,
	p.StoreOneInstruction(tableVarDisplace(IN_NEW_LIST,0))

	index := uint32(0)
	for {
		if p.curToken.T == RBRACK {
			break
		}
		Expression(p, 0)
		if p.curToken.T == COMMA { //[1,2,3]
			p.StoreOneInstruction(tableVarDisplace(IN_LIST_APPEND, index))
			p.GetToken()
		} else if p.curToken.T == COLON { //["string":a]
			p.GetToken()
			Expression(p,0)
			p.StoreOneInstruction(tableVarDisplace(IN_LIST_ADD, index))
		} else { //[1]
			p.StoreOneInstruction(tableVarDisplace(IN_LIST_APPEND, index))
			break
		}
		index ++
	}
	Assert(p,RBRACK)
	p.GetToken()
}

// a = abc[1]
func listLed (p *Parser) {
	Expression(p,0)
	Assert(p,RBRACK)
	p.StoreOneInstruction(tableVarDisplace(IN_LIST_GET,0))
	p.GetToken()
}

//单目运算符的，- 取负 nud
func OneOperandsNegateNud (p *Parser) {
	fmt.Println(ReturnOpCode(p))
	nud := Instruction[ReturnOpCode(p)].nud
	p.GetToken()
	nud(p)
	p.StoreOneInstruction(ledInstructions(IN_NEGATE))

}
//单目运算，！ 取反 nud
func OneOperandsReverseNud (p *Parser) {
	nud := Instruction[ReturnOpCode(p)].nud
	p.GetToken()
	nud(p)
	p.StoreOneInstruction(ledInstructions(IN_NOT))

}

//双目运算符的led方法
func TwoOperandsLed (p *Parser)  {
	op := Instruction[TypeToOInstruction(p.preToken)]
	Expression(p,op.level)
	p.StoreOneInstruction(ledInstructions(op.b))
}

//语法分析核心
func Expression (p *Parser,rbp level) {
	// aSwTe
	//刚进来时curToken是操作数w、preToken是运算符S
	nud := Instruction[ReturnOpCode(p)].nud
	//表达式开头只能是操作数、前缀运算符，必然存在nud方法
	if nud == nil {
		fmt.Println(ReturnOpCode(p))
		panic("语句错误，不存在nud方法")
	}
	p.GetToken()  //curToken指向运算符T
	nud(p)
	for rbp < Instruction[ReturnOpCode(p)].level {
		led := Instruction[ReturnOpCode(p)].led
		p.GetToken()   // curToken 指向下一个操作数e
		led(p)
	}

}

func ReturnPreoPCODE (P *Parser) string {
	return TypeToOInstruction(P.preToken)
}

func ReturnOpCode(p *Parser) string {
	return TypeToOInstruction(p.curToken)
}

//输入token的类型，输出对应的指令
func TypeToOInstruction (t Token) string {
	switch t.T {
	case EOF:
		return "NIL"
	case ADD:
		return "+"
	case SUB:
		return "-"
	case ID:   //标识符加载
		return "ID"
	case INT:
		fallthrough
	case STRING:
		fallthrough
	case FLOAT:
		fallthrough
	case TRUE:
		fallthrough
	case FALSE:
		return "CONST"
	case MUL:
		return "*"
	case QUO:
		return "/"
	case REM:
		return "%"
	case AND :
		return "&"
	case OR:
		return "|"
	case XOR:
		return "^"
	case SHL:
		return "<<"
	case SHR:
		return ">>"
	case AND_NOT:
		return "&^"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACK:
		return "["
	case RBRACK:
		return "]"
	case LSS:
		return "<"
	case GTR:
		return ">"
	case NOT:
		return "!"
	case LOR:
		return "||"
	case LAND:
		return "&&"
	case RETURN:
		return "return"
	case COMMA:
		return ","
	case EQL:
		return "=="
	}
	return "NIL"
}

//ID所使用的nud方法，
//① 查询ID是否存在，存在则返回index，不存在报错
func idNud (p *Parser)  {
	VarType,index,ok := findVarInLocalAndUpValueAndGlobal(p)
	if ok != false {
		//写入VarType,index的变量入栈指令
		p.StoreOneInstruction(varTypeLoadInstruct(VarType,uint8(index)))
		return
	}
	//ID不存在
	panic(fmt.Sprint("ID不存在：",p.preToken.Value))
}

//int类型字面量的nud方法，先加载到常量表拿到index，后生成常量加载到栈顶指令
func constNud (p *Parser) {
	op := Instruction[TypeToOInstruction(p.preToken)].b
	consVarIndex := p.CompileUnit.storeConst(TokenChangeValue(&p.preToken))
	p.CompileUnit.StoreOneInstruction(tableVarDisplace(op,consVarIndex))
}

//根据varType生成对应的加载入栈的指令
func varTypeLoadInstruct(t LocalVarType,index uint8) uint32 {
	switch t {
	case GLOBAL_VARIABLE:
		return insDisplace(IN_OVER_LOAD,index)
	case LOCAL_VARIABLE:
		return insDisplace(IN_MOVE,index)
	}

	//upvalue 入栈指令
	return insDisplace(IN_UPVALUE_LOAD,index)
}

//写入将栈顶值传给type ， index 的位置的指令
func assignmentInstruction (p *Parser,varType LocalVarType,index uint8) {
	switch varType {
	case GLOBAL_VARIABLE:
		p.StoreOneInstruction(tableVarDisplace(IN_OUT_OVER_LOAD,uint32(index)))
	case LOCAL_VARIABLE:
		p.StoreOneInstruction(tableVarDisplace(IN_OUT_MOVE,uint32(index)))
	case UPVALUE_VARIABLE:
		//
	}
}


//大部分led运算符的指令生成方法
func ledInstructions (i uint8) uint32 {
	ins := uint32(0)
	ins = ins | uint32(i) << 24
	return ins
}

//符号表加载、存储、跳转、创建list指令
func tableVarDisplace (i uint8,index uint32) uint32 {
	ins := uint32(0)
	ins = ins | uint32(i) << 24
	ins = ins | index
	return ins
}

//变量、upValue、全局变量加载指令、return
func insDisplace(i uint8,index uint8) uint32 {
	ins := uint32(0)
	ins = ins | uint32(i) << 24
	ins = ins  | uint32(index)
	return ins
}

//用于新的局部变量入栈
func push () uint32 {
	ins := uint32(0)
	ins  = ins | uint32(IN_PUSH) << 24
	return ins
}

func pop () uint32 {
	ins := uint32(0)
	ins  = ins | uint32(IN_POP) << 24
	return ins
}

//以局部变量、闭包、全局变量的顺序查找变量ID
//index == -1 没找到则添加
func returnVarindex (p *Parser) (LocalVarType,uint8) {
	// level != 0
	if lType,n,ok := findVarInLocalAndUpValueAndGlobal(p) ; ok {
		return lType,n
	}

	if p.CompileUnit.ScopeLevel == 0 {
		return GLOBAL_VARIABLE,p.ObjModule.InsertModuleVariable(returnIdName(p))
	}

	//此时作用域！= 0，且没找到变量，则在当前作用域创建局部变量
	return LOCAL_VARIABLE,p.CompileUnit.AddLocalVar(returnIdName(p))
}

func findVarInLocalAndUpValueAndGlobal(p *Parser) (LocalVarType,uint8,bool) {
	//全局作用域没有局部变量、UPVALUE
	if p.CompileUnit.ScopeLevel == 0 {
		//全局变量查找
		index,b :=  p.ObjModule.FindVariable(returnIdName(p))
		return GLOBAL_VARIABLE,index,b
	}
	//查找局部变量
	if index,ok := p.CompileUnit.findLocalVar(returnIdName(p)) ; ok != false {

		return LOCAL_VARIABLE,index,ok
	}
	//查找闭包

	//全局变量查找
	if index,ok := p.ObjModule.FindVariable(returnIdName(p)) ; ok != false {
		return GLOBAL_VARIABLE,index,true
	}
	//以上都不存在
	return NILL_VARIABLE,0,false
}

//返回preToken的标识符string形式
func returnIdName (p *Parser) string {
	return p.preToken.Value.(string)
}