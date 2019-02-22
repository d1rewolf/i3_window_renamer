package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"go.i3wm.org/i3"
)

func getWindowIDs(name string) []string {
	log.Printf("Searching for window by name: [%s]", name)

	cmd := exec.Command("xdotool", "search", "--onlyvisible", "--name", "--classname", "--class", name)
	output, err := cmd.Output()

	if err != nil {
		log.Printf("Error occurred when searching for ids for: %s", name)
	}

	// trim dangling newline
	trimmedOutput := strings.TrimSuffix(string(output), "\n")

	// split on newline
	ids := strings.Split(trimmedOutput, "\n")

	return ids
}

func setWindowName(name string, ids []string) string {
	var result string
	for _, oneID := range ids {
		binResult, err := exec.Command("xdotool", "set_window", "--name", name, oneID).CombinedOutput()
		result += string(binResult)

		if err != nil {
			log.Printf("Set error occurred. Result: %v Error: %v", result, err)
		}
	}

	return result
}

func removeBadCharactersFromTitle(name string) string {
	// Certain characters cause xdotool to fail. Replace with wildcards
	chars := []string{"$", "?", "&", "(", ")"}
	cleansedName := name

	for _, badChar := range chars {
		cleansedName = strings.Replace(cleansedName, badChar, ".", -1)
	}

	return cleansedName
}

func main() {
	// We listen for title change events and then change the title to what we wish. This
	// change itself will trigger another event, so we save id and title in the following
	// savedIDs and savedTitle to detect these follow-on events and prevent an endless
	// loop of events
	var savedIDs []string
	var savedName string

	recv := i3.Subscribe(i3.WindowEventType)

EVENT_LOOP:
	for recv.Next() {
		ev := recv.Event().(*i3.WindowEvent)

		log.Printf("--------------------- %s event ---------------------", ev.Change)

		if ev.Change == "title" {
			id := ev.Container.ID
			instance := ev.Container.WindowProperties.Instance
			class := ev.Container.WindowProperties.Class
			name := ev.Container.Name

			log.Printf("id: [%v] savedID: [%v] title: [%v] savedTitle: [%v]", id, savedIDs, name, savedName)

			cleansedName := removeBadCharactersFromTitle(name)

			windowIDs := getWindowIDs(cleansedName)

			if savedName == name {
				for _, windowID := range windowIDs {
					for _, savedID := range savedIDs {
						log.Printf("Comparing %s to %s", windowID, savedID)
						if windowID == savedID {
							log.Println("Stopping attempt on same savedId and savedTitle")
							continue EVENT_LOOP
						}
					}
				}
			}

			log.Printf("Id: %v Instance: %v Class %v", id, instance, class)

			newName := fmt.Sprintf("%s %s", instance, name)

			log.Printf("=====================> windowIDs: %s\n", windowIDs)

			result := setWindowName(newName, windowIDs)

			// save id and window name to prevent looping title events
			savedIDs = windowIDs
			savedName = newName

			log.Printf("Result: %s\n", result)
		}
	}
	log.Fatal(recv.Close())
}
