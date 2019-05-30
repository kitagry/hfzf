package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell"

	"gopkg.in/yaml.v2"
)

var (
	text    = ""
	row     = 0
	mapData interface{}

	style = tcell.StyleDefault
)

func main() {
	var byteData []byte
	if len(os.Args) >= 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		byteData, err = ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		var err error
		byteData, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Println(err)
			return
		}
	}

	byteDatas := bytes.Split(byteData, []byte("---\n"))
	interfaceDatas := make([]interface{}, 0)
	for _, d := range byteDatas {
		var tmp interface{}
		err := yaml.Unmarshal(d, &tmp)
		if err != nil {
			log.Println(err)
			return
		}
		interfaceDatas = append(interfaceDatas, tmp)
	}
	mapData = interfaceDatas

	s, err := tcell.NewScreen()
	if err != nil {
		log.Println(err)
		return
	}

	if err = s.Init(); err != nil {
		log.Println(err)
		return
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	quit := make(chan struct{})

	s.Show()

	data := fuzzyFind(text, mapData)

	s.Show()
	output(s, data)

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyDelete, tcell.KeyBackspace, tcell.KeyBackspace2:
					textRune := []rune(text)
					if len(textRune) != 0 {
						textRune = textRune[0 : len(textRune)-1]
						text = string(textRune)
					}
				case tcell.KeyRune:
					text += string(ev.Rune())
				}
				data := fuzzyFind(text, mapData)
				s.Clear()
				output(s, data)
				s.Sync()
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	<-quit

	s.Fini()
}

func output(s tcell.Screen, data interface{}) {
	putln(s, "> "+text)
	switch da := data.(type) {
	case []interface{}:
		for i, das := range da {
			buffer := new(bytes.Buffer)
			encoder := yaml.NewEncoder(buffer)
			err := encoder.Encode(das)
			if err != nil {
				log.Println(err)
				return
			}

			str := buffer.String()
			strs := strings.Split(str, "\n")
			strs = strs[:len(strs)-1]
			for _, el := range strs {
				putln(s, el)
			}

			if i != len(da)-1 {
				putln(s, "---")
			}
		}
	}
	row = 0
}

func putln(s tcell.Screen, str string) {
	puts(s, style, 1, row, str)
	row++
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	places := pointPlace(str, text)
	stRunes := []rune(str)
	for i, sr := range stRunes {
		if in(i, places) {
			s.SetContent(x, y, sr, []rune(""), style.Foreground(tcell.Color100))
		} else {
			s.SetContent(x, y, sr, []rune(""), style)
		}

		bsr := []byte(string(sr))
		if len(bsr) > 1 {
			x += 2
		} else {
			x++
		}
	}
}
