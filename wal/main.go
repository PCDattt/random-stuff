package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type Calculator struct {
	value float32
}

const logFile = "log.txt"

func main() {
	rootCmd := &cobra.Command{
		Use: "calculator",
		Short: "Calculator CLI",
	}
	rootCmd.AddCommand(readCmd())
	rootCmd.AddCommand(clearCmd())
	
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
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

func readCmd() *cobra.Command {
	return &cobra.Command{
		Use: "read",
		Short: "Read log and calculate",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.OpenFile(logFile, os.O_RDWR, 0664)
			if err != nil {
				return err
			}
			fileScanner := bufio.NewScanner(f)
			c := &Calculator{}
			err = readLog(fileScanner, c)
			if err != nil {
				return err
			}
			fmt.Printf("Value after read log: %f\n", c.value)
			stdScanner := bufio.NewScanner(os.Stdin)
			err = scanAndCalculate(stdScanner, f, c)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func clearCmd() *cobra.Command {
	return &cobra.Command{
		Use: "clear",
		Short: "Clear log and calculate",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.OpenFile(logFile, os.O_RDWR, 0664)
			if err != nil {
				return err
			}
			err = clearLog(f)
			if err != nil {
				return err
			}
			fmt.Println("Log cleared")
			c := &Calculator{}
			stdScanner := bufio.NewScanner(os.Stdin)
			err = scanAndCalculate(stdScanner, f, c)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
