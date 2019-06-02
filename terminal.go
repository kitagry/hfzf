package main

import (
	"bytes"
	"log"
	"strings"

	"github.com/gdamore/tcell"
	"gopkg.in/yaml.v2"
)

type Terminal struct {
	screen tcell.Screen
	style  tcell.Style
	row    int
}

func NewTerminal() (Terminal, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return Terminal{}, err
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	return Terminal{
		screen: s,
		style:  tcell.StyleDefault,
		row:    0,
	}, nil
}

func (t Terminal) Init() error {
	return t.screen.Init()
}

func (t Terminal) Clear() {
	t.screen.Clear()
}

func (t Terminal) Sync() {
	t.screen.Sync()
}

func (t Terminal) Show() {
	t.screen.Show()
}

func (t Terminal) Fini() {
	t.screen.Fini()
}

func (t Terminal) PollEvent() tcell.Event {
	return t.screen.PollEvent()
}

func (t *Terminal) Output(text string, data interface{}) {
	t.putln("> "+text, []int{0}, tcell.ColorPurple)
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
				places := PointPlace(el, text)
				t.putln(el, places, tcell.ColorLightGreen)
			}

			if i != len(da)-1 {
				t.putln("---", []int{}, tcell.Color100)
			}
		}
	}
	t.row = 0
}

func (t *Terminal) putln(str string, highlightPlaces []int, color tcell.Color) {
	t.puts(1, t.row, str, highlightPlaces, color)
	t.row++
}

func (t *Terminal) puts(x, y int, str string, highlightPlaces []int, color tcell.Color) {
	stRunes := []rune(str)
	for i, sr := range stRunes {
		if in(i, highlightPlaces) {
			t.screen.SetContent(x, y, sr, []rune(""), t.style.Foreground(color))
		} else {
			t.screen.SetContent(x, y, sr, []rune(""), t.style)
		}

		bsr := []byte(string(sr))
		if len(bsr) > 1 {
			x += 2
		} else {
			x++
		}
	}
}

func in(el int, array []int) bool {
	for _, ar := range array {
		if ar == el {
			return true
		}
	}
	return false
}
