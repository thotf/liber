package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

//每调用一次GetToken，进行一次词法分析，输出词法单元或词素

type Parser struct {
	curToken  Token
	preToken  Token
	curFile   *bufio.Reader
	c         byte
	ObjModule *Module    //当前正在编译的模块
	CompileUnit *CompileUnit  //当前编译的单元
	line      int		//正在编译第几行
}

func (p *Parser) GetToken () Token {
	if p.toKenInitialization() == io.EOF {
		p.curToken.T = EOF
		p.curToken.Value = nil
		p.curToken.line = p.line
		return p.curToken
	}
	//编译其他字符、字符串
	p.parserCharacter()
	if p.curToken.T != ILLEGAL {
		return p.curToken
	}
	//编译数字
	if unicode.IsDigit(rune(p.c)) {
		p.parserNumToken(rune(p.c))
		return p.curToken
	}

	//编译标识符、变量

	p.parserIdentifier()
	return p.curToken
}

func (p *Parser) getChar () error {
	c,err := p.curFile.ReadByte()
	if err == io.EOF {
		return err
	}
	p.c = c
	return nil
}

func (p *Parser) getRune () (rune,error) {
	r,_,err := p.curFile.ReadRune()
	if err == io.EOF {
		return r,err
	}
	return r,nil
}

//r :文件缓冲
//fname : 文件名
func InitNewParser (r *bufio.Reader,fname string) *Parser {
	return &Parser{
		curToken: Token{},
		curFile:  r,
		ObjModule:InitModule(fname),
		CompileUnit:nil,
		line:1,
	}
}

//token初始化
func (p *Parser)toKenInitialization() error {
	p.preToken = p.curToken
	p.curToken = NewOneToken()
	err := p.getChar()
	if err := p.sKipSpaces(); err != nil {
		return err
	}//跳过所有空格
	p.curToken.line = p.line
	return err
}

//查看下一个字符
func (p *Parser) LookNextChar (s string) bool {
	//bytes,err := p.curFile.Peek(len([]byte(s)))
	bytes,err := p.curFile.Peek(1)
	if err == io.EOF {
		fmt.Println("文件结束")
	}else if err != nil {
		panic("出错")
	}
	//if err != nil {
	//	panic(123)
	//}
	if string(bytes) == s {
		return true
	}
	return false
}

//编译字符串
func  (p *Parser) parserString() string {
	buf := new(bytes.Buffer)
	for {
		if rStr, _ := p.getRune(); string(rStr) != "\"" {
			buf.WriteRune(rStr)
			continue
		}
		return buf.String()
	}
}

//编译数字
func  (p *Parser) parserNumToken(c rune) {
	buf := new(bytes.Buffer)

	IsFloat := false
	//是数字或者浮点字面量
	for  {
		buf.WriteRune(c)
		if unicode.IsDigit(rune(p.viewNextOneChar())) || p.viewNextOneChar() == ('.') {
			c, _ = p.getRune()
			if c == rune('.') {
				//一个数字字面量只能有一个'.'
				if IsFloat == true {
					panic("【ERROR】float 只能有一个.")
				}
				IsFloat = true
			}
		}else {
			break
		}
	}

	if IsFloat {
		p.curToken.T = FLOAT
		p.curToken.Value,_ = strconv.ParseFloat(buf.String(),64)
		return
	}

	p.curToken.T = INT
	p.curToken.Value,_ = strconv.Atoi(buf.String())

	return
}

func (p *Parser) viewNextOneChar() byte {
	b,err := p.curFile.Peek(1)
	if err != nil {
		return ' '
	}
	return b[0]
}

//编译字符逻辑
func (p *Parser) parserCharacter() {
	switch p.c {
	case '+' :
		if p.LookNextChar("+") {
			p.curToken.T = INC
			p.getChar()
		}else {
			p.curToken.T = ADD
		}
	case ':' :
		p.curToken.T = COLON
	case '-' :
		if p.LookNextChar("-") {
			p.curToken.T = DEC
			p.getChar()
		}else {
			p.curToken.T = SUB
		}
	case '*' :
		p.curToken.T = MUL
	case '/' :
		p.curToken.T = QUO
	case '%' :
		if p.LookNextChar("^") {
			p.curToken.T = AND_NOT
			p.getChar()
		}else {
			p.curToken.T = REM
		}
	case '&' :
		if p.LookNextChar("&") {
			p.curToken.T = LAND
			p.getChar()
		}else {
			p.curToken.T = AND
		}
	case '|' :
		if p.LookNextChar("|") {
			p.curToken.T = LOR
			p.getChar()
		}else {
			p.curToken.T = OR
		}
	case '^' :
		p.curToken.T = XOR
	case '<' :
		if p.LookNextChar("<") {
			p.curToken.T = SHL
			p.getChar()
		}else if p.LookNextChar("=") {
			p.getChar()
			p.curToken.T = LEQ
		} else {
			p.curToken.T = LSS
		}
	case '>' :
		if p.LookNextChar(">") {
			p.curToken.T = SHR
			p.getChar()
		}else if p.LookNextChar("=") {
			p.curToken.T = GEQ
			p.getChar()
		} else {
			p.curToken.T = GTR
		}
	case '=' :
		if p.LookNextChar("=") {
			p.curToken.T = EQL
			p.getChar()
		}else {
			p.curToken.T = ASSIGN
		}
	case '!' :
		if p.LookNextChar("=") {
			p.curToken.T = NEQ
			p.getChar()
		}else {
			p.curToken.T = NOT
		}
	case '(' :
		p.curToken.T = LPAREN
	case '[' :
		p.curToken.T = LBRACK
	case '{' :
		p.curToken.T = LBRACE
	case ',' :
		p.curToken.T = COMMA
	case '.' :
		p.curToken.T = PERIOD
	case ')' :
		p.curToken.T = RPAREN
	case ']' :
		p.curToken.T = RBRACK
	case '}' :
		p.curToken.T = RBRACE
	case '"' :
		str := p.parserString()  //此时已经读完了最后的"
		p.curToken.T = STRING
		p.curToken.Value = str
	}
}

