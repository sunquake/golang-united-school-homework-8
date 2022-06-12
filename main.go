package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var item Item
var items []Item

func Perform(args Arguments, writer io.Writer) error {
	oper, ok := args["operation"]
	if !ok {
		return errors.New("-operation flag has to be specified")
	}
	fName, ok := args["fileName"]
	if !ok {
		return errors.New("-fileName flag has to be specified")
	}
	file, _ := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	b, _ := ioutil.ReadAll(file)
	fmt.Println("+++++++++++", string(b))
	json.Unmarshal(b, &items)
	fmt.Println("--------unmarshal:", len(items))

	switch oper {
	case "add":
		if args["item"] == "" {
			return errors.New("-item flag has to be specified")
		}
		if item.Id == "" {
			return errors.New("id is missing")
		}
		for _, v := range items {
			if v.Id == item.Id {
				return fmt.Errorf("Item with id %s already exists", v.Id)
			}
		}
		items = append(items, item)
		js, _ := json.Marshal(items)
		file.Write(js)
	case "list":
		_, err := writer.Write(b)
		return err
	case "findById":
		if item.Id == "" {
			return errors.New("id is missing")
		}
		for _, v := range items {
			if v.Id == item.Id {
				js, _ := json.Marshal(v)
				writer.Write(js)
				return nil
			}
		}
		writer.Write([]byte(""))
		return fmt.Errorf("Item with id %s not found", item.Id)
	case "remove":
		if item.Id == "" {
			return errors.New("id is missing")
		}
		temp := []Item{}
		for _, v := range items {
			if v.Id != item.Id {
				temp = append(temp, v)
			}
		}
		if len(items) == len(temp) {
			fmt.Fprintf(writer, "Item with id %s not found", item.Id)
			return fmt.Errorf("Item with id %s not found", item.Id)
		}
		js, _ := json.Marshal(temp)
		file.Write(js)
	default:
		return fmt.Errorf("Operation %s not allowed!", oper)
	}
	return nil
}

func parseArgs() Arguments {
	op := flag.String("operation", "", "")
	itm := flag.String("item", "", "a json")
	fN := flag.String("fileName", "", "")
	flag.Parse()
	json.Unmarshal([]byte(*itm), &item)
	return Arguments{"operation": *op, "id": item.Id, "item": *itm, "fileName": *fN}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
