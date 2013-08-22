package main

import (
	"net/http"
	"os/exec"
)

const (
	shepPath = "/Users/georgemacrorie/personal/shep.git"
)

func infoRefsRecievePackHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	rw.Header().Add("Pragma", "no-cache")
	rw.Header().Add("Cache-Control", "no-cache, max-age=0, must-revalidate")
	rw.Header().Add("Content-Type", "application/x-git-receive-pack-result")
}

func receivePackHandler(rw http.ResponseWriter, req *http.Request) {

}

func main() {

	cmd := gitReceivePackCommand()

	http.HandleFunc("/_git/info/refs", infoRefsRecievePackHandler)
	http.HandleFunc("/_git/git-receive-pack", receivePackHandler)

	http.ListenAndServe(":9292", nil)
}
