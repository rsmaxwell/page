package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rsmaxwell/page/internal/basic/version"
	"github.com/rsmaxwell/page/internal/myerror"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func error() {
	dir, _ := os.Getwd()
	fmt.Printf("<p>page, version: %s</p>\n", version.Version())
	fmt.Printf("<p>Current Working Directory: %s</p>\n", dir)
}

func main() {

	fmt.Printf("Content-type: text/html\n\n")

	prefix := ""
	value, exists := os.LookupEnv("PREFIX")
	if exists {
		prefix = value
	}

	requestURI, exists := os.LookupEnv("REQUEST_URI")
	if !exists {
		myerror.New("environment variable 'REQUEST_URI' not found").Handle()
		os.Exit(1)
	}

	u, err := url.Parse(requestURI)
	if err != nil {
		myerror.New(err.Error()).Add("could not parse REQUEST_URI: " + requestURI).Handle()
		os.Exit(1)
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
		myerror.New("too many zooms: " + strings.Join(zooms, ",")).Handle()
	}

	files := q["image"]
	if len(files) < 1 {
		myerror.New("no files: " + requestURI).Handle()
		os.Exit(1)
	} else if len(files) > 1 {
		myerror.New("too many files: " + strings.Join(files, ",")).Handle()
	}

	filename := files[0]

	_, err = os.Stat(prefix + filename)
	if err != nil {
		myerror.New("could not stat file: " + prefix + filename).Add("prefix: " + prefix).Add("filename: " + filename).Handle()
		os.Exit(1)
	}

	prefixDirectory := filepath.Dir(prefix + filename)
	directory := strings.ReplaceAll(filepath.Dir(filename), "\\", "/")

	children, err := ioutil.ReadDir(prefixDirectory)
	if err != nil {
		log.Fatal(err)
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
		myerror.New("file not found: " + filename).Handle()
		os.Exit(1)
	}

	previousButton := ""
	if found > 0 {
		prev := filelist[found-1]
		previousButton = " <div class=\"center-left\">" +
			"<a href=\"" + directory + "/" + prev.Name() + "\">" +
			"<img src=\"images/previous.png\" >" +
			"</a>" +
			"</div> \n"
	}

	nextButton := ""
	if found < len(filelist) {
		next := filelist[found+1]
		nextButton = " <div class=\"center-right\">" +
			"<a href=\"" + directory + "/" + next.Name() + "\">" +
			"<img src=\"images/next.png\" >" +
			"</a>" +
			"</div> \n"
	}

	zoomButton := ""
	image := ""
	if zoom == "scale" {
		zoomButton = " <div class=\"top-center\"><img src=\"images/minus.png\"></div> \n"
		image = " <img src=\"" + filename + "\" class=\"center-fit\" > \n"
	} else {
		zoomButton = " <div class=\"top-center\"><img src=\"images/plus.png\"></div> \n"
		image = " <img src=\"" + filename + "\" class=\"center-fit\" > \n"
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
