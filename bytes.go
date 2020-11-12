package main

import "fmt"

// Sizes of bytes
const (
	BYTES = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
)

// HumanFileSize returns a readable file size string
func HumanFileSize(size int) string {
	var byteString string
	if size >= GIGABYTE {
		size /= GIGABYTE
		byteString = "G"
	} else if size >= MEGABYTE {
		size /= MEGABYTE
		byteString = "M"
	} else if size >= KILOBYTE {
		size /= KILOBYTE
		byteString = "K"
	} else {
		byteString = "B"
	}
	output := fmt.Sprintf("%4d%2s", size, byteString)
	return output
}
