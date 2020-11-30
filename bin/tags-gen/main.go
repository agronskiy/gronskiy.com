// Reads the "posts" directory recursively, analyses YAML frontmatter, extracts tags
// and then generates necessary family of files into "tags" directory

// TODO(agronskiy, 2020-11): create a "Parallelizer" package, opensource

package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"text/template"

	front "github.com/gernest/front"
)

type (
	input  string
	output string
)

type listTemplateData struct {
	Tag     string
	Lang    string
	PageURL string
}

const listYamlTemplate = `---
title: "Alexey Gronskiy's posts :: tags: {{.Tag}} :: languages: {{.Lang}} "
tag: {{.Tag}}
lang: {{.Lang}}
{{- if ne .PageURL "" }}
url: {{ .PageURL }}
{{- end }}
mathjax: true
---
<!-- Generated automatically -->`

type tagsTemplateData struct {
	TagsList []string
}

const tagsYamlTemplate = `---
tags_all:
{{- range $tag := .TagsList}}
  - {{ $tag }}
{{- end}}
---
<!-- Generated automatically -->`

const (
	dataDir = "../../data"
	tagsDir = "../../content/tags"
)

func makeRunner() (chan<- input, <-chan output) {
	var (
		numCPU = runtime.NumCPU()

		inputQueue        = make(chan input, numCPU)
		intermediateQueue = make(chan output, numCPU)
		outputQueue       = make(chan output, numCPU)

		counterCh = make(chan int)

		numOpenWorkers = 0
	)

	stopGracefully := func() {
		for {
			// First, drain remaining results, and only then stop.
			select {
			case out := <-intermediateQueue:
				outputQueue <- out
			default:
				close(outputQueue)
				return
			}
		}
	}

	// Create actual runner. It:
	// 1. spawns workers
	// 2a. listens to their output
	// 2b. does bookkeeping, counts how many are still working and closes output channel
	go func() {
		for i := 0; i < numCPU; i++ {
			go worker(inputQueue, intermediateQueue, counterCh)
		}

		for {
			select {
			case out := <-intermediateQueue:
				outputQueue <- out
			case n := <-counterCh:
				numOpenWorkers += n
				if numOpenWorkers > 0 {
					continue
				}
				stopGracefully()
			}
		}
	}()

	return inputQueue, outputQueue
}

func worker(
	inputQueue <-chan input,
	intermediateQueue chan<- output,
	counterCh chan<- int,
) {
	counterCh <- 1
	defer func() {
		counterCh <- -1
	}()

	for path := range inputQueue {

		m := front.NewMatter()
		m.Handle("---", front.YAMLHandler)
		file, err := os.Open(string(path))
		f, _, err := m.Parse(file)
		if err != nil {
			panic(err)
		}

		if unlisted, ok := f["unlisted"]; ok {
			switch unlisted := unlisted.(type) {
			case bool:
				if unlisted {
					return
				}
			case string:
				if unlisted == "true" {
					return
				}
			default:
				log.Printf("Tag `unlisted` in %v unrecognized!")
			}
		}

		if tags, ok := f["tags"]; ok {
			switch tags := tags.(type) {
			case []interface{}:
				for _, tag := range tags {
					tag := tag.(string)
					intermediateQueue <- output(tag)
				}
			case string:
				intermediateQueue <- output(tags)
			default:
			}
		}
	}
}

func processOutput(outputQueue <-chan output) {
	t := template.Must(template.New("yaml").Parse(listYamlTemplate))

	os.MkdirAll(filepath.Join(tagsDir), os.ModePerm)
	file, _ := os.Create(filepath.Join(tagsDir, "_index.md"))
	t.Execute(file, listTemplateData{Tag: "all", Lang: "all", PageURL: "posts"})

	for _, lang := range []string{"en", "ru"} {
		os.MkdirAll(filepath.Join(tagsDir, lang), os.ModePerm)
		file, _ := os.Create(filepath.Join(tagsDir, lang, "_index.md"))
		t.Execute(file, listTemplateData{Tag: "all", Lang: lang, PageURL: filepath.Join("/posts", lang)})
	}

	tags := make(map[string]bool)

	// This will create files under 'tags/...'
	for res := range outputQueue {
		tag := string(res)
		if tags[tag] {
			continue
		}

		tags[tag] = true

		os.MkdirAll(filepath.Join(tagsDir, tag), os.ModePerm)
		file, _ := os.Create(filepath.Join(tagsDir, tag, "_index.md"))
		t.Execute(file, listTemplateData{Tag: tag, Lang: "all"})

		for _, lang := range []string{"en", "ru"} {
			os.MkdirAll(filepath.Join(tagsDir, tag, lang), os.ModePerm)
			file, _ := os.Create(filepath.Join(tagsDir, tag, lang, "_index.md"))
			t.Execute(file, listTemplateData{Tag: tag, Lang: lang})
		}
	}

	// This will create all tags under data.
	tagsList := make([]string, 0, len(tags))
	for k := range tags {
		if tags[k] {
			tagsList = append(tagsList, k)
		}
	}
	sort.Strings(tagsList)

	dataT := template.Must(template.New("yaml").Parse(tagsYamlTemplate))

	os.MkdirAll(filepath.Join(tagsDir), os.ModePerm)
	tagsFile, _ := os.Create(filepath.Join(dataDir, "tags_all.yaml"))
	dataT.Execute(tagsFile, tagsTemplateData{TagsList: tagsList})

	log.Println("All tags written.")
	log.Println(tagsList)
}

func main() {
	inputQueue, outputQueue := makeRunner()

	go func() {
		// Rather unoptimal but ok
		paths := make([]string, 0)
		pathsGlob, err := filepath.Glob("../../content/posts/*/*.md")
		if err != nil {
			return
		}
		paths = append(paths, pathsGlob...)

		pathsGlob, err := filepath.Glob("../../content/posts/*/*/*.md")
		if err != nil {
			return
		}
		paths = append(paths, pathsGlob...)

		pathsGlob, err := filepath.Glob("../../content/posts/*/*/*/*.md")
		if err != nil {
			return
		}
		paths = append(paths, pathsGlob...)

		for _, path := range paths {
			inputQueue <- input(path)
		}

		paths, err = filepath.Glob("../../content/posts/*/*/*.md")
		if err != nil {
			return
		}
		for _, path := range paths {
			inputQueue <- input(path)
		}

		close(inputQueue)
	}()

	processOutput(outputQueue)
}
