package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

const baseCheckURL = "https://github.com/torvalds/linux/releases/tag/"
const version = "v3.18"

func handleKernelRelease(w http.ResponseWriter, r *http.Request) {
	versionArr, ok := r.URL.Query()["version"]
	version := ""
	if ok {
		version = versionArr[0]
	}
	checkURL := fmt.Sprintf("%s%s", baseCheckURL, version)
	if testUrl(checkURL) {
		io.WriteString(w, fmt.Sprintf("%s Released", version))
	} else {
		io.WriteString(w, fmt.Sprintf("%s Not released yet", version))
	}
}

func TLSConfig() (*tls.Config, error) {
	certDir := "/system/etc/security/cacerts"
	fi, err := os.Stat(certDir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%q not a dir", certDir)
	}
	pool := x509.NewCertPool()
	cfg := &tls.Config{RootCAs: pool}

	f, err := os.Open(certDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, _ := f.Readdirnames(-1)
	for _, name := range names {
		pem, err := ioutil.ReadFile(filepath.Join(certDir, name))
		if err != nil {
			return nil, err
		}
		pool.AppendCertsFromPEM(pem)
	}
	return cfg, nil
}

func androidDial(network, addr string) (net.Conn, error) {
	c, err := net.Dial(network, "192.30.252.128:443")
	if err != nil {
		return nil, err
	}
	return c, err
}

func testUrl(url string) bool {
	tlsConfig, err := TLSConfig()
	if err != nil {
		log.Print(err)
	}

	tr := &http.Transport{
		Dial:            androidDial,
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: tr}
	r, err := client.Head(url)
	if err != nil {
		return false
	}
	return r.StatusCode == http.StatusOK
}

func main() {
	http.HandleFunc("/", handleKernelRelease)
	http.ListenAndServe(":8000", nil)
}
