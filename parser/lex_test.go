package parser

import (
	"os"
)

//func TestParser_GetToken(t *testing.T) {
//	f := OpenFile()
//	aFile := bufio.NewReader(f)
//
//	p := InitNewParser(aFile,f.Name())
//
//	z := p.GetToken()
//	for i:=0 ; z.T!=EOF ; i ++ {
//		fmt.Println(z)
//		z = p.GetToken()
//	}
//}
//
//func TestInitNewParser(t *testing.T) {
//	f := OpenFile()
//	aFile := bufio.NewReader(f)
//
//	p := InitNewParser(aFile,f.Name())
//	p.ReadyToCompile()
//}



func OpenFile () *os.File {
	f,err := os.Open("E:\\GO_worker\\liber\\abc.txt")
	if err != nil {
		panic(err)
	}
	return f
}
