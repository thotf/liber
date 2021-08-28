package cli

import "os"

func OpenFile () *os.File {
	f,err := os.Open("./abc.txt")
	if err != nil {
		panic(err)
	}
	return f
}

