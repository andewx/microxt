package common

import (
	"log"
	"math/rand"
	"os"
	"strings"
)

const APP_NAME = "microxt"

func GetProjectDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get project working directory")
	}
	paths := strings.Split(wd, PATH_SEPARATOR)
	i := 0
	j := 0
	for i = 0; i < len(paths); i++ {
		cur := paths[i]
		if cur == APP_NAME {
			if i < len(paths)-1 {
				j = i + 1
				break
			} else {
				j = i
				break
			}
		}
	}
	return strings.Join(paths[0:j], PATH_SEPARATOR)
}

func ProjectRelativePath(relative_path string) string {
	return GetProjectDir() + PATH_SEPARATOR + relative_path
}

// Linux Cd command
func Cd(path string) string {
	paths := strings.Split(path, PATH_SEPARATOR)
	return strings.Join(paths[0:len(paths)-1], PATH_SEPARATOR)
}

const (
	BIG_ENDIAN = iota
	LITTLE_ENDIAN
)

func RandUint32() uint32 {
	return rand.Uint32()
}

func Int32(bytes []byte, endianess int) int32 {
	var n int32
	if endianess == BIG_ENDIAN {
		n = int32(bytes[0])<<24 + int32(bytes[1])<<16 + int32(bytes[2])<<8 + int32(bytes[3])
	} else {
		n = int32(bytes[3])<<24 + int32(bytes[2])<<16 + int32(bytes[1])<<8 + int32(bytes[0])
	}
	return n
}

func Int16(bytes []byte, endianess int) int16 {
	var n int16
	if endianess == BIG_ENDIAN {
		n = int16(bytes[0])<<8 + int16(bytes[1])
	} else {
		n = int16(bytes[1])<<8 + int16(bytes[0])
	}
	return n
}

func Int8(bytes []byte) int8 {
	return int8(bytes[0])
}

func GetBytes32(x int, endianess int) []byte {
	if endianess == BIG_ENDIAN {
		return []byte{byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)}
	} else {
		return []byte{byte(x), byte(x >> 8), byte(x >> 16), byte(x >> 24)}
	}
}

func GetBytes16(x int, endianess int) []byte {
	if endianess == BIG_ENDIAN {
		return []byte{byte(x >> 8), byte(x)}
	} else {
		return []byte{byte(x), byte(x >> 8)}
	}
}

func GetBytes8(x int) []byte {
	return []byte{byte(x)}
}
