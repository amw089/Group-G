package main

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
	"strings"
)

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
	p.desType = strings.ToUpper(dhs[8:10])
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

	// Open File
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Get the content
	dhs, err := DecodeFileToHexString(f)
	if err != nil {
		panic(err)
	}

	typesMap := populateType()

	numPartitions := 0
	part1, status1 := newPartition(dhs[446*2 : 462*2])
	if status1 == 0 {
		numPartitions++
		part1.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
		printPartition(*part1,typesMap[part1.desType])
		fmt.Println("-----------------------------")
	}
	part2, status2 := newPartition(dhs[462*2 : 478*2])
	if status2 == 0 {
		numPartitions++
		part2.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
		printPartition(*part2,typesMap[part2.desType])
		fmt.Println("-----------------------------")
	}
	part3, status3 := newPartition(dhs[478*2 : 494*2])
	if status3 == 0 {
		numPartitions++
		part3.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
		printPartition(*part3,typesMap[part3.desType])
		fmt.Println("-----------------------------")
	}
	part4, status4 := newPartition(dhs[494*2 : 510*2])
	if status4 == 0 {
		numPartitions++
		part4.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
		printPartition(*part4,typesMap[part4.desType])
		fmt.Println("-----------------------------")
	}
	fmt.Println("Number of partitions: ", numPartitions)


}

func DecodeFileToHexString(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	const hextable = "0123456789abcdef"
	dst := make([]byte, len(buffer)*2)
	j := 0
	for _, v := range buffer {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		j += 2
	}
	
	return string(dst), nil
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

func populateType() map[string]string {

	file, err := os.Open("mbr_partition_types.csv")

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
		types[key] = arguments[1]
	}

	file.Close()

	return types
}

