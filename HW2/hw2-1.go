package main

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
	"strings"
)

var SECTORSIZE int64 = 512

type partition struct {
	name        string
	boot        string
	startCHS    string
	desType     string
	endCHS      string
	startSector string
	size        string
}

func newPartition(dhs string) (*partition, int) {
	p := partition{}
	p.name = ""
	p.boot = dhs[0:2]
	p.startCHS = dhs[2:8]
	p.desType = dhs[8:10]
	p.endCHS = dhs[10:16]
	p.startSector = dhs[16:24]
	p.size = dhs[24:32]

	i := 0
	if p.size == "00000000" {
		i = -1
	}
	
	return &p, i
}

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

		// Get the content
		dhs, err := DecodeHexString(f,LBA(0),512)
		if err != nil {
			panic(err)
		}

		numPartitions := 0
		part1, status1 := newPartition(dhs[446*2 : 462*2])
		if status1 == 0 {
			numPartitions++
			part1.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
			printPartition(*part1,MBRMap[part1.desType])
			fmt.Println("-----------------------------")
		}
		part2, status2 := newPartition(dhs[462*2 : 478*2])
		if status2 == 0 {
			numPartitions++
			part2.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
			printPartition(*part2,MBRMap[part2.desType])
			fmt.Println("-----------------------------")
		}
		part3, status3 := newPartition(dhs[478*2 : 494*2])
		if status3 == 0 {
			numPartitions++
			part3.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
			printPartition(*part3,MBRMap[part3.desType])
			fmt.Println("-----------------------------")
		}
		part4, status4 := newPartition(dhs[494*2 : 510*2])
		if status4 == 0 {
			numPartitions++
			part4.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
			printPartition(*part4,MBRMap[part4.desType])
			fmt.Println("-----------------------------")
		}
		fmt.Println("Number of partitions: ", numPartitions)
	
	} else if mode == "-gpt" {
		
		GPTMap := populateType("gpt")
		entrySize := 128
		
		f.Seek(LBA(2),0)
		buffer := make([]byte, entrySize)
		
		entry := 1
		for entry < 5 {
			f.Read(buffer)

			fmt.Println("Partition ", entry)

			chunk := buffer[56:128]
			fmt.Printf("Name: %s\n",chunk)

			chunk = buffer[0:16]
			str := DecodeHexString2(chunk)
			GUID := ToLittleEndian(str[0:8],4)+"-"+ToLittleEndian(str[8:12],2)+"-"+ToLittleEndian(str[12:16],2)+"-"+str[16:20]+"-"+str[20:]
			println("GUID: "+GUID)
			println("Type: "+ GPTMap[GUID])

			chunk = buffer[32:40]
			str = DecodeHexString2(chunk)
			x := ToLittleEndian(str,8)
			sLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Starting LBA: %d\n", sLBA)

			chunk = buffer[40:48]
			str = DecodeHexString2(chunk)
			x = ToLittleEndian(str,8)
			eLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Ending LBA: %d\n", eLBA)
			println("---------------------------------------------")
			if err != nil {
				panic(err)	
			}
			entry++
		}
		fmt.Println("Number of partitions: ", entry-1)

	} 


}

func LBA(lba int64) int64 {
	return lba * SECTORSIZE
}

func DecodeHexString2(buffer []byte) string {
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

func DecodeHexString(out *os.File, lba int64, size int) (string, error) {

	// Only the first 512 bytes are used
	buffer := make([]byte, size)

	out.Seek(lba, 0)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	const hextable = "0123456789ABCDEF"
	
	dst := make([]byte, len(buffer)*2)
	j := 0
	for _, v := range buffer {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		j += 2
	}
	
	return string(dst), nil
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

func printPartition(part partition, pType string) {

	line := part.name + "\n" + "Boot: "
	if part.boot == "80" {
		line += "bootable"
	} else {
		line += "non-bootable"
	}

	
	line += "\n" + "Type: " + pType +"\n"+"LBA: "

	s1 := part.startSector[6:8]
	s2 := part.startSector[4:6]
	s3 := part.startSector[2:4]
	s4 := part.startSector[0:2]

	sSect, err2 := strconv.ParseInt(s1+s2+s3+s4, 16, 64)

	line += strconv.FormatInt(sSect, 10) + "\n" + "Size: "

	s1 = part.size[6:8]
	s2 = part.size[4:6]
	s3 = part.size[2:4]
	s4 = part.size[0:2]

	sSect, err2 = strconv.ParseInt(s1+s2+s3+s4, 16, 64)
	if err2 != nil {
		panic(err2)
	}

	line += strconv.FormatInt(sSect, 10)

	println(line)
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

