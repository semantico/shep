package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var (
	shepPath = flag.String("g", "/Users/georgemacrorie/personal/shep.git", "Shep bare git repo")
	port     = flag.String("p", "9292", "Port on which shep listens")
)

func generateServiceMessage() string {
	packet := "# service=git-receive-pack"
	prefix := IntToHexString4(len(packet) + 4)
	return fmt.Sprintf("%s%s0000", prefix, packet)
}

func infoRefsRecievePackHandler(rw http.ResponseWriter, req *http.Request) {
	setCommonHeadersOnResponse(rw)
	rw.Header().Add("Content-Type", "application/x-git-receive-pack-advertisement")

	strings.NewReader(generateServiceMessage()).WriteTo(rw)

	cmd := exec.Command("git", "receive-pack", "--stateless-rpc", "--advertise-refs", *shepPath)

	cmd.Stdout = rw

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()
}

func receivePackHandler(rw http.ResponseWriter, req *http.Request) {
	setCommonHeadersOnResponse(rw)
	rw.Header().Add("Content-Type", "application/x-git-receive-pack-advertisement")

	cmd := exec.Command("git", "receive-pack", "--stateless-rpc", *shepPath)

	cmd.Stdin = req.Body
	cmd.Stdout = rw

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	cmd.Wait()
}

func setCommonHeadersOnResponse(rw http.ResponseWriter) {
	rw.Header().Add("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	rw.Header().Add("Pragma", "no-cache")
	rw.Header().Add("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

func main() {
	flag.Parse()

	http.HandleFunc("/_git/info/refs", infoRefsRecievePackHandler)
	http.HandleFunc("/_git/git-receive-pack", receivePackHandler)

	fileSystem, err := NewGitFileSystem(*shepPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "favicon not found", 404)
	})
	http.Handle("/", http.FileServer(fileSystem))

	http.ListenAndServe(":"+*port, nil)
}
