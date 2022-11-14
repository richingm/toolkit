package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

var (
	TOOLKIT = "toolkit"
)

func checkFile() (string, string) {
	var dstDir string
	flag.StringVar(&dstDir, "d", "", "set build directory")
	var config string
	pwd, _ := os.Getwd()
	file := fmt.Sprintf("%s/%s", pwd, "conf.yml")
	flag.StringVar(&config, "c", file, "Build requires dependent yaml files")
	flag.Parse()
	if dstDir == "" {
		panic("please set dst dir")
	}
	makeDir(dstDir)
	if _, err := os.Stat(config); os.IsNotExist(err) {
		panic(err)
	}
	return dstDir, config
}

func main() {
	dstDir, config := checkFile()

	// read file
	yamlContext, err := ioutil.ReadFile(config)
	if err != nil {
		panic(err)
	}

	// to yaml struct
	c := &Configuration{}
	err = yaml.Unmarshal(yamlContext, c)
	if err != nil {
		panic(err)
	}

	// check
	if c.App != TOOLKIT {
		panic("fatal: not toolkit yaml")
	}

	// 创建README.md文件
	createFiles(dstDir, c.Files)

	// 创建目录
	makeDirs(dstDir, c.Dirs)

	fmt.Println("Success")
}

// create dirs
func makeDirs(currentDir string, dirs []Dirs) {
	if len(dirs) == 0 {
		return
	}
	for _, dir := range dirs {
		var build strings.Builder
		build.WriteString(currentDir)
		build.WriteString(dir.Name)
		makeDir(build.String())
		makeDirs(build.String(), dir.Dirs)
		createFiles(build.String(), dir.Files)
	}
}

func makeDir(dstDir string) {
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err := os.Mkdir(dstDir, 0755)
		if err != nil {
			panic(err)
		}
		err = os.Chmod(dstDir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

// create files
func createFiles(currentDir string, files []File) {
	if len(files) == 0 {
		return
	}
	for _, file := range files {
		var build strings.Builder
		build.WriteString(currentDir)
		build.WriteString("/")
		build.WriteString(file.Name)
		filePath := build.String()
		_, err := os.Stat(filePath)
		if err == nil {
			os.Remove(filePath)
		}
		createFile(filePath, file)
	}
}

// create file
func createFile(filePath string, file File) {
	handle, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer handle.Close()
	if len(file.Info) != 0 {
		writer := bufio.NewWriter(handle)
		for _, s := range file.Info {
			writer.WriteString(s)
			writer.WriteString("\n")
		}
		writer.Flush()
	}
}

type Configuration struct {
	App   string
	Files []File
	Dirs  []Dirs
}

type Dirs struct {
	Name  string
	Files []File
	Dirs  []Dirs
}

type File struct {
	Name string
	Info []string
}
