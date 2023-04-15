package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var logFile *os.File

func main() {

	args := os.Args[1:]

	if len(args) != 1 {
		return
	}

	path, _, productName, _ := getUnityProjectInfo(args[0])

	if path == "" {
		writeToLog("Not a Unity project directory: " + args[0])
		return
	}

	moveToTrash(filepath.Join(path, "Logs"), filepath.Join(path, "Logs (from "+productName+")"))
	moveToTrash(filepath.Join(path, "Library"), filepath.Join(path, "Library (from "+productName+")"))
	moveToTrash(filepath.Join(path, "obj"), filepath.Join(path, "obj (from "+productName+")"))
}

func moveToTrash(path string, renamedPath string) {

	f, err := os.Open(path)
	f.Close()

	if err != nil {
		return
	}

	err = os.Rename(path, renamedPath)

	if err != nil {
		writeErrorToLog(err)
		return
	}

	directory, _ := os.Executable()
	directory, _ = filepath.Split(directory)

	_, err = exec.Command(filepath.Join(directory, "trash"), renamedPath).Output()

	if err != nil {
		writeErrorToLog(err)
		return
	}
}

func getUnityProjectInfo(path string) (string, string, string, string) {

	directory, _ := os.Executable()
	directory, _ = filepath.Split(directory)

	out, err := exec.Command(filepath.Join(directory, "get-unity-project"), path).Output()

	if err != nil {
		writeErrorToLog(err)
		return "", "", "", ""
	}

	str := string(out[:])

	str = strings.Trim(str, "\n")

	if str == "" {
		return "", "", "", ""
	}

	bits := strings.Split(str, "\t")

	return bits[0], bits[1], bits[2], bits[3]
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
