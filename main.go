package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	oper := args["operation"]
	if oper == "" {
		return errors.New("-operation flag has to be specified")
	}
	fName := args["fileName"]
	if fName == "" {
		return errors.New("-fileName flag has to be specified")
	}
	b, _ := os.ReadFile(fName)
	if oper == "list" {
		writer.Write(b)
		return nil
	}

	items := []Item{}
	json.Unmarshal(b, &items)
	id := args["id"]

	if oper == "findById" {
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		for _, v := range items {
			if v.Id == id {
				js, _ := json.Marshal(v)
				writer.Write(js)
				return nil
			}
		}
		writer.Write([]byte{})
		return nil
	}
	file, _ := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	switch oper {
	case "add":
		if args["item"] == "" {
			return errors.New("-item flag has to be specified")
		}
		item := Item{}
		json.Unmarshal([]byte(args["item"]), &item)
		if item.Id == "" {
			break
		}
		for _, v := range items {
			if v.Id == item.Id {
				fmt.Fprintf(writer, "Item with id %s already exists", v.Id)
				return nil
			}
		}
		items = append(items, item)
		js, _ := json.Marshal(items)
		file.Write(js)

	case "remove":
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		temp := []Item{}
		for _, v := range items {
			if v.Id != id {
				temp = append(temp, v)
			}
		}
		if len(items) == len(temp) {
			fmt.Fprintf(writer, "Item with id %s not found", id)
			return nil
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
	id := flag.String("id", "", "")
	flag.Parse()
	return Arguments{"operation": *op, "id": *id, "item": *itm, "fileName": *fN}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
