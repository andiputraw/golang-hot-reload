package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"plugin"
	"sort"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type DrawPlugin interface {
	Draw(screen *ebiten.Image)
}

type TmpStruct struct {}

func getAnyStruct() interface{}{
	return TmpStruct{}
}

func loadPlugin[T DrawPlugin](src string) (T, error) {
	plug, err := plugin.Open(src)
		
	if err != nil {
		tmp, _  := getAnyStruct().(T)
		return tmp, err
	}
	pluginLookup , err := plug.Lookup("Export")
	
	if err != nil {
		tmp, _  := getAnyStruct().(T)

		return tmp, errors.New("fail to lookup symbol")
	}
	
	pluginSymbol, ok := pluginLookup.(T)
	
	if(!ok){
		return pluginSymbol, errors.New("fail to cast struct")
	}

	return pluginSymbol, nil
}

func filter_dir(files []fs.FileInfo)[]fs.FileInfo{
	filtered_files := make([]fs.FileInfo, 0)

	for _, file := range files {
		if strings.Contains(file.Name(), "draw"){
			filtered_files = append(filtered_files, file)		
		} 
	}

	return filtered_files
}

func getNewestPluginPath() string{
	dirSrc := "./plugin"
	dir, err := ioutil.ReadDir(dirSrc)

	if err != nil {
		log.Panicf("Cannot read dir: %s", err.Error())
	}

	dir = filter_dir(dir)
	sort.Slice(dir, func(i, j int) bool {
		return dir[i].Name() > dir[j].Name()
	})
	plugPath := dir[0]
	return fmt.Sprintf("%s/%s", dirSrc, plugPath.Name())
}

type Game struct{
	drawPlugin DrawPlugin
	drawPluginPath string
	drawPluginLastUpdated time.Time

}

func (g *Game) Update() error {
	fileInfo, _ := os.Stat("./plugin")
	currentModifiedTime := fileInfo.ModTime()
	
	if(currentModifiedTime != g.drawPluginLastUpdated){
		plugPath := getNewestPluginPath()

		if plugPath == g.drawPluginPath {
			g.drawPluginLastUpdated = currentModifiedTime
			g.drawPluginPath = plugPath			
			return nil
		}
		draw, err := loadPlugin[DrawPlugin](plugPath)
		if err != nil {
			log.Printf("Fail to reload plugin, no reload will happened %s", err.Error())
			g.drawPluginLastUpdated = currentModifiedTime
			g.drawPluginPath = plugPath
			return nil
		}
		log.Printf("Reloaded plugin... %s", plugPath)
		g.drawPlugin = draw
		g.drawPluginPath = plugPath
		g.drawPluginLastUpdated = currentModifiedTime
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawPlugin.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}


func main(){
	plugPath := getNewestPluginPath()
	dirSrc := "./plugin"
	
	dirInfo, _ := os.Stat(dirSrc)
	lastUpdated := dirInfo.ModTime()
	
	drawPlugin, err := loadPlugin[DrawPlugin](plugPath)
	if err != nil {
		log.Fatalf("Fail to load plugin %s", err.Error())
	}
	log.Printf("Loaded plugin... %s", plugPath)

	// timeout := 0
	// for(timeout < 100){
	

	// 	time.Sleep(3 * time.Second)
	// 	timeout += 1
	// }

	game := Game{
		drawPlugin: drawPlugin,
		drawPluginLastUpdated: lastUpdated,
		drawPluginPath: plugPath,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
	
}