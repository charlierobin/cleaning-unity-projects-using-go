package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var logFile *os.File

func main() {

	args := os.Args[1:]

	if len(args) != 1 {
		return
	}

	pathOfTempFile := args[0]

	content, err := os.ReadFile(pathOfTempFile)

	if err != nil {
		writeErrorToLog(err)
	}

	str := strings.Trim(string((content)), "\n")

	lines := strings.Split(str, "\n")

	caption := lines[0]

	options := []string{}
	indexes := []int{}

	for i := 1; i < len(lines); i++ {
		options = append(options, lines[i])
	}

	fmt.Println("")

	var entered string

	for {

		fmt.Println("Available options:")
		fmt.Println("")

		for i := 0; i < len(options); i++ {

			if contains(indexes, i) {
				fmt.Println(">", i+1, options[i], "<")
			} else {
				fmt.Println(" ", i+1, options[i], " ")
			}
		}

		fmt.Println("")

		fmt.Print("❓ " + caption + " (Enter a number or pattern or “*” (all), prefix “-” (minus) to deselect, 0 to finish) ")
		fmt.Scanln(&entered)

		fmt.Println("")

		if entered == "" {
			continue
		}

		if entered == "0" {
			break
		}

		remove := false

		if entered[0:1] == "-" {
			remove = true
			entered = entered[1:]
		}

		if entered == "*" {

			if remove {

				for i := 0; i < len(options); i++ {
					if contains(indexes, i) {
						indexes = removeFromArray(indexes, i)
					}
				}
			} else {

				for i := 0; i < len(options); i++ {
					if !contains(indexes, i) {
						indexes = append(indexes, i)
					}
				}
			}
		}

		index, err := strconv.Atoi(entered)

		if err != nil {

			for i := 0; i < len(options); i++ {
				if strings.Contains(options[i], entered) {
					if remove {
						if contains(indexes, i) {
							indexes = removeFromArray(indexes, i)
						}
					} else {
						if !contains(indexes, i) {
							indexes = append(indexes, i)
						}
					}
				}
			}
		} else {

			if index > 0 && index <= len(options) {

				index--

				if remove {
					if contains(indexes, index) {
						indexes = removeFromArray(indexes, index)
					}
				} else {
					if !contains(indexes, index) {
						indexes = append(indexes, index)
					}
				}
			}
		}
	}

	tempFile := ""

	for i := 0; i < len(indexes); i++ {
		tempFile = tempFile + strconv.Itoa(indexes[i]) + "\n"
	}

	bytes := []byte(tempFile)

	err = os.WriteFile(pathOfTempFile+"-response", bytes, 0644)

	if err != nil {
		writeErrorToLog(err)
	}
}

func writeToLog(message string) {

	if logFile == nil {

		directory, _ := os.Executable()
		directory, _ = filepath.Split(directory)

		logFile, err := os.OpenFile(filepath.Join(directory, "logfile.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}

		defer logFile.Close()

		log.SetOutput(logFile)
	}

	log.Println(message)
}

func writeErrorToLog(err error) {
	writeToLog(err.Error())
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func removeFromArray(array []int, e int) []int {

	var new []int

	for _, a := range array {
		if a != e {
			new = append(new, a)
		}
	}
	return new
}
