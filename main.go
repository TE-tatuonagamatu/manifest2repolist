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
	Name  string `xml:"name,attr"`
	Fetch string `xml:"fetch,attr"`
}

type defaultx struct {
	Revision string `xml:"revision,attr"`
	Remote   string `xml:"remote,attr"`
}

type project struct {
	Groups string `xml:"groups,attr"`
	Name   string `xml:"name,attr"`
	Path   string `xml:"path,attr"`
}
type manifest struct {
	Remote   remote    `xml:"remote"`
	Default  defaultx  `xml:"default"`
	Projects []project `xml:"project"`
}

func parseRepo(xmldoc []byte) ([]string, error) {
	mt := manifest{}
	err := xml.Unmarshal(xmldoc, &mt)
	if err != nil {
		return nil, err
	}
	list := []string{}
	for _, project := range mt.Projects {
		list = append(list, project.Name)
	}
	return list, nil
}

func manifest2repolist(dirpath string, finfo os.FileInfo) ([]string, error) {
	f, err := os.Open(filepath.Join(dirpath, finfo.Name()))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data := make([]byte, finfo.Size())
	_, err = f.Read(data)
	if err != nil {
		return nil, err
	}

	return parseRepo(data)
}

func uniq(repolist []string) []string {
	m := make(map[string]bool)
	uniq := []string{}

	for _, l := range repolist {
		if !m[l] {
			m[l] = true
			uniq = append(uniq, l)
		}
	}
	return uniq
}

func main() {

	dirPath := os.Args[1]
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	repolist := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".xml") {
			list, err := manifest2repolist(dirPath, file)
			if err != nil {
				log.Fatal(err)
			}
			repolist = append(repolist, list...)
		}
	}

	repolist = uniq(repolist)
	for _, name := range repolist {
		fmt.Printf("%s\n", name)
	}
}
