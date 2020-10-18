package main

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
	"strings"
)

var SECTORSIZE int64 = 512

func main() {

	if len(os.Args) > 2 {
		if mode := os.Args[1]; mode != "-mbr" {
			if mode != "-gpt" {
			println("----Usage: go run hw2-1.go [MODES] <filename>")
			println("MODES:\n -mbr     mbr analysis mode\n -gpt     gpt analysis mode" )
			os.Exit(0)
			}
		}
	} else {
		println("+++Usage: go run hw2-1.go [MODES] <filename>")
			println("MODES:\n -mbr     mbr analysis mode\n -gpt     gpt analysis mode" )
			os.Exit(0)
	}
	mode := os.Args[1]
	file := os.Args[2]
	
	// Open File
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if mode == "-mbr" {
		MBRMap := populateType("mbr")
		entrySize := 16
		// Get the content

		buffer := make([]byte, entrySize)
		f.Seek(446,0)

		entry := 1
		numPartions := 0
		for entry < 5 {
			_, err := f.Read(buffer)
			if err != nil {
				break	
			}
			if DecodeHexString(buffer[8:12]) != "00000000" {
				numPartions = entry
			} 
			entry++
		}

		
		fmt.Println("Number of partitions: ", numPartions)
		println("---------------------------------------------")
		
		f.Seek(446,0)
		entry = 1
		for entry < numPartions+1 {
			f.Read(buffer)

			fmt.Println("Partition ", entry)

			chunk := buffer[0:1]
			if DecodeHexString(chunk) == "80" {
				println("Boot: bootable")
			} else {
				println("Boot: non-bootable")
			}

			chunk = buffer[4:5]
			str := DecodeHexString(chunk)
			println("Type: "+ MBRMap[str])

			chunk = buffer[8:12]
			str = DecodeHexString(chunk)
			x := ToLittleEndian(str,4)
			sLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Address LBA: %d\n", sLBA)

			chunk = buffer[12:16]
			str = DecodeHexString(chunk)
			x = ToLittleEndian(str,4)
			eLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Sectors in Partition: %d\n", eLBA)
			println("---------------------------------------------")
			if err != nil {
				panic(err)	
			}
			entry++
		}
	
	} else if mode == "-gpt" {
		
		GPTMap := populateType("gpt")
		entrySize := 128
		
		f.Seek(LBA(2),0)
		buffer := make([]byte, entrySize)
		
		entry := 1
		numPartions := 0
		for entry <= 128 {
			_, err := f.Read(buffer)
			if err != nil {
				break	
			} 
			if DecodeHexString(buffer[0:16]) != "00000000000000000000000000000000" {
				numPartions = entry
			} 
			entry++
		}

		f.Seek(LBA(2),0)
		
		fmt.Println("Number of partitions: ", numPartions)
		println("---------------------------------------------")

		entry = 1
		for entry < numPartions+1 {
			f.Read(buffer)

			fmt.Println("Partition ", entry)

			chunk := buffer[56:128]
			fmt.Printf("Name: %s\n",chunk)

			chunk = buffer[0:16]
			str := DecodeHexString(chunk)
			GUID := ToLittleEndian(str[0:8],4)+"-"+ToLittleEndian(str[8:12],2)+"-"+ToLittleEndian(str[12:16],2)+"-"+str[16:20]+"-"+str[20:]
			println("GUID: "+GUID)
			println("Type: "+ GPTMap[GUID])

			chunk = buffer[32:40]
			str = DecodeHexString(chunk)
			x := ToLittleEndian(str,8)
			sLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Starting LBA: %d\n", sLBA)

			chunk = buffer[40:48]
			str = DecodeHexString(chunk)
			x = ToLittleEndian(str,8)
			eLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Ending LBA: %d\n", eLBA)
			println("---------------------------------------------")
			if err != nil {
				panic(err)	
			}
			entry++
		}

	} 


}

func LBA(lba int64) int64 {
	return lba * SECTORSIZE
}

func DecodeHexString(buffer []byte) string {
	const hextable = "0123456789ABCDEF"
	
	dst := make([]byte, len(buffer)*2)
	j := 0
	for _, v := range buffer {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		j += 2
	}
	return string(dst)
}

func ToLittleEndian(str string, bytes int) string {
	result := ""
	j := 0
	len := bytes*2
	for i := 0; i < bytes; i++ {
		temp := str[len-j-2:len-j]
		result += temp
		j += 2
	}
	return result
}

func populateType(pType string) map[string]string {

	types := make(map[string]string)
	if(pType == "mbr") {
		file, err := os.Open("mbr_partition_types.csv")

		if err != nil {
	//		log.Fatalf("failed opening file: %s", err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			aType := scanner.Text()
			arguments := strings.Split(aType,",")
			key := arguments[0]
			types[key] = arguments[1]
		}

		file.Close()
		return types

	} else if(pType == "gpt") {
		file, err := os.Open("gpt_partition_guids.csv")

		if err != nil {
	//		log.Fatalf("failed opening file: %s", err)
		}
		types := make(map[string]string)
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			aType := scanner.Text()
			arguments := strings.Split(aType,",")
			key := arguments[0]

			types[key] = arguments[1]+" - "+arguments[2]
		}

		file.Close()
		return types
	}

	return types
}

