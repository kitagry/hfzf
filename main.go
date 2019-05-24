package main

import (
	"bytes"
	"fmt"
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
      fmt.Println(err)
      return
    }
    defer file.Close()

    decoder = yaml.NewDecoder(file)
  } else {
    decoder = yaml.NewDecoder(os.Stdin)
  }

  err := decoder.Decode(&mapData)
	if err != nil {
		fmt.Println(err)
		return
	}

	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = s.Init(); err != nil {
		fmt.Println(err)
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
				case tcell.KeyCtrlL:
					s.Sync()
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
				s.Sync()
				output(s, data)
				s.Show()
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
		fmt.Println(err)
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

func fuzzyFind(keyword string, data map[interface{}]interface{}) map[interface{}]interface{} {
	threshold := len([]rune(text))

	result := make(map[interface{}]interface{})
	for k, value := range data {
		if smithWaterman(k.(string), keyword) >= threshold {
			result[k] = value
			continue
		}

		switch v := value.(type) {
		case string:
			// fmt.Println(v, "is string type")
			if smithWaterman(v, keyword) >= threshold {
				result[k] = v
			}
		case []interface{}:
			// fmt.Println(v, "is []interface{} type")
			var tmpData []interface{}
			for _, el := range v {
				if smithWaterman(el.(string), keyword) >= threshold {
					tmpData = append(tmpData, el)
				}
			}

			if len(tmpData) != 0 {
				result[k] = tmpData
			}
		case map[interface{}]interface{}:
			// fmt.Println(v, "is map[interface{}]interface{} type")
			tmp := fuzzyFind(keyword, v)
			if len(tmp) != 0 {
				result[k] = tmp
			}
		default:
			fmt.Println(v)
		}
	}
	return result
}

func smithWaterman(s1, s2 string) int {
	s1Rune := []rune(s1)
	s2Rune := []rune(s2)
	gap := 0
	match := 1
	mismatch := 1

	matrix := make([][]int, len(s1Rune)+1)
	for i := 0; i < len(s1Rune)+1; i++ {
		matrix[i] = make([]int, len(s2Rune)+1)
	}

	maxScore := 0
	for i := 1; i < len(s1Rune)+1; i++ {
		for j := 1; j < len(s2Rune)+1; j++ {
			s1Gap := matrix[i-1][j] - gap
			s2Gap := matrix[i][j-1] - gap

			match := matrix[i-1][j-1] + match
			if s1Rune[i-1] != s2Rune[j-1] {
				match = matrix[i-1][j-1] - mismatch
			}

			matrix[i][j] = max(s1Gap, s2Gap, match, 0)
			maxScore = max(maxScore, matrix[i][j])
		}
	}

	return maxScore
}

func max(s ...int) int {
	maxInt := s[0]
	for _, el := range s {
		if el > maxInt {
			maxInt = el
		}
	}
	return maxInt
}
