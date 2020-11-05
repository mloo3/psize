package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	versionStr   = "prints version"
	errorStr     = "Error: %s"
	total_blocks = 40
	bar          = "█"
)

type State struct {
	bar         string
	totalBlocks int
}

var curState = State{bar: "█"}

func sum(files []os.FileInfo) int {
	total := 0
	for _, f := range files {
		total += int(f.Size())
	}
	return total
}

func shortenString(s string) string {
	strLength := len([]rune(s))
	if strLength > 30 {
		return "..." + s[strLength-30:]
	}
	return s
}

func renderBar(padding int, totalSize int, curSize int) string {
	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}
	width -= padding

	fraction := (float64(curSize) / float64(totalSize))
	barSize := int(fraction * float64(width))

	barString := fmt.Sprintf(strings.Repeat(bar, barSize))
	return barString
}
func ls(pathname string) {
	files, err := ioutil.ReadDir(pathname)
	if err != nil {
		log.Fatalf(errorStr, err)
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].Size() < files[j].Size() {
			return false
		}
		if files[i].Size() > files[j].Size() {
			return true
		}
		return files[i].Name() < files[j].Name()
	})
	totalSize := sum(files)
	for _, f := range files {
		// fraction := (float64(f.Size()) / float64(totalSize))
		// cur_size := int(fraction * total_blocks)
		if f.IsDir() {
			color.Set(color.FgCyan)
		} else {
			color.Set(color.FgRed)
		}
		curString := fmt.Sprintf("%-40s %5s|", shortenString(f.Name()), HumanFileSize(int(f.Size())))
		curLen := len([]rune(curString))
		curString += renderBar(curLen, totalSize, int(f.Size()))
		fmt.Println(curString)
	}
	color.Unset()
	fmt.Printf("Total Size: %v\n", HumanFileSize(totalSize))
}

func main() {
	version := flag.Bool("version", false, versionStr)
	flag.BoolVar(version, "v", false, versionStr)

	flag.Parse()

	pathname := "./"

	if *version {
		fmt.Println("dsize v 1.0.0")
	} else if len(flag.Args()) == 1 {
		pathname = flag.Args()[0]
		ls(pathname)
	} else {
		ls(pathname)
	}
}
