package main

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell"

	"gopkg.in/yaml.v2"
)

var (
	text    = ""
	row     = 0
	mapData map[interface{}]interface{}

	style = tcell.StyleDefault
)

func main() {
	var decoder *yaml.Decoder
	if len(os.Args) >= 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		decoder = yaml.NewDecoder(file)
	} else {
		decoder = yaml.NewDecoder(os.Stdin)
	}

	err := decoder.Decode(&mapData)
	if err != nil {
		log.Println(err)
		return
	}

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

func output(s tcell.Screen, data map[interface{}]interface{}) {
	putln(s, "> "+text)
	buffer := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		log.Println(err)
		return
	}

	str := buffer.String()
	strs := strings.Split(str, "\n")
	for _, el := range strs {
		putln(s, el)
	}
	row = 0
}

func putln(s tcell.Screen, str string) {
	puts(s, style, 1, row, str)
	row++
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	st := []rune(str)
	if len(st) > 0 {
		s.SetContent(x, y, st[0], st[1:], style)
	} else {
		s.SetContent(x, y, []rune(" ")[0], []rune(""), style)
	}
}
