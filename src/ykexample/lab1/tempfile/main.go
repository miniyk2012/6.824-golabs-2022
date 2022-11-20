package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	file, _ := ioutil.TempFile("", "task")
	fmt.Println(file.Name())
	fmt.Fprint(file, "hello")
}
