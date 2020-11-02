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
)

const (
	versionStr   = "prints version"
	errorStr     = "Error: %s"
	total_blocks = 100
)

func sum(files []os.FileInfo) int {
	total := 0
	for _, f := range files {
		total += int(f.Size())
	}
	return total
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
	total_size := sum(files)
	for _, f := range files {
		fraction := (float64(f.Size()) / float64(total_size))
		cur_size := int(fraction * total_blocks)
		if f.IsDir() {
			color.Set(color.FgCyan)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Printf("%-40s %5v %-50s\n", f.Name(), f.Size(), strings.Repeat("█", cur_size))
	}
	color.Unset()
	fmt.Printf("Total Size: %v\n", total_size)
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
