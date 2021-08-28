package parser

type TokenT  int

const (
	// Special tokens
	ILLEGAL TokenT = iota
	EOF    //文件结束标识符仅词法分析时使用
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	STRING // "abc"

	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...
	COLON    // :

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	ID        //标识符

	//keywords
	KEY
	BREAK
	CONST
	CONTINUE
	ELSE
	IF
	FALSE
	TRUE
	FOR
	DO
	WHILE
	NULL
	DEF
	RETURN
	//模块倒入
	IMPORT
)


type Token struct {
	T  		TokenT   //token 类型
	Value	interface{}
	line    int    //属于第几行
}

func NewOneToken () Token {
	return Token{
		T:     ILLEGAL,
		Value: nil,
		line:  0,
	}
}

//关键字定义
//type keyWord struct {
//	name	string
//	t       TokenT
//
//}
//var keyWords = []keyWord{
//	{"break",BREAK},
//	{"const",CONST},
//	{"continue",CONTINUE},
//	{"else",ELSE},
//	{"if",IF},
//	{"false",FALSE},
//	{"true",TRUE},
//	{"for",FOR},
//	{"do",DO},
//	{"while",WHILE},
//	{"null",NULL},
//	{"import",IMPORT},
//}

var keyWords = map[string]TokenT{
	"break":BREAK,
	"const":CONST,
	"continue":CONTINUE,
	"else":ELSE,
	"if":IF,
	"false":FALSE,
	"true":TRUE,
	"for":FOR,
	"do":DO,
	"while":WHILE,
	"null":NULL,
	"import":IMPORT,
	"def":DEF,
	"return":RETURN,
}
