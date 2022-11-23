package main

import (
	"fmt"
	"os"
)

const apiUrl = "https://fanyi.youdao.com/openapi.do?keyfrom=blog125&key=21376174&type=data&doctype=json&version=1.1&q="

func main() {
	var word string
	if len(os.Args) != 2 {
		fmt.Print("Usage: youdao <Word> ")
	} else {
		word = os.Args[1]
	}

	fmt.Println(word)
}

func getEntry() *Entry {
	var entry = new(Entry)
	return entry
}

func printEntry() {

}

type Entry struct {
}
