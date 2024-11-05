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
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	godotenv.Load()

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

	posts := make(map[string]*Post)

	// find new/updated documents based on time
	now := time.Now().Format(timeLayout)

	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.HasPrefix(path, ".git") || strings.HasPrefix(path, blogRepositoryDir) {
			return nil
		}

		if slices.Contains(ignoreFiles, info.Name()) {
			return nil
		}

		if err == nil && r.MatchString(info.Name()) && now == info.ModTime().Format(timeLayout) {
			post, err := fileInfoToBlogPost(path, info)
			if err != nil {
				return err
			}

			posts[post.FileName] = post
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(posts) == 0 {
		log.Println("Found no new posts. Exit.")
		os.Exit(0)
	}

	err = writeToBlog(posts)
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

	blogFileName := fi.Name()

	segs := strings.Split(fi.Name(), ".")

	name := strings.ReplaceAll(segs[0], "-", " ")

	cs := cases.Title(language.English)
	name = cs.String(name)

	p := &Post{
		Title:    name,
		Date:     modTime,
		Category: postCategory,
		Layout:   postLayout,
		FileName: blogFileName,
		Content:  string(content),
	}

	return p, nil
}

func writeToBlog(postsToWrite map[string]*Post) error {
	t := template.Must(template.New("blogpost").Parse(postTemplate))

	currentPosts, err := findAllCurrentPosts(blogRepositoryDir)
	if err != nil {
		return nil
	}

	// the blog repository should already existed
	for postName, post := range postsToWrite {
		var (
			file *os.File
			err  error
		)
		// a post already existed
		if fileNameInBlog, exist := currentPosts[postName]; exist {
			// re assign the initial date of the post
			post.Date = fileNameInBlog[:10]

			file, err = os.OpenFile(fmt.Sprintf("%s/%s", blogRepositoryDir, fileNameInBlog), os.O_RDWR, os.ModeAppend)
		} else {
			filename := post.Date + "-" + postName
			file, err = os.Create(fmt.Sprintf("%s/%s", blogRepositoryDir, filename))
		}

		defer file.Close()

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

func findAllCurrentPosts(dir string) (map[string]string, error) {
	currentPosts := map[string]string{}

	r := regexp.MustCompile(docRegex)

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err == nil && r.MatchString(info.Name()) {
			// trim the date prefix YYYY-MM-DD-
			name := info.Name()[11:]

			currentPosts[name] = info.Name()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return currentPosts, nil
}
