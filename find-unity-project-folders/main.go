package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var logFile *os.File
var directoriesToSkip []string
var verbose bool = false

func main() {

	directory, _ := os.Executable()
	directory, _ = filepath.Split(directory)

	config, err := ioutil.ReadFile(filepath.Join(directory, "config.json"))

	if err != nil {
		writeErrorToLog(err)
	}

	json.Unmarshal(config, &directoriesToSkip)

	args := os.Args[1:]

	if len(args) == 0 {
		return
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "-v" {
			verbose = true
		}
		if i > 0 && args[i] != "-v" {
			directoriesToSkip = append(directoriesToSkip, args[i])
		}
	}

	err = filepath.Walk(args[0], tester)

	if err != nil {
		writeErrorToLog(err)
	}
}

func tester(path string, info os.FileInfo, err error) error {

	if err != nil {
		writeErrorToLog(err)
		return filepath.SkipDir
	}

	if info.IsDir() {

		if info.Name()[0:1] == "." {
			if verbose {
				writeToLog("Skipping folder beginning with . " + path)
			}
			return filepath.SkipDir
		}

		if contains(directoriesToSkip, path) {
			if verbose {
				writeToLog("Skipping folder in config.json " + path)
			}
			return filepath.SkipDir
		}

		path, version, productName, companyName := getUnityProjectInfo(path)

		if path != "" {

			fmt.Printf("%s\t%s\t%s\t%s\n", path, version, productName, companyName)

			return filepath.SkipDir
		}
	}

	return nil
}

func getUnityProjectInfo(path string) (string, string, string, string) {

	directory, _ := os.Executable()
	directory, _ = filepath.Split(directory)

	out, err := exec.Command(filepath.Join(directory, "get-unity-project"), path).Output()

	if err != nil {
		writeErrorToLog(err)
		return "", "", "", ""
	}

	if len(out) == 0 {
		return "", "", "", ""
	}

	str := string(out[:])

	str = strings.Trim(str, "\n")

	bits := strings.Split(str, "\t")

	return bits[0], bits[1], bits[2], bits[3]
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
