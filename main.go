package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	var exp expense

	set := flag.NewFlagSet("", flag.ExitOnError)
	set.StringVar(&exp.description, "desc", "", "description of expense")
	set.IntVar(&exp.amount, "amount", 0, "amount of expense")
	set.IntVar(&exp.id, "id", 0, "id of expense")

	err := set.Parse(os.Args[2:])
	if err != nil {
		panic(err)
	}

	//fmt.Println(exp.description, exp.amount, os.Args[1])
	file, err := os.OpenFile("data.csv", os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	reader := csv.NewReader(file)

	expenses, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	comand := os.Args[1]
	switch comand {
	case "add":
		exp.date = time.Now()
		if len(expenses) > 0 {
			exp.id, err = strconv.Atoi(expenses[len(expenses)-1][0])
			if err != nil {
				panic(err)
			}
		}
		exp.id++
		_, err = file.WriteString(fmt.Sprintf("%d,%s,%d,%s\n", exp.id, exp.description, exp.amount, exp.date.String()))
		if err != nil {
			panic(err)
		}
	case "delete":
		for i, expense := range expenses {
			id, err := strconv.Atoi(expense[0])
			if err != nil {
				panic(err)
			}
			if exp.id == id {
				expenses = expenses[:i+copy(expenses[i:], expenses[i+1:])]
			}
		}
		err = file.Truncate(0)
		if err != nil {
			panic(err)
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			panic(err)
		}
		writer := csv.NewWriter(file)
		err = writer.WriteAll(expenses)
		if err != nil {
			panic(err)
		}
	case "list":
		for _, expense := range expenses {
			for _, item := range expense {
				fmt.Print(item, "\t")
			}
			fmt.Println()
		}
	case "update":
		for i, expense := range expenses {
			id, err := strconv.Atoi(expense[0])
			if err != nil {
				panic(err)
			}
			if exp.id == id {
				if exp.description != "" {
					expenses[i][1] = exp.description
				}
				if exp.amount != 0 {
					expenses[i][2] = strconv.Itoa(exp.amount)
				}
			}
		}
		err = file.Truncate(0)
		if err != nil {
			panic(err)
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			panic(err)
		}
		writer := csv.NewWriter(file)
		err = writer.WriteAll(expenses)
		if err != nil {
			panic(err)
		}
	}
}

type expense struct {
	id          int
	description string
	amount      int
	date        time.Time
}
