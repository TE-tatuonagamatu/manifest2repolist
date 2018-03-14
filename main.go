package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type remote struct {
	Name string	`xml:"name,attr"`
	Fetch string `xml:"fetch,attr"`
}

type defaultx struct {
	Revision string `xml:"revision,attr"`
	Remote string `xml:"remote,attr"`
}

type project struct {
	Groups string `xml:"groups,attr"`
	Name string `xml:"name,attr"`
	Path string `xml:"path,attr"`
}
type manifest struct {
	Remote   remote `xml:"remote"`
	Default defaultx	`xml:"default"`
	Projects []project `xml:"project"`
}

var m map[string]bool
var uniq []string

func parseRepo(xmldoc []byte) error {
	mt := manifest{}
	err := xml.Unmarshal(xmldoc, &mt)
	if err != nil {
		return err
	}
	for _, project := range mt.Projects {
		if !m[project.Name] {
			m[project.Name] = true
			uniq = append(uniq, project.Name)
		}
	
		//fmt.Printf("%s\n", project.Name)
	}
	return nil
}

func manifest2repolist(dirpath string, finfo os.FileInfo) error {
	f, err := os.Open(filepath.Join(dirpath, finfo.Name()))
	if err != nil {
		return err
	}
	defer f.Close()
	data := make([]byte, finfo.Size())
	_, err = f.Read(data)
	if err != nil {
		return err
	}

	return parseRepo(data)
}

func main() {
	m = make(map[string]bool)
	uniq = [] string{}

	dirPath := os.Args[1]
    files, err := ioutil.ReadDir(dirPath)
    if err != nil {
        log.Fatal(err)
	}

	for _, file := range files {
		if (strings.HasSuffix(file.Name(), ".xml")) {
			err := manifest2repolist(dirPath, file)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, name := range uniq {
		fmt.Printf("%s\n", name)
	}
}
