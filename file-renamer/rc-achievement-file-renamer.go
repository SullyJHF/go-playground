package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	prependLocked := flag.Bool("locked", false, "Whether to prepend LOCKED_ to the filename")
	baseDir := flag.String("dir", "", "Base dir for achievement images")
	flag.Parse()
	if *baseDir == "" {
		log.Fatal("You must provide a base directory using the -dir flag")
	}

	var achDir string
	if *prependLocked {
		achDir = path.Join(*baseDir, "Unachieved")
	} else {
		achDir = path.Join(*baseDir, "Achieved")
	}
	items, err := os.ReadDir(achDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		fileName := item.Name()
		absPath := path.Join(achDir, fileName)
		// remove files starting with "._"
		if strings.HasPrefix(fileName, "._") {
			fmt.Println(fileName, " has prefix, ._")
			fmt.Println("Removing", absPath)
			os.Remove(absPath)
			continue
		}

		if strings.Contains(fileName, "desktop.ini") {
			os.Remove(absPath)
			continue
		}

		fmt.Println(fileName)
		upperFileName := strings.ToUpper(fileName)
		underscored := strings.Replace(upperFileName, " ", "_", -1)
		noAmp := strings.Replace(underscored, "&", "", -1)
		noExcl := strings.Replace(noAmp, "!", "", -1)
		noComma := strings.Replace(noExcl, ",", "", -1)
		singleUnderscore := strings.Replace(noComma, "__", "_", -1)
		noDash := strings.Replace(singleUnderscore, "-", "_", -1)
		noApo := strings.Replace(noDash, "'", "", -1)
		var finalFilename string
		if *prependLocked && !strings.HasPrefix(noApo, "LOCKED_") {
			finalFilename = "LOCKED_" + noApo
		} else {
			finalFilename = noApo
		}
		fmt.Println(finalFilename)
		os.Rename(absPath, path.Join(achDir, finalFilename))
	}
}
