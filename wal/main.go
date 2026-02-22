package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Calculator struct {
	value float32
}

const logFile = "log.txt"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	c := &Calculator{}
	f, err := os.OpenFile(logFile, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Println("1: Calculate")
	fmt.Println("2: Read log and calculate")
	fmt.Println("3: Clear log and calculate")
	scanner.Scan()
	line := scanner.Text()
	switch line {
	case "1":
		err := scanAndCalculate(scanner, f, c)
		if err != nil {
			panic(err)
		}
	case "2":
		fileScanner := bufio.NewScanner(f)
		err := readLog(fileScanner, c)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Value after read log: %f\n", c.value)
		err = scanAndCalculate(scanner, f, c)
		if err != nil {
			panic(err)
		}
	case "3":
		err := clearLog(f)
		if err != nil {
			panic(err)
		}
		err = scanAndCalculate(scanner, f, c)
		if err != nil {
			panic(err)
		}
	}
}

func (c *Calculator) Calculate(str string) error {
	operator := str[:1]
	num, err := strconv.Atoi(str[1:])
	if err != nil {
		return err
	}
	n := float32(num)
	if operator == "-" {
		c.value -= n
	}
	if operator == "+" {
		c.value += n
	}
	if operator == "*" {
		c.value *= n
	}
	if operator == "/" {
		c.value /= n
	}
	return nil
}

func writeLog(w io.Writer, content string) error {
	_, err := w.Write([]byte(content + "\n"))
	return err
}

func readLog(s *bufio.Scanner, c *Calculator) error {
	for s.Scan() {
		str := s.Text()
		err := c.Calculate(str)
		if err != nil {
			return err
		}
	}
	return nil
}

func clearLog(f *os.File) error {
	err := f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	return nil
}

func scanAndCalculate(s *bufio.Scanner, f *os.File, c *Calculator) error {
	for s.Scan() {
		line := s.Text()
		err := writeLog(f, line)
		if err != nil {
			return err
		}
		err = c.Calculate(line)
		fmt.Printf("Current value: %f\n", c.value)
		if err != nil {
			return err
		}
	}
	return nil
}
