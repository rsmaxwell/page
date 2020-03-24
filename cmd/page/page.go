package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rsmaxwell/page/internal/config"
	"github.com/rsmaxwell/page/internal/version"
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
	fmt.Fprintf(os.Stderr, "config.DocumentRoot:"+config.DocumentRoot+"\n")

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

	diaries := q["diary"]
	if len(diaries) < 1 {
		fmt.Fprintf(os.Stderr, "Missing 'diary' key in requestURI query\n")
		fmt.Fprintf(os.Stderr, "    requestURI: %s\n", requestURI)
		fmt.Fprintf(os.Stderr, "    q: %+v\n", q)
		os.Exit(1)
	} else if len(diaries) > 1 {
		fmt.Fprintf(os.Stderr, "too many 'diary' keys in requestURI\n")
		fmt.Fprintf(os.Stderr, "    requestURI: %s\n", requestURI)
		for i, d := range diaries {
			fmt.Fprintf(os.Stderr, "    diary[%d]: %s\n", i, d)
		}
	}
	diary := diaries[0]

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

	imagefile := filepath.Join(config.DocumentRoot, "diaries/pages", diary, filename)
	_, err = os.Stat(imagefile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not stat image file\n")
		fmt.Fprintf(os.Stderr, "    image file: %s\n", imagefile)
		fmt.Fprintf(os.Stderr, "    config.DocumentRoot: %s\n", config.DocumentRoot)
		fmt.Fprintf(os.Stderr, "    'diaries/pages': %s\n", "diaries/pages")
		fmt.Fprintf(os.Stderr, "    diary: %s\n", diary)
		fmt.Fprintf(os.Stderr, "    filename: %s\n", filename)
	}

	imageDirectory := filepath.Dir(imagefile)

	children, err := ioutil.ReadDir(imageDirectory)
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

	foundIndex := -1
	for i, f := range filelist {
		if filepath.Base(filename) == f.Name() {
			foundIndex = i
		}
	}

	if foundIndex < 0 {
		fmt.Fprintf(os.Stderr, "file not found: %s\n", filename)
	}

	previousHTML := ""
	if foundIndex > 0 {
		prev := filelist[foundIndex-1]
		previousURL := config.CgiProgram + "?diary=" + diary + "&image=" + prev.Name()
		previousHTML = "    <div class=\"center-left\"> \n" +
			"        <a href=\"" + previousURL + " \"> \n" +
			"            <img src=\"" + config.DiariesRoot + "/controls/previous.png\"> \n" +
			"        </a>\n" +
			"    </div>\n\n"
	}

	nextHTML := ""
	if foundIndex < len(filelist) {
		next := filelist[foundIndex+1]
		nextURL := config.CgiProgram + "?diary=" + diary + "&image=" + next.Name()
		nextHTML = "    <div class=\"center-right\">\n" +
			"        <a href=\"" + nextURL + " \"> \n" +
			"            <img src=\"" + config.DiariesRoot + "/controls/next.png\"> \n" +
			"        </a> \n" +
			"    </div> \n\n"
	}

	imageRef := config.DiariesRoot + "/pages/" + diary + "/" + filename

	var imageHTML string

	thisURL := config.CgiProgram + "?diary=" + diary + "&image=" + filename
	imageHTML = "  <img class=\"zoom\" src=\"" + imageRef + "\" > \n\n"

	fmt.Fprintf(os.Stderr, "thisURL: %s\n", thisURL)
	fmt.Fprintf(os.Stderr, "imageHTML: %s\n", imageHTML)

	// Write out the html
	content := "<!DOCTYPE html> \n" +
		"<html> \n" +
		"<head> \n" +
		"<meta charset=\"utf-8\"> \n" +
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"> \n" +
		"\n" +
		"    <link rel=\"stylesheet\" type=\"text/css\" href=\"css/page.css\"> \n" +
		"\n" +
		"    <title>Pages</title> \n" +
		"</head> \n" +
		"<body> \n" +
		"<div class=\"imgbox\"> \n" +
		imageHTML +
		previousHTML +
		nextHTML +
		"</div> \n" +
		"\n" +
		"    <script src=\"scripts/wheelzoom.js\"></script> \n" +
		"    <script> \n" +
		"	    wheelzoom(document.querySelector('img.zoom')); \n" +
		"    </script> \n" +
		"</body> \n" +
		"</html> \n"

	fmt.Print(content)
}
