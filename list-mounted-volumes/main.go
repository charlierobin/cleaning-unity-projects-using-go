package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var logFile *os.File

func main() {

	ignore := []string{"/", "/dev", "/System/Volumes/Data", "/private/var/vm", "/System/Volumes/Data/home", "/Volumes/Recovery"}

	out, err := exec.Command("df").Output()

	if err != nil {
		writeErrorToLog(err)
	}

	str := string(out[:])

	lines := strings.Split(str, "\n")

	for _, element := range lines {

		bits := strings.Split(element, "%   ")

		if len(bits) > 1 {

			entry := bits[len(bits)-1]

			if !contains(ignore, entry) {
				fmt.Println(entry)
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
