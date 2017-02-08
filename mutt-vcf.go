package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func scan(filename string, bufferSize int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var (
		buffer []byte = make([]byte, bufferSize)
		offset int
		key    string
		vCard  map[string]string
	)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		for i, j := 0, 0; j < n; j++ {
			if buffer[j] == 0x3A {
				key = string(buffer[i:j])
				if key == "BEGIN" {
					vCard = make(map[string]string)
				}
				i = j + 1
			} else if buffer[j] == 0x0A {
				vCard[key] = string(buffer[i : j-1])
				if key == "END" {
					alias(vCard["FN"], vCard["EMAIL;TYPE=INTERNET"])
					alias(vCard["FN"], vCard["EMAIL;TYPE=INTERNET;TYPE=HOME"])
					alias(vCard["FN"], vCard["EMAIL;TYPE=INTERNET;TYPE=WORK"])
				}
				i = j + 1
				offset = i
			}
		}
		_, err = file.Seek(int64(offset-n), io.SeekCurrent)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func alias(nickname /* longname, */, address string) {
	if address != "" {
		if nickname == "" {
			nickname = address
		}
		fmt.Printf("alias \"%v\" <%v>\n", nickname, address)
	}
}

func main() {
	bufferSize := flag.Int("bytes", 1024, "The buffer size in bytes.")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage:\n\t$", os.Args[0], "[-bytes 1024] /path/to/file.vcf [<< /path/to/mutt/aliases]")
		return
	}
	scan(flag.Arg(0), *bufferSize)
}
