package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func compileDraw() string {
	draw_path := fmt.Sprintf("./plugin/%d_draw.so", time.Now().Unix())

	// 			  		 go    build    -buildmode=plugin   -o   .            ./src/draw/draw.go
	err := exec.Command("go", "build", "-buildmode=plugin", "-o", draw_path, "./src/draw/draw.go").Run()
	if err != nil {
		log.Fatalf("Failed to build plugin. %s", err.Error())
	}
	log.Println("Rebuilding, ", draw_path)
	return draw_path
}

func filterDir(dir []fs.FileInfo, prefix string) []fs.FileInfo {
	filtered := make([]fs.FileInfo, 0)
	for _, file := range dir {
		if strings.Contains(file.Name(), prefix) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func deleteDir(path,prefix, except string){
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Cannot read %s : %s", path, err.Error())
	}
	prefix_dir := filterDir(dir, prefix)
	for _ , file := range prefix_dir {
		filePath := fmt.Sprintf("%s/%s", path, file.Name())
		if strings.Contains(filePath, except){
			continue
		}
		err := os.Remove(filePath)
		if err != nil {
			log.Panicf("Cannot delete %s: %s", filePath, err.Error())
		}
	}
}

func getPrefix(name string) string {
	return strings.Split(name, "/")[1]
}

func main(){
	_ = compileDraw()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Cannot create watcher: %s", err.Error())
	}
	lastBuild := time.Now()

	go func(){
		for {
			select {
			case event,ok := <- watcher.Events:
				if !ok {
					return
				}
				if time.Since(lastBuild) < 3 * time.Second {
					continue
				}
				var except string
				if event.Has(fsnotify.Write){
					except = compileDraw()
				}
				lastBuild = time.Now()
				deleteDir("./plugin",getPrefix(event.Name), except )
			case err, ok := <- watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Got error from watcher: %s", err.Error() )
			}
		}
	}()
	draw_file := "./src/draw/draw.go"
	err = watcher.Add(draw_file)
	if err != nil {
		log.Fatalf("Cannot watch %s: %s",draw_file, err.Error())
	}

	log.Printf("Watching %s\n", draw_file)

	<-make(chan struct{})
}

