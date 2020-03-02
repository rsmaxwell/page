package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rsmaxwell/page/internal/version"

	"github.com/rsmaxwell/page/internal/config"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {

	fmt.Printf("Content-type: text/html\n\n")

	config := config.New()
	fmt.Fprintf(os.Stderr, "config.Prefix:"+config.Prefix+"\n")

	fmt.Fprintf(os.Stderr, "---[ page: %s ]------------\n", version.Version())

	requestURI, exists := os.LookupEnv("REQUEST_URI")
	if !exists {
		fmt.Fprintf(os.Stderr, "environment variable 'REQUEST_URI' not found\n")
	}

	u, err := url.Parse(requestURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse requestURI\n")
		fmt.Fprintf(os.Stderr, "    requestURI: %s\n", requestURI)
	}

	q := u.Query()
	if len(q) < 1 {
		fmt.Fprintf(os.Stderr, "requestURI had empty query\n")
		fmt.Fprintf(os.Stderr, "    requestURI: %s\n", requestURI)
		fmt.Fprintf(os.Stderr, "    query: %+v\n", q)
	}

	zooms := q["zoom"]
	zoom := "scale"
	if len(zooms) < 1 {
		zoom = "scale"
	} else if len(zooms) == 1 {
		value := zooms[0]
		validZooms := []string{"scale", "orig"}
		if contains(validZooms, strings.ToLower(value)) {
			zoom = value
		}
	} else {
		fmt.Fprintf(os.Stderr, "too many 'zooms' keys in requestURI\n")
		fmt.Fprintf(os.Stderr, "    requestURI: %s\n", requestURI)
		for i, z := range zooms {
			fmt.Fprintf(os.Stderr, "    zoom[%d]: %s\n", i, z)
		}
	}

	files := q["image"]
	if len(files) < 1 {
		fmt.Fprintf(os.Stderr, "Missing 'image' key in requestURI query\n")
		fmt.Fprintf(os.Stderr, "    requestURI: %s\n", requestURI)
		fmt.Fprintf(os.Stderr, "    q: %+v\n", q)
		os.Exit(1)
	} else if len(files) > 1 {
		fmt.Fprintf(os.Stderr, "too many 'image' keys in requestURI\n")
		fmt.Fprintf(os.Stderr, "    requestURI: %s\n", requestURI)
		for i, f := range files {
			fmt.Fprintf(os.Stderr, "    file[%d]: %s\n", i, f)
		}
	}

	filename := files[0]

	imagefile := filepath.Join(config.Prefix, filename)
	_, err = os.Stat(imagefile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not stat image file\n")
		fmt.Fprintf(os.Stderr, "    image file: %s\n", imagefile)
		fmt.Fprintf(os.Stderr, "    config.Prefix: %s\n", config.Prefix)
		fmt.Fprintf(os.Stderr, "    filename: %s\n", filename)
	}

	prefixDirectory := filepath.Dir(imagefile)

	children, err := ioutil.ReadDir(prefixDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	// list the files with the same parent, sorted by name
	validExtensions := []string{".jpg", ".jpeg", ".png"}
	var filelist = make([]os.FileInfo, 0)
	for _, child := range children {
		extension := filepath.Ext(child.Name())
		if contains(validExtensions, strings.ToLower(extension)) {
			filelist = append(filelist, child)
		}
	}

	sort.Slice(filelist, func(i, j int) bool {
		return filelist[i].Name() < filelist[j].Name()
	})

	found := -1
	for i, f := range filelist {
		if filepath.Base(filename) == f.Name() {
			found = i
		}
	}

	if found < 0 {
		fmt.Fprintf(os.Stderr, "file not found: %s\n", filename)
	}

	previousButton := ""
	if found > 0 {
		prev := filelist[found-1]
		previousFile := filepath.Join(config.Prefix, prev.Name())
		previousButton = " <div class=\"center-left\">" +
			"<a href=\"" + previousFile + "\">" +
			"<img src=\"images/previous.png\" >" +
			"</a>" +
			"</div> \n"

		_, err = os.Stat(previousFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not stat previous file\n")
			fmt.Fprintf(os.Stderr, "    previousFile: %s\n", previousFile)
			fmt.Fprintf(os.Stderr, "    config.Prefix: %s\n", config.Prefix)
			fmt.Fprintf(os.Stderr, "    prev.Name(): %s\n", prev.Name())
		}
	}

	nextButton := ""
	if found < len(filelist) {
		next := filelist[found+1]
		nextFile := filepath.Join(config.Prefix, next.Name())
		nextButton = " <div class=\"center-right\">" +
			"<a href=\"" + nextFile + "\">" +
			"<img src=\"images/next.png\" >" +
			"</a>" +
			"</div> \n"

		_, err = os.Stat(nextFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not stat previous file\n")
			fmt.Fprintf(os.Stderr, "    nextFile: %s\n", nextFile)
			fmt.Fprintf(os.Stderr, "    config.Prefix: %s\n", config.Prefix)
			fmt.Fprintf(os.Stderr, "    next.Name(): %s\n", next.Name())
		}
	}

	zoomButton := ""
	image := ""
	if zoom == "scale" {
		zoomButton = " <div class=\"top-center\"><img src=\"images/minus.png\"></div> \n"
		image = " <img src=\"" + imagefile + "\" class=\"center-fit\" > \n"
	} else {
		zoomButton = " <div class=\"top-center\"><img src=\"images/plus.png\"></div> \n"
		image = " <img src=\"" + imagefile + "\" > \n"
	}

	// Write out the html
	content := "<!DOCTYPE html> \n" +
		"<html> \n" +
		"<head> \n" +
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"> \n" +
		"<link rel=\"stylesheet\" type=\"text/css\" href=\"../css/diary.css\"> \n" +
		"</head> \n" +
		"<body> \n" +
		"<div class=\"imgbox\"> \n" +
		image +
		previousButton +
		zoomButton +
		nextButton +
		"</div> \n" +
		"</body> \n" +
		"</html> \n"

	fmt.Print(content)
}
