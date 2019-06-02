package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/gdamore/tcell"

	"gopkg.in/yaml.v2"
)

var (
	text    = ""
	mapData interface{}
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

	t, err := NewTerminal()
	if err != nil {
		log.Println(err)
		return
	}
	defer t.Fini()

	if err = t.Init(); err != nil {
		log.Println(err)
		return
	}

	quit := make(chan struct{})

	t.Show()

	data := FuzzyFind(text, mapData)

	t.Show()
	t.Output(text, data)

	go func() {
		for {
			ev := t.PollEvent()
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
				data = FuzzyFind(text, mapData)
				t.Clear()
				t.Output(text, data)
				t.Sync()
			case *tcell.EventResize:
				t.Sync()
			}
		}
	}()

	<-quit

	t.Fini()

	encoder := yaml.NewEncoder(os.Stdout)
	switch dat := data.(type) {
	case []interface{}:
		for _, d := range dat {
			err = encoder.Encode(d)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