//编译标识符
//所有标识符都强制改为小写
func (p *Parser) parserIdentifier() {
	p.readID()
	p.isKeyWord()
}

func (p *Parser) readID() {
	buf := new(bytes.Buffer)

	buf.WriteByte(p.c)
	c := p.viewNextOneChar()
	for {
		if unicode.IsDigit(rune(c)) || unicode.IsLetter(rune(c)) || c == '_' {
			p.getChar()
			buf.WriteByte(p.c)
			c = p.viewNextOneChar()
			continue
		}
		break
	}

	p.curToken.T = ID
	p.curToken.Value = string(bytes.ToLower(buf.Bytes()))
}

func (p *Parser) isKeyWord() {
	if _,ok := keyWords[p.curToken.Value.(string)];ok {

		p.curToken.T = keyWords[p.curToken.Value.(string)]
		return
	}
}

//跳过所有空格
func (p *Parser) sKipSpaces () error {
	for unicode.IsSpace(rune(p.c)) {
		if p.c == '\n' {
			p.line ++
		}
		err := p.getChar()
		if err != nil {
			return  err
		}
	}
	return nil
}

//初始化main CompileUnit
func (p *Parser) initMain () {
	p.CompileUnit = InitCompileUnit(0,nil,p)
	p.ObjModule.p = p
}

//准备编译
//找import
//循环加载每个import
func (p *Parser) ReadyToCompile () {
	p.initMain()
	funcLoading(p)
	p.runMainCompile()
}

//编译main函数
func (p *Parser) runMainCompile () {
	//这个getToken必须拿到外面来，如果进入下一行时，在expression中已经拿到了下一行第一个token
	p.GetToken()

	for {
		p.analysisOneSentence()
		if p.curToken.T == EOF {
			//跳出for循环
			goto Loop
		}
	}
	Loop:
}

//解析一条语句
func (p *Parser) analysisOneSentence () {
	switch p.curToken.T {
	case IMPORT:
		//其他模块加载
	case ID:
		//数据处理语句
		p.dataProcessing()
	case FOR:

	case IF:
		p.ifStatement()

	case DEF:
		p.defStatement()

	case RETURN:
		returnParser(p)

	case WHILE:
		p.whileStatement()
	}

}

func (p *Parser) defStatement () {
	defStatement(p)
}

func (p *Parser) whileStatement () {
	whileStatement(p)
}

//数据处理 : 包括赋值、函数调用
func (p *Parser) dataProcessing () {

	if p.LookNextChar("(") {
		numberOfParameters := funcCall(p)
		//在一行直接调用一个函数时，函数所有返回值都没有接受者，舍弃
		for i:=int8(0) ; i < numberOfParameters ; i ++ {  //有几个返回值就弹出几个
			p.StoreOneInstruction(pop())
		}
	}

	p.GetToken()      //读取下一个token 可能是 = 、!=、+= .....
	//preToken为标识符名
	p.assignmentStatement()
}

func (p *Parser) ifStatement () {
	ifStatement(p)
}

//赋值语句处理
func (p *Parser) assignmentStatement () {

	switch p.curToken.T {
	case ASSIGN:			//赋值语句
		assignmentStatementHandle(p)

	case COMMA:
		unpacKingHandle(p)

	case LBRACK:            //切片赋值
		sectionAssignmentStatementHandle(p)
	}
}

//查找变量
//① 如果此时作用域为全局作用域0，那一定是模块变量，直接在模块中找
//② 如果此时在局部作用域中先在局部变量中找，没有找upValue
//func (p *Parser) findVariable () int {
//	vName := p.curToken.Value.(string)
//
//	if p.CompileUnit.ScopeLevel == 0 {
//		//全局作用域
//		return p.findVariableInModule(vName)
//	}
//
//	//这里return 0 是测试用的写完，upValue会删除
//	return 0
//}

func (p *Parser) findVariableInModule(name string) (uint8,bool) {
	return p.ObjModule.FindVariable(name)
}

func (p *Parser) StoreOneInstruction (i uint32) {
	p.CompileUnit.StoreOneInstruction(i)
}

