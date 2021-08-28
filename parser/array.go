package parser

import "fmt"

const (
	DEFAULT_ARRAY_SIZE = 10
	EXPANSION_FACTOR   = 1.25
)

type array struct {
	buf 	[]Value
	len     int
	positon	int
}

func NewAarray () array {
	return array{
		buf: make([]Value,DEFAULT_ARRAY_SIZE),
		len: DEFAULT_ARRAY_SIZE,
	}
}

func (arr *array) append (v Value) {
	if arr.positon < arr.len {
		arr.buf[arr.positon] = v
		arr.positon ++
		return
	}

	//扩容
	nlen := int(float64(arr.len) * EXPANSION_FACTOR)
	nbuf := make([]Value,0,nlen)
	copy(nbuf,arr.buf)
	arr.buf = nbuf
	arr.len = nlen

	arr.buf[arr.positon] = v
	arr.positon ++
	return
}

func (arr *array) set (index int,v Value) {
	if index >=0 && index < arr.len {
		arr.buf[index] = v
		return
	}
	fmt.Printf("array越界,array[%d]\n",index)
	panic("")
}

func (arr *array) find (index int) Value {
	if index >=0 && index < arr.len {
		return arr.buf[index]
	}
	fmt.Printf("array越界,array[%d]\n",index)
	panic("")
}