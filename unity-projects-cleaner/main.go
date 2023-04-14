package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var logFile *os.File

type projectInfo struct {
	path    string
	version string
	name    string
	company string
	status  string
}

var unityProjects []projectInfo

func main() {

	volumes := selectScanLocations()

	if len(volumes) == 0 {
		fmt.Println("Nothing was selected")
		return
	}

	for i := 0; i < len(volumes); i++ {

		fmt.Println("Scanning", volumes[i], "…")
		fmt.Println("")

		unityProjects = append(unityProjects, scanVolume(volumes[i])...)
	}

	if len(unityProjects) == 0 {
		fmt.Println("No Unity projects were found")
		return
	}

	noun := "projects"

	if len(unityProjects) == 1 {
		noun = "project"
	}

	fmt.Println("Found", len(unityProjects), "Unity", noun)
	fmt.Println("")

	var entered string

	fmt.Println("Would you like to review a list of them? (y/n)")
	fmt.Scanln(&entered)
	fmt.Println("")

	if entered == "y" {
		unityProjects = selectProjects()
	}

	if len(unityProjects) == 0 {
		fmt.Println("No Unity projects selected")
		return
	}

	fmt.Println("Selected:")
	fmt.Println("")

	for i := 0; i < len(unityProjects); i++ {
		fmt.Printf("%s\t%s\t%s\t%s\n", unityProjects[i].path, unityProjects[i].version, unityProjects[i].name, unityProjects[i].company)
	}

	fmt.Println("")

	fmt.Println("Are you sure you want to clean them? (y/n)")
	fmt.Scanln(&entered)
	fmt.Println("")

	if entered != "y" {
		return
	}

	fmt.Println("Cleaning…")
	fmt.Println("")

	for i := 0; i < len(unityProjects); i++ {

		result := clean(unityProjects[i].path)

		if result {
			unityProjects[i].status = "Cleaned"
		}
	}

	fmt.Println("Done")
}

func selectProjects() []projectInfo {

	var picked []projectInfo

	tempFileData := "Please pick the projects you want to clean" + "\n"

	for i := 0; i < len(unityProjects); i++ {
		tempFileData = tempFileData + unityProjects[i].path + "\t" + unityProjects[i].version + "\t" + unityProjects[i].name + "\t" + unityProjects[i].company + "\n"
	}

	tempFileDataBytes := []byte(tempFileData)

	tempFile, err := os.CreateTemp(os.TempDir(),"")

	err = os.WriteFile(tempFile.Name(), tempFileDataBytes, 0644)

	if err != nil {
		writeErrorToLog(err)
		return picked
	}

	cmd := exec.Command("osascript",
		"-e", `tell application "Terminal" to activate`,
		"-e", `tell application "System Events" to keystroke "t" using {command down}`,
		"-e", `tell application "Terminal" to do script "`+executableDirectory("picker")+` `+tempFile.Name()+`" in front window`)

	err = cmd.Run()

	if err != nil {
		writeErrorToLog(err)
		return picked
	}

	origFileName := tempFile.Name()

	filenameResponse := origFileName + "-response"

	found := false

	for !found {

		time.Sleep(1 * time.Second)

		f, err := os.Open(filenameResponse)
		f.Close()

		if err == nil {
			found = true
		}
	}

	content, err := ioutil.ReadFile(filenameResponse)

	if err != nil {
		writeErrorToLog(err)
		return picked
	}

	err = os.Remove(origFileName)

	if err != nil {
		writeErrorToLog(err)
	}

	err = os.Remove(filenameResponse)

	if err != nil {
		writeErrorToLog(err)
	}

	str := strings.Trim(string((content)), "\n")

	indexes := strings.Split(str, "\n")

	for i := 0; i < len(indexes); i++ {

		if indexes[i] != "" {

			index, err := strconv.Atoi(indexes[i])

			if err != nil {
				writeErrorToLog(err)
			}

			picked = append(picked, unityProjects[index])
		}
	}

	return picked
}

func selectScanLocations() []string {

	var locations []string

	out, err := exec.Command(executableDirectory("list-mounted-volumes")).Output()
	// out, err := exec.Command("./list-mounted-volumes").Output()

	if err != nil {
		writeErrorToLog(err)
		return locations
	}

	str := string(out[:])

	str = strings.Trim(str, "\n")

	volumes := strings.Split(str, "\n")

	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		writeErrorToLog(err)
		return locations
	}

	volumes = append(volumes, userHomeDir+"/Documents")

	tempFileData := "Please pick the places you want to scan" + "\n"

	for i := 0; i < len(volumes); i++ {
		tempFileData = tempFileData + volumes[i] + "\n"
	}

	tempFileDataBytes := []byte(tempFileData)

	tempFile, err := os.CreateTemp(os.TempDir(),"")

	err = os.WriteFile(tempFile.Name(), tempFileDataBytes, 0644)

	if err != nil {
		writeErrorToLog(err)
		return locations
	}

	cmd := exec.Command("osascript",
		"-e", `tell application "Terminal" to activate`,
		"-e", `tell application "System Events" to keystroke "t" using {command down}`,
		"-e", `tell application "Terminal" to do script "`+executableDirectory("picker")+` `+tempFile.Name()+`" in front window`)

	err = cmd.Run()

	if err != nil {
		writeErrorToLog(err)
		return locations
	}

	origFileName := tempFile.Name()

	filenameResponse := origFileName + "-response"

	found := false

	for !found {

		time.Sleep(1 * time.Second)

		f, err := os.Open(filenameResponse)
		f.Close()

		if err != nil {
			// fmt.Println("not found")
		} else {
			// fmt.Println("found")
			found = true
		}
	}

	content, err := ioutil.ReadFile(filenameResponse)

	if err != nil {
		writeErrorToLog(err)
		return locations
	}

	err = os.Remove(origFileName)

	if err != nil {
		writeErrorToLog(err)
	}

	err = os.Remove(filenameResponse)

	if err != nil {
		writeErrorToLog(err)
	}

	str = strings.Trim(string((content)), "\n")

	indexes := strings.Split(str, "\n")

	for i := 0; i < len(indexes); i++ {

		if indexes[i] != "" {

			index, err := strconv.Atoi(indexes[i])

			if err != nil {
				writeErrorToLog(err)
			}

			locations = append(locations, volumes[index])
		}
	}

	return locations
}

func scanVolume(path string) []projectInfo {

	var foundProjects []projectInfo

	out, err := exec.Command(executableDirectory("find-unity-project-folders"), path).Output()

	if err != nil {
		writeErrorToLog(err)
		return foundProjects
	}

	str := string(out[:])

	str = strings.Trim(str, "\n")

	lines := strings.Split(str, "\n")

	for i := 0; i < len(lines); i++ {
		bits := strings.Split(lines[i], "\t")
		if len(bits) == 4 {
			foundProjects = append(foundProjects, projectInfo{bits[0], bits[1], bits[2], bits[3], ""})
		}
	}

	return foundProjects
}

func clean(path string) bool {

	out, err := exec.Command(executableDirectory("clean-unity-project"), path).Output()

	if err != nil {
		writeErrorToLog(err)
		return false
	}

	str := string(out[:])

	str = strings.Trim(str, "\n")

	return true
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

func executableDirectory(name string) string {

	directory, _ := os.Executable()
	directory, _ = filepath.Split(directory)

	return filepath.Join(directory, name)
}
