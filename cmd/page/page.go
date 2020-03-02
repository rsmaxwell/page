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
		fmt.Fprintf(os.Stderr, "environment variable 'REQUEST_URI' not foundIndex\n")
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

	cgi := config.DocumentRoot + "/diaries/pages/page"
	_, err = os.Stat(cgi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not stat cgi file: %s\n", cgi)
	}

	previousHTML := ""
	if foundIndex > 0 {
		prev := filelist[foundIndex-1]
		previousURL := cgi + "?diary=" + diary + "&image=" + prev.Name()
		previousHTML = " <div class=\"center-left\">" +
			"<a href=\"" + previousURL + "\">" +
			"<img src=\"/diaries/images/previous.png\" >" +
			"</a>" +
			"</div> \n"
	}

	nextHTML := ""
	if foundIndex < len(filelist) {
		next := filelist[foundIndex+1]
		nextURL := cgi + "?diary=" + diary + "&image=" + next.Name()
		nextHTML = " <div class=\"center-right\">" +
			"<a href=\"" + nextURL + "\">" +
			"<img src=\"/diaries/images/next.png\" >" +
			"</a>" +
			"</div> \n"
	}

	zoomHTML := ""
	image := ""
	if zoom == "scale" {
		zoomHTML = " <div class=\"top-center\"><img src=\"images/minus.png\"></div> \n"
		image = " <img src=\"" + imagefile + "\" class=\"center-fit\" > \n"
	} else {
		zoomHTML = " <div class=\"top-center\"><img src=\"images/plus.png\"></div> \n"
		image = " <img src=\"" + imagefile + "\" > \n"
	}

	// Write out the html
	content := "<!DOCTYPE html> \n" +
		"<html> \n" +
		"<head> \n" +
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"> \n" +
		"<link rel=\"stylesheet\" type=\"text/css\" href=\"/diaries/css/diary.css\"> \n" +
		"</head> \n" +
		"<body> \n" +
		"<div class=\"imgbox\"> \n" +
		image +
		previousHTML +
		zoomHTML +
		nextHTML +
		"</div> \n" +
		"</body> \n" +
		"</html> \n"

	fmt.Print(content)
}
