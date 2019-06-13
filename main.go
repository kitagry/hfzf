package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

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

	go t.SetKeymap(quit)

	<-quit

	t.Fini()

	encoder := yaml.NewEncoder(os.Stdout)
	data = FuzzyFind(string(t.Keyword.Text), mapData)
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
