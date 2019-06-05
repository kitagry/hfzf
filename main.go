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

	data := FuzzyFind(string(t.Keyword.Text), mapData)
	t.Output(data)

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
					t.Keyword.Delete()
				case tcell.KeyLeft:
					t.Keyword.MoveLeft()
				case tcell.KeyRight:
					t.Keyword.MoveRight()
				case tcell.KeyRune:
					t.Keyword.Input(ev.Rune())
				}
				data = FuzzyFind(string(t.Keyword.Text), mapData)
				t.Clear()
				t.Output(data)
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
