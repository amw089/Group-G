///////////////////////////////
///// CSC 443 - Forensics /////
/// MBR, GPT, FAT Analysis  ///
///////////////////////////////

package main

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
	"strings"
)
// Global Variables //
var DEBUG bool = false
var SECTORSIZE int64 = 512
var EXTENSIONFILTER string = "JPG"

func main() {
// Print Usage at the begining
	PrintUsage()
// Parsing the command line for methods, -mbr -gpt -fat
	if len(os.Args) > 2 {
		if mode := os.Args[1]; mode != "-1" {
			if mode != "-2" {
				os.Exit(0)
			}
		}
	} else {
		os.Exit(0)
	}
// Set the mode and file  
	mode := os.Args[1]
	file := os.Args[2]

// Open File for reading
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
// Reading std_in for an offset
offset := PrintOffsetUsage()
// Seek the offset of the wanted partition, create a 512 buffer to get the complete sector, and reading accordingly
// This is the start of Boot Sector. Here we analyze the info of the partition //
		f.Seek(int64(offset),0)
		buffer := make([]byte, 512)
		f.Read(buffer)
// Decoding the Bytes/Sector
		chunk := buffer[11:13]
		str := DecodeHexString(chunk) 	
		x := ToLittleEndian(str,2)
		bPERs, err := strconv.ParseInt(x, 16, 64)
		println("Bytes/Sector: ", bPERs)
// Decoding the Sectors/Cluster
		chunk = buffer[13:14]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,1)
		temp, err := strconv.ParseInt(x, 16, 64)
		println("Sectors/Cluster: ", temp)
// Decoding the size of the reseverd area, and printing in bytes
// This gives us the start of the 1st FAT
		chunk = buffer[14:16]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,2)
		sizeReserved, err := strconv.ParseInt(x, 16, 64)
		println("Size of Reserved Area in Sectors: ", sizeReserved)
		println("Start Address of 1st FAT: ", sizeReserved)
// Decoding the number of FATs
		chunk = buffer[16:17]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,1)
		numFATs, err := strconv.ParseInt(x, 16, 64)
		println("# of FATs: ", numFATs)		
// Decoding the Sectors per FAT
		chunk = buffer[36:40]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,2)
		secPerFAT, err := strconv.ParseInt(x, 16, 64)
		println("Sectors/FAT: ", secPerFAT)
// Decoding the Cluster Address of the root Directory
		chunk = buffer[44:48]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,4)
		clusterRootDir, err := strconv.ParseInt(x, 16, 64)
		println("Cluster Address of Root Directory: ", clusterRootDir)
// Decoding the Starting Sector address of the Data Section. Reserved Size + the multiplication of the sectors per FAT and the Number of FATs
		startSecAddRootDir := sizeReserved + (secPerFAT*numFATs)
		println("Starting Sector Address of the Data Section: ", startSecAddRootDir)
// Skiping to root directory and skiping the folders
		f.Seek(((startSecAddRootDir+1)*512)+int64(offset),0)
		currSectAddress := ((startSecAddRootDir+1)*512)+int64(offset)
		if DEBUG {
			println("Address For the file info in root directory: ",currSectAddress)
		}
	if mode == "-1" {
		///////////////
	} else if mode == "-2" { 
		
//	
	BeginSignatureBuffer := make([]byte, 3)
	EndSignatureBuffer := make([]byte, 2)
	previousAddress := startSecAddRootDir*512

		for {
			// Start of data section
			f.Seek(previousAddress,0)
			nameOfFile := "file" + strconv.Itoa(int(previousAddress))+".jpg"
			_, err := f.Read(BeginSignatureBuffer)
			if err != nil {
				println("------------END OF ISO----------")
	//			panic(err)	
			}
			previousAddress += 3
			BeginSignature := DecodeHexString(BeginSignatureBuffer) 	

			if BeginSignature == "FFD8FF" {
				recoveredFile, err := os.OpenFile(nameOfFile,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					println("------------CANT OPEN FILE---------")
	//				panic(err)	
				}
				defer f.Close()
//////
				if _, err := recoveredFile.Write(BeginSignatureBuffer); err != nil {
					println("------------CANT WRITE TO FILE---------")
//					println(err)	
				}

				for {
					_, err := f.Read(EndSignatureBuffer)
					if err != nil { 
						panic(err)	
					}
					if _, err := recoveredFile.Write(EndSignatureBuffer); err != nil {
						println(err)
					}
					EndSignature := DecodeHexString(EndSignatureBuffer)
					if EndSignature == "FFD9" {
						break
					}

				}
				
			} 
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

func PrintUsage() {
	println("---------------------------------------------")
	println("Usage: go run hw2-1.go [MODES] <filename>")
	println("MODES:\n -mbr     mbr analysis mode\n -gpt     gpt analysis mode\n -fat     fat analysis mode")
	println("---------------------------------------------")
}

func PrintOffsetUsage() int {
	println("----------------------------------------------------")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a byte Offset if needed to jump to a specific Partition, if not, leave blank.\nTo find the partitions run -mbr or -gpt\nOffset: ")
	text, _ := reader.ReadString('\n')	
	offset, _ := strconv.Atoi(strings.Replace(text, "\n","",1))
	fmt.Println("Offset of Partitions: ",offset)
	println("----------------------------------------------------")
	return offset
}
