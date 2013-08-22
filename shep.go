package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const (
	shepPath = "/Users/georgemacrorie/personal/shep.git"
)

func generateServiceMessage() string {
	packet := "# service=git-receive-pack"
	prefix := intToHexString4(len(packet) + 4)
	return fmt.Sprintf("%s%s0000", prefix, packet)
}

func infoRefsRecievePackHandler(rw http.ResponseWriter, req *http.Request) {
	setCommonHeadersOnResponse(rw)
	rw.Header().Add("Content-Type", "application/x-git-receive-pack-advertisement")

	strings.NewReader(generateServiceMessage()).WriteTo(rw)

	cmd := exec.Command("git", "receive-pack", "--stateless-rpc", "--advertise-refs", shepPath)

	cmd.Stdout = rw

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()
}

func receivePackHandler(rw http.ResponseWriter, req *http.Request) {
	setCommonHeadersOnResponse(rw)
	rw.Header().Add("Content-Type", "application/x-git-receive-pack-advertisement")

	cmd := exec.Command("git", "receive-pack", "--stateless-rpc", shepPath)

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

func padHexString(hex string) string {
	if len(hex) >= 4 {
		return hex
	}
	padding := ""
	paddingLength := 4 - len(hex)
	for i := 0; i < paddingLength; i++ {
		padding += "0"
	}
	return padding + hex
}

func encodeIntToBytes(i int) (bytes []byte, err error) {
	temp := make([]byte, 100)
	length := binary.PutUvarint(temp, uint64(i))
	bytes = make([]byte, length)
	if lost := copy(bytes, temp); lost == length {
		err = nil
	} else {
		bytes = temp
		err = errors.New("Bytes got lost when copying: " + fmt.Sprintf("%d", lost))
	}
	return
}

func intToHexString4(i int) string {
	bytes, _ := encodeIntToBytes(i)
	return padHexString(hex.EncodeToString(bytes))
}

func main() {
	http.HandleFunc("/_git/info/refs", infoRefsRecievePackHandler)
	http.HandleFunc("/_git/git-receive-pack", receivePackHandler)

	http.ListenAndServe(":9292", nil)
}
