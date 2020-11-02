package fmtbyte

import "fmt"

const (
	BYTES = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
)

func convert(size int) string {
	byteString := ""
	if size >= GIGABYTE {
		size /= 10
		byteString = "G"
	} else if size >= MEGABYTE {
		size /= 10
		byteString = "M"
	} else if size >= KILOBYTE {
		size /= 10
		byteString = "K"
	} else {
		byteString = "B"
	}
	output := fmt.Sprint("%d%s", size, byteString)
	return output
}
