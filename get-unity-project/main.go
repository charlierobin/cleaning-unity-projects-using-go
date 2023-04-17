package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var logFile *os.File

func main() {

	args := os.Args[1:]

	if len(args) != 1 {
		return
	}

	path := args[0]

	if !isUnityFolder(path) {
		return
	}

	productName, companyName, version := getUnityProjectInfo(path)

	if version != "" {
		fmt.Printf("%s\t%s\t%s\t%s\n", path, version, productName, companyName)
	}
}

func isUnityFolder(path string) bool {

	f, err := os.Open(filepath.Join(path, "Assets"))
	f.Close()

	if err != nil {
		return false
	}

	f, err = os.Open(filepath.Join(path, "Packages"))
	f.Close()

	if err != nil {
		return false
	}

	f, err = os.Open(filepath.Join(path, "ProjectSettings"))
	f.Close()

	if err != nil {
		return false
	}

	f, err = os.Open(filepath.Join(path, "UserSettings"))
	f.Close()
	
	if err != nil {
		return false
	}

	return true
}

func getUnityProjectInfo(path string) (string, string, string) {

	content, err := ioutil.ReadFile(filepath.Join(path, "ProjectSettings/ProjectVersion.txt"))

	if err != nil {
		writeErrorToLog(err)
		return "", "", ""
	}

	version := "⚠️ Unknown"

	lines := strings.Split(string(content), "\n")

	if len(lines) > 0 {

		version = lines[0]

		version = strings.ReplaceAll(version, "m_EditorVersion: ", "")
	}
	
	content, err = ioutil.ReadFile(filepath.Join(path, "ProjectSettings/ProjectSettings.asset"))

	if err != nil {
		writeErrorToLog(err)
		return "", "", ""
	}

	productName := "⚠️ Unknown"
	companyName := "⚠️ Unknown"

	lines = strings.Split(string(content), "\n")

	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "  productName: ") {
			productName = strings.ReplaceAll(lines[i], "  productName: ", "")
			break
		}
	}

	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "  companyName: ") {
			companyName = strings.ReplaceAll(lines[i], "  companyName: ", "")
			break
		}
	}

	return productName, companyName, version
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
