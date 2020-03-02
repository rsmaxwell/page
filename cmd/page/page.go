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
	"github.com/rsmaxwell/page/internal/debug"
)

var (
	pkg          = debug.NewPackage("main")
	functionMain = debug.NewFunction(pkg, "main")
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
	f := functionMain

	fmt.Printf("Content-type: text/html\n\n")

	config := config.New()

	f.Infof(fmt.Sprintf("---[ page: %s ]------------", version.Version()))

	requestURI, exists := os.LookupEnv("REQUEST_URI")
	if !exists {
		f.DebugInfo("help! info")
		f.Fatalf("environment variable 'REQUEST_URI' not found")
	}

	u, err := url.Parse(requestURI)
	if err != nil {
		f.Fatalf("could not parse REQUEST_URI: " + requestURI)
	}

	q := u.Query()

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
		f.Fatalf("too many zooms: " + strings.Join(zooms, ","))
	}

	files := q["image"]
	if len(files) < 1 {
		f.Fatalf("no files: " + requestURI)
		os.Exit(1)
	} else if len(files) > 1 {
		f.Fatalf("too many files: " + strings.Join(files, ","))
	}

	filename := files[0]

	imagefile := filepath.Join(config.Prefix, filename)
	_, err = os.Stat(imagefile)
	if err != nil {
		f.Fatalf("could not stat file: " + imagefile + ", prefix: " + config.Prefix + ", filename: " + filename)
	}

	prefixDirectory := filepath.Dir(imagefile)

	children, err := ioutil.ReadDir(prefixDirectory)
	if err != nil {
		f.Fatalf(err.Error())
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
		f.Fatalf("file not found: " + filename)
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
	}

	zoomButton := ""
	image := ""
	if zoom == "scale" {
		zoomButton = " <div class=\"top-center\"><img src=\"images/minus.png\"></div> \n"
		image = " <img src=\"" + imagefile + "\" class=\"center-fit\" > \n"
	} else {
		zoomButton = " <div class=\"top-center\"><img src=\"images/plus.png\"></div> \n"
		image = " <img src=\"" + imagefile + "\" class=\"center-fit\" > \n"
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
