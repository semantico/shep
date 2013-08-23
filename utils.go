package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

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

func IntToHexString4(i int) string {
	bytes, err := encodeIntToBytes(i)
	if err != nil {
		log.Fatal(err)
	}
	return padHexString(hex.EncodeToString(bytes))
}
