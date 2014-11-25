package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const baseCheckURL = "https://github.com/torvalds/linux/releases/tag/"
const version = "v3.18"

func handleKernelRelease(w http.ResponseWriter, r *http.Request) {
	checkURL := fmt.Sprintf("%s%s", baseCheckURL, version)
	if testUrl(checkURL) {
		io.WriteString(w, fmt.Sprintf("%s Released", version))
	} else {
		io.WriteString(w, fmt.Sprintf("%s Not released yet", version))
	}
}

func testUrl(url string) bool {
	r, err := http.Head(url)
	if err != nil {
		log.Print(err)
		return false
	}
	return r.StatusCode == http.StatusOK
}

func main() {
	http.HandleFunc("/", handleKernelRelease)
	http.ListenAndServe(":8000", nil)
}
