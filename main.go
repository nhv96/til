package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

var (
	blogPostFormat    = "%s-%s"
	blogRepositoryDir string
	postCategory      string
	postLayout        string

	ignoreFiles = []string{"README.md"}

	timeLayout = "2006-01-02"

	docRegex = "\\S+.(md|markdown)"

	postTemplate = `---
layout: {{ .Layout }}
title:  {{printf "%q" .Title }}
date:   {{ .Date }}
category: {{ .Category }}
---
{{ .Content }}
`
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	blogRepositoryDir = os.Getenv("BLOG_DIR")
	postCategory = os.Getenv("POST_CATEGORY")
	postLayout = os.Getenv("POST_LAYOUT")

	ignoreFiles = append(ignoreFiles, blogRepositoryDir)

	log.Println(fmt.Sprintf(`
	Blog dir: %s
	Post category: %s
	Post layout: %s
	`, blogRepositoryDir, postCategory, postLayout))
}

type Post struct {
	FileName string
	Layout   string
	Title    string
	Date     string
	Category string
	Content  string
}

func main() {
	r := regexp.MustCompile(docRegex)

	fileMap := make(map[string]*Post)

	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.HasPrefix(path, ".git") || strings.HasPrefix(path, blogRepositoryDir) {
			return nil
		}

		if slices.Contains(ignoreFiles, info.Name()) {
			return nil
		}

		if err == nil && r.MatchString(info.Name()) {
			post, err := fileInfoToBlogPost(path, info)
			if err != nil {
				return err
			}

			fileMap[post.FileName] = post
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(fileMap) == 0 {
		log.Println("Found no new posts. Exit.")
		os.Exit(0)
	}

	err = writeToBlog(fileMap)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func fileInfoToBlogPost(path string, fi fs.FileInfo) (*Post, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	modTime := fi.ModTime().Format(timeLayout)

	blogFileName := fmt.Sprintf(blogPostFormat, modTime, fi.Name())

	segs := strings.Split(fi.Name(), ".")

	name := strings.ReplaceAll(segs[0], "-", " ")

	name = strings.ToTitle(name)

	p := &Post{
		Title:    name,
		Date:     fi.ModTime().Format(timeLayout),
		Category: postCategory,
		Layout:   postLayout,
		FileName: blogFileName,
		Content:  string(content),
	}

	return p, nil
}

func writeToBlog(posts map[string]*Post) error {
	t := template.Must(template.New("blogpost").Parse(postTemplate))

	for postPath, post := range posts {
		// the repository dir should already existed
		file, err := os.Create(fmt.Sprintf("%s/%s", blogRepositoryDir, postPath))
		if err != nil {
			return err
		}

		err = t.Execute(file, post)
		if err != nil {
			return err
		}
	}
	return nil
}
