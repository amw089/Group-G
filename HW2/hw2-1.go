package main

import (
	"fmt"
	"os"
	"strconv"
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

	numPartitions := 0
	part1, status1 := newPartition(dhs[446*2 : 462*2])
	if status1 == 0 {
		numPartitions++
		part1.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
		printPartition(*part1)
		fmt.Println("-----------------------------")
	}
	part2, status2 := newPartition(dhs[462*2 : 478*2])
	if status2 == 0 {
		numPartitions++
		part2.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
		printPartition(*part2)
		fmt.Println("-----------------------------")
	}
	part3, status3 := newPartition(dhs[478*2 : 494*2])
	if status3 == 0 {
		numPartitions++
		part3.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
		printPartition(*part3)
		fmt.Println("-----------------------------")
	}
	part4, status4 := newPartition(dhs[494*2 : 510*2])
	if status4 == 0 {
		numPartitions++
		part4.name = fmt.Sprintf("%s%d", "Partition", numPartitions)
		printPartition(*part4)
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

func printPartition(part partition) {

	line := part.name + "\n" + "Boot: "
	if part.boot == "80" {
		line += "active"
	} else {
		line += "inactive"
	}

	line += "\n" + "Starting CHS: "
	bsC := part.startCHS[4:6]
	bsH := part.startCHS[2:4]
	bsS := part.startCHS[0:2]

	C, err := strconv.ParseInt(bsC, 16, 64)
	H, err := strconv.ParseInt(bsH, 16, 64)
	S, err := strconv.ParseInt(bsS, 16, 64)

	CHS := "(" + strconv.FormatInt(C, 10) + "," + strconv.FormatInt(H, 10) + "," + strconv.FormatInt(S, 10) + ")"
	line += CHS + "\n" + "Type: "
	// type
	line += "\n" + "Ending CHS: "

	beC := part.endCHS[4:6]
	beH := part.endCHS[2:4]
	beS := part.endCHS[0:2]

	C, err = strconv.ParseInt(beC, 16, 64)
	H, err = strconv.ParseInt(beH, 16, 64)
	S, err = strconv.ParseInt(beS, 16, 64)

	if err != nil {
		panic(err)
	}

	CHS = "(" + strconv.FormatInt(C, 10) + "," + strconv.FormatInt(H, 10) + "," + strconv.FormatInt(S, 10) + ")"
	line += CHS + "\n" + "Starting Sector: "

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
