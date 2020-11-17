package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/karrick/godirwalk"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	countStr   = "shows count number of files"
	dirStr     = "do not shows size of directories"
	errorStr   = "Error: %s"
	reverseStr = "shows files in ascending order"
	versionStr = "prints version"
)

var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

// FileInfo allows us to store the nested file size
type FileInfo struct {
	name  string
	size  int64
	isDir bool
	path  string
}

// Config has all the adjustable configurations for psize
type Config struct {
	bar       string
	dirSize   bool
	pathname  string
	reverse   bool
	showCount int
	width     int
	writer    io.Writer
}

func defaultConfigs() Config {
	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}

	config := Config{bar: "█",
		dirSize:   true,
		pathname:  "./",
		reverse:   false,
		showCount: 10,
		width:     width,
		writer:    os.Stdout,
	}
	return config
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sum(c Config, files []FileInfo) int64 {
	var total int64
	for _, f := range files {
		total += f.size
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

func getDirSize(c Config, path string) (int64, error) {
	var size int64
	err := godirwalk.Walk(c.pathname+path, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			st, err := os.Stat(osPathname)
			if err != nil {
				return err
			}
			if !de.IsDir() {
				size += st.Size()
			}
			return err
		},
	})
	return size, err
}

func renderBar(c Config, padding int, totalSize int, curSize int) string {
	curWidth := c.width - padding

	fraction := (float64(curSize) / float64(totalSize))
	barSize := int(fraction * float64(curWidth))

	barString := fmt.Sprintf(strings.Repeat(c.bar, barSize))
	return barString
}

func populateFiles(c Config, files []os.FileInfo) []FileInfo {
	retFiles := make([]FileInfo, len(files))

	wg := sync.WaitGroup{}

	for i, f := range files {
		retFiles[i].name = f.Name()
		if c.dirSize {
			wg.Add(1)
			go func(i int) {
				retFiles[i].size, _ = getDirSize(c, files[i].Name())
				wg.Done()
			}(i)
		} else {
			retFiles[i].size = f.Size()
		}
		retFiles[i].isDir = f.IsDir()
	}
	wg.Wait()
	// for i := 0; i < len(files); i++ {
	// 	fmt.Println(retFiles[i].size)
	// }
	return retFiles
}

func sortFiles(c Config, files *[]FileInfo) {
	sort.Slice(*files, func(i, j int) bool {
		if (*files)[i].size < (*files)[j].size {
			if c.reverse {
				return true
			}
			return false
		}
		if (*files)[i].size > (*files)[j].size {
			if c.reverse {
				return false
			}
			return true
		}
		return (*files)[i].name < (*files)[j].name
	})
}

func ls(c Config) string {
	byteString := ""

	fileinfos, err := ioutil.ReadDir(c.pathname)
	if err != nil {
		log.Fatalf(errorStr, err)
	}

	files := populateFiles(c, fileinfos)
	sortFiles(c, &files)

	totalSize := int(sum(c, files))
	var curSize int64
	for _, f := range files[0:min(c.showCount, len(files))] {
		curSize = f.size

		curString := fmt.Sprintf("%-40s %5s|", shortenString(f.name), HumanFileSize(float32(curSize)))
		curLen := len([]rune(curString))
		curString += renderBar(c, curLen, totalSize, int(f.size)) + "\n"
		if f.isDir {
			curString = cyan(curString)
		} else {
			curString = yellow(curString)
		}
		byteString += curString
	}
	sizeString := fmt.Sprintf("Total Size: %v", HumanFileSize(float32(totalSize)))
	byteString += sizeString
	return byteString
}

func main() {

	dirSize := flag.Bool("dirsize", false, dirStr)
	flag.BoolVar(dirSize, "d", false, dirStr)

	reverse := flag.Bool("reverse", false, reverseStr)
	flag.BoolVar(reverse, "r", false, reverseStr)

	showCount := flag.Int("count", 10, countStr)
	flag.IntVar(showCount, "c", 10, countStr)

	version := flag.Bool("version", false, versionStr)
	flag.BoolVar(version, "v", false, versionStr)

	flag.Parse()

	var byteString string
	config := defaultConfigs()

	if *version {
		byteString = "dsize v 1.0.0"
	} else {
		if *dirSize {
			config.dirSize = false
		}
		if path := flag.Arg(0); path != "" {
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}
			if strings.HasPrefix(path, "~/") {
				usr, _ := user.Current()
				dir := usr.HomeDir
				path = dir
			}
			config.pathname = path
		}
		if *reverse {
			config.reverse = true
		}
		if *showCount != 10 {
			config.showCount = *showCount
		}
		byteString = ls(config)
	}
	byteString += "\n"
	io.WriteString(config.writer, byteString)
}
