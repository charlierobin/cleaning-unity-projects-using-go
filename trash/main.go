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

	args := os.Args[1:]

	if len(args) == 1 {

		toDelete := args[0]

		command := "tell application \"Finder\" to delete POSIX file \"" + toDelete + "\""

		out, err := exec.Command("osascript", "-e", command, "-e", "set fileAlias to result as alias", "-e", "set posixPath to POSIX path of fileAlias").Output()

		if err != nil {
			writeToLog("osascript: " + err.Error())
			return
		}

		str := string(out[:])

		str = strings.Trim(str, "\n")

		fmt.Println(str)
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
