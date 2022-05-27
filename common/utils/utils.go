package utils

import (
	"fmt"
	urlparser "net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FolderExists checks if a folder exists
func FolderExists(folderpath string) bool {
	_, err := os.Stat(folderpath)
	return !os.IsNotExist(err)
}

// HasStdin determines if the user has piped input
func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	isPipedFromChrDev := (mode & os.ModeCharDevice) == 0
	isPipedFromFIFO := (mode & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}

func GetWindowWidth() int {
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0
	}
	return w
}

func ListAllFileByName(ext string, dir string) []string {
	filelist := make([]string, 0, 0)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "."+ext) == true {
			filelist = append(filelist, path)
		}
		return nil
	})
	return filelist
}

func UrlFormat(url string, path string) string {
	url = strings.TrimRight(url, "/")
	path = strings.TrimLeft(path, "/")
	return url + "/" + path
}

func UrlQueryFormat(url string, query string) string {
	parseURL, err := urlparser.Parse(url)
	if err != nil {
		return url
	}
	params := urlparser.Values{}
	for _, v := range strings.Split(query, "&") {
		index := strings.Index(v, "=")
		key := v[:index]
		value := v[index+1:]
		value, err = urlparser.QueryUnescape(value)
		if err != nil {
			fmt.Println(err)
		}
		params.Add(key, value)
	}
	if parseURL.RawQuery != "" {
		parseURL.RawQuery += "&" + params.Encode()
	} else {
		parseURL.RawQuery = params.Encode()
	}
	return parseURL.String()
}

func HeadersToString(headers map[string]string) string {
	lines := ""
	for k, v := range headers {
		line := k + ": " + v + "\n"
		lines += line
	}
	return lines
}

func WriteHTML(filename string) {
	prefix := `<!DOCTYPE html>
	<html>
	<meta charset="utf-8">
	<title>AliveWeb Report</title>
	<head>
		<style type="text/css">
			table.hovertable {
				font-family: verdana,arial,sans-serif;
				font-size:11px;
				color:#333333;
				border-width: 1px;
				border-color: #999999;
				border-collapse: collapse;
			}
			table.hovertable th {
				background-color:#c3dde0;
				border-width: 1px;
				padding: 8px;
				border-style: solid;
				border-color: #a9c6c9;
			}
			table.hovertable tr {
				background-color:#d4e3e5;
			}
			table.hovertable td {
				border-width: 1px;
				padding: 8px;
				border-style: solid;
				border-color: #a9c6c9;
			}
		</style>
	</head>
	<body>
	`
	f, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(prefix)
}
