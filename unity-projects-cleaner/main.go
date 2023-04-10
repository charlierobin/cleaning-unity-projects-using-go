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

type projectInfo struct {
	path    string
	version string
	name    string
	company string
	status  string
}

var unityProjects []projectInfo

func main() {

	fmt.Println("Pick the places you want to scan:")
	fmt.Println("")

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

		for i := 0; i < len(unityProjects); i++ {
			fmt.Printf("%d:\t%s\t%s\t%s\t%s\n", i+1, unityProjects[i].path, unityProjects[i].version, unityProjects[i].name, unityProjects[i].company)
		}

		fmt.Println("")
	}

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
}

func selectScanLocations() []string {

	locations := []string{}

	out, err := exec.Command(executableDirectory("list-mounted-volumes")).Output()

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

	var entered int = 1

	for entered > 0 {

		fmt.Println("Available locations:")

		for i := 0; i < len(volumes); i++ {
			fmt.Println(i+1, volumes[i])
		}

		fmt.Println("")
		fmt.Println("Currently selected:")

		if len(locations) == 0 {
			fmt.Println("(none)")
		} else {
			for i := 0; i < len(locations); i++ {
				fmt.Println(locations[i])
			}
		}

		fmt.Println("")

		fmt.Print("❓ Enter a number (0 to finish) ")
		fmt.Scanln(&entered)

		fmt.Println("")

		if entered > 0 && entered <= len(volumes) {
			locations = append(locations, volumes[entered-1])
			volumes = append(volumes[:entered-1], volumes[entered:]...)
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
