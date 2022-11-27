package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/mozillazg/go-pinyin"
)

// Perhaps in the future should change to web scraping as the api is deprecated and might become unavailable in the future
const apiUrl = "https://fanyi.youdao.com/openapi.do?keyfrom=blog125&key=21376174&type=data&doctype=json&version=1.1&q="

// Declare all tones
var tones = [][]string {
	{"ā", "ē", "ī", "ō", "ū", "ǖ"},
	{"á", "é", "í", "ó", "ú", "ǘ"},
	{"ǎ", "ě", "ǐ", "ǒ", "ǔ", "ǚ"},
	{"à", "è", "ì", "ò", "ù", "ǜ"},
	{"a", "e", "i", "o", "u", "ü"},
}

func main() {
	var word string
	if len(os.Args) != 2 {
		fmt.Print("Usage: youdao <Word>")
		os.Exit(1)
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
	cont, err := io.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(string(cont)), &entry)

	if err != nil {
		fmt.Printf("Couldn't get that word for you...\nFull debug: ")
		log.Fatal(err)
	}
	return entry
}

func printEntry(entry *Entry) {
	// Check if there's a result
	if len(entry.Basic.Explains) == 0 {
		fmt.Println("Word not found!")
		return
	}
	// Print pronunciation
	if entry.Basic.PhoneticUS != "" && entry.Basic.PhoneticUK != "" {
		fmt.Println(entry.Query+" |", color.MagentaString("美："), color.YellowString("["+entry.Basic.PhoneticUS+"]"),
			"|", color.MagentaString("英："), color.YellowString("["+entry.Basic.PhoneticUK+"]"))
	} else {
		// Need to format this
		// fmt.Println(entry.Query+" |", "["+entry.Basic.Phonetic+"]", getWordPinyin(entry.Basic.Phonetic))
		fmt.Print(entry.Query + " ")
		printPronColor(entry.Basic.Phonetic)
	}
	fmt.Println()
	// Print "explains" and pinyin
	a := pinyin.NewArgs()
	a.Style = pinyin.Tone
	color.Magenta("简明：")
	for i, s := range entry.Basic.Explains {
		var o string
		if len(entry.Basic.Explains) <= 1 {
			o = ""
		} else {
			o = color.GreenString(strconv.Itoa(i+1) + ". ")
		}
		fmt.Print(o + s + " ")
		// printPinyinSent(pinyin.Pinyin(s, a))
		fmt.Print("\n")
	}

	// Print "web"
	if entry.Web != nil {
		fmt.Println()
		color.Magenta("短语：")
		for i, web := range entry.Web {
			o := strconv.Itoa(i + 1)
			fmt.Println(color.GreenString(o+"."), color.YellowString(web.Key))
			for j, s := range web.Value {
				if j == len(web.Value)-1 {
					fmt.Println(s)
				} else {
					fmt.Print(s + ", ")
				}
			}
		}
	}
	fmt.Println()
}

func printPronColor(pron string) {
	firstTone := color.New(color.FgBlue).PrintfFunc()
	secondTone := color.New(color.FgGreen).PrintfFunc()
	thirdTone := color.New(color.FgYellow).PrintfFunc()
	fourthTone := color.New(color.FgRed).PrintfFunc()
	fifthTone := color.New(color.FgHiBlack).PrintfFunc()

	wordsArr := getPinyinWords(pron)
	if wordsArr == nil {
		fmt.Println("["+pron+"]")
	} else {
		fmt.Print("[")
		for i, s := range wordsArr {
			tone := getWordPinyin(s)
			switch(tone) {
			case 1: firstTone(s)
			case 2: secondTone(s)
			case 3: thirdTone(s)
			case 4: fourthTone(s)
			case 5: fifthTone(s)
			}
			if len(wordsArr) > 1 && i != len(wordsArr) -1 {
				fmt.Print(" ")				
			}
		}
		fmt.Print("]\n")
	}

}

func printPinyinSent(sent [][]string) {
	fmt.Print(sent)
}

func getWordPinyin(pron string) int {
	var tone int

	out:
	for i, t := range tones {
		for _, s := range t {
			if strings.Contains(pron, s) {
				tone = i + 1
				break out
			}
		}
	}
	return tone
}

func isPinyin(pron string) bool {
	for i:=0;i<4;i++ {
		for _, s := range tones[i] {
			if strings.Contains(pron, s) {
				return true
			}
		}
	}
	return false
}

func getPinyinWords(pron string) []string {
	if isPinyin(pron) {
		if strings.Contains(pron, " ") {
			return strings.Split(pron, " ")
		} else {
			return []string{pron}
		}
	} else {
		return nil
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
