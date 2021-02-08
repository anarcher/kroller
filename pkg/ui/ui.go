package ui

import (
	"fmt"
	"strings"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

func PrintTitle(title string, verbose bool) {
	if verbose {
		color.Yellow(title)
	}
}

func Print(title string, verbose bool) {
	if verbose {
		fmt.Println(title)
	}
}

func PrintBanner(title string) {
	bannerString := banner.PrintS(title)
	color.Yellow(bannerString)
	fmt.Println("")
}

func AskForConfirm() (bool, error) {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}

	response = strings.ToLower(strings.TrimSpace(response))

	ok := []string{"y", "yes"}
	no := []string{"n", "no"}

	if inString(ok, response) {
		return true, nil
	} else if inString(no, response) {
		return false, nil
	}

	fmt.Println("Please type yes or no and then press enter:")
	return AskForConfirm()
}

func inString(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
