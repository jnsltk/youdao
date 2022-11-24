package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const apiUrl = "https://fanyi.youdao.com/openapi.do?keyfrom=blog125&key=21376174&type=data&doctype=json&version=1.1&q="

func main() {
	var word string
	if len(os.Args) != 2 {
		fmt.Print("Usage: youdao <Word> ")
	} else {
		word = os.Args[1]
	}

	entry := getEntry(word)

	printEntry(entry)
}

func getUrl(word string) string {
	return apiUrl + url.QueryEscape(word)
}

func getEntry(word string) *Entry {
	var entry = new(Entry)

	resp, err := http.Get(getUrl(word))
	cont, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(string(cont)), &entry)

	if err != nil {
		fmt.Printf("Couldn't get that word for you...\nFull debug: ")
		log.Fatal(err)
	}
	return entry
}

func printEntry(entry *Entry) {
	// Print pronunciation
	if entry.Basic.PhoneticUS != "" && entry.Basic.PhoneticUK != "" {
		fmt.Println("美：[" + entry.Basic.PhoneticUS + "] | 英：[" + entry.Basic.PhoneticUK + "]")
	} else {
		fmt.Println("[" + entry.Basic.Phonetic + "]")
	}
	fmt.Println()
	// Print "explains"
	fmt.Println("简明：")
	for i, s := range entry.Basic.Explains {
		var o string
		if len(entry.Basic.Explains) <= 1 {
			o = ""
		} else {
			o = strconv.Itoa(i+1) + ". "
		}
		fmt.Println(o + s)
	}
	fmt.Println()
	// Print "web"
	fmt.Println("短语：")
	for i, web := range entry.Web {
		o := strconv.Itoa(i + 1)
		fmt.Println(o + ". " + web.Key)
		for j, s := range web.Value {
			if j == len(web.Value)-1 {
				fmt.Println(s)
			} else {
				fmt.Print(s + ", ")
			}
		}
	}
}

type Entry struct {
	Query       string   `json:"query"`
	Translation []string `json:"translation"`
	Basic       struct {
		Phonetic   string   `json:"phonetic"`
		PhoneticUS string   `json:"us-phonetic"`
		PhoneticUK string   `json:"uk-phonetic"`
		Explains   []string `json:"explains"`
	}
	Web []Web `json:"web"`
}

type Web struct {
	Value []string `json:"value"`
	Key   string   `json:"key"`
}
