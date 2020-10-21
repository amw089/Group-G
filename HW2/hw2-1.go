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
var NAMEFILTER string = "NORWAY"

func main() {
// Print Usage at the begining
	PrintUsage()
// Parsing the command line for methods, -mbr -gpt -fat
	if len(os.Args) > 2 {
		if mode := os.Args[1]; mode != "-mbr" {
			if mode != "-gpt" {
				if mode != "-fat" {
					os.Exit(0)
				}
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
// Start of MBR analysis //
// Populate the MBR type Map for storing all of the MBR types for identification
// Set the entry size to 16
	if mode == "-mbr" {
		MBRMap := populateType("mbr")
		entrySize := 16
// Create the buffer size according to the entry size
// Seek the start of the partition info area
		buffer := make([]byte, entrySize)
		f.Seek(446,0)
// Itiration to check for the number of available partitions
// We only have 4 in an MBR
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
// Printing the Number of partitions available 
		fmt.Println("Number of partitions: ", numPartions)
		println("---------------------------------------------")
// Seeking to the start of the partition area again to traverse through the partitions 
// an get the needed info		
		f.Seek(446,0)
		entry = 1
		for entry < numPartions+1 {
// Read the first chunk of info with the size buffer 16, and printing a generic name for the partition
			f.Read(buffer)
			fmt.Println("Partition ", entry)
// Checking slices of the 16 byte chunk read //
// First is the bootable or not, 80 is bootable
			chunk := buffer[0:1]
			if DecodeHexString(chunk) == "80" {
				println("Boot: bootable")
			} else {
				println("Boot: non-bootable")
			}
// Reading the type in hex, looking for the type name in the Map of types, and printing it
			chunk = buffer[4:5]
			str := DecodeHexString(chunk)
			println("Type: "+ MBRMap[str])
// Reading for the LBA address, the hex is decoded and converted to big endian for printing as an integer 
			chunk = buffer[8:12]
			str = DecodeHexString(chunk)
			x := ToLittleEndian(str,4)
			sLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Address LBA: %d\n", sLBA)
// Reading and printing the Sectors in Partition
			chunk = buffer[12:16]
			str = DecodeHexString(chunk)
			x = ToLittleEndian(str,4)
			eLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Sectors in Partition: %d\n", eLBA)
			println("---------------------------------------------")
// Catching any errors 			
			if err != nil {
				panic(err)	
			}
// Increment the counter, and looping back
			entry++
		}
// Start of GPT analysis //
// Populate the GPT type Map for storing all of the GPT types for identification
// Set the entry size to 128
	} else if mode == "-gpt" {
		GPTMap := populateType("gpt")
		entrySize := 128
// Seek to LBA 2 to read the partitions, and create buffer acording to entry size 128
		f.Seek(LBA(2),0)
		buffer := make([]byte, entrySize)
// Itiration to check for the number of available partitions
// We only have 128 in a GPT 
		entry := 1
		numPartions := 0
		for entry <= 128 {
			_, err := f.Read(buffer)
			if err != nil {
				break	
			} 
// If the first 16 bytes are not zeros, count the entry as having a partition 
			if DecodeHexString(buffer[0:16]) != "00000000000000000000000000000000" {
				numPartions = entry
			} 
			entry++
		}
// Seeking to the start of the partition area again to traverse through the partitions 
// an get the needed info
		f.Seek(LBA(2),0)
		
		fmt.Println("Number of partitions: ", numPartions)
		println("---------------------------------------------")
// Traversing through all of the available partitions
		entry = 1
		for entry < numPartions+1 {
			f.Read(buffer)
// Print the partiton number and name. GPT has names
			fmt.Println("Partition ", entry)
			chunk := buffer[56:128]
			fmt.Printf("Name: %s\n",chunk)
// Slice the first 16 bytes to get the GUID type. It is decoded by converting the first 3 chunks to little endian and the last 2 stay in big endian
// Printing the GUID type, looking for the type in the GPT type Map, and printing it 
			chunk = buffer[0:16]
			str := DecodeHexString(chunk)
			GUID := ToLittleEndian(str[0:8],4)+"-"+ToLittleEndian(str[8:12],2)+"-"+ToLittleEndian(str[12:16],2)+"-"+str[16:20]+"-"+str[20:]
			println("GUID: "+GUID)
			println("Type: "+ GPTMap[GUID])
// Decoding the starting LBA 
			chunk = buffer[32:40]
			str = DecodeHexString(chunk)
			x := ToLittleEndian(str,8)
			sLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Starting LBA: %d\n", sLBA)
// Decoding the ending LBA
			chunk = buffer[40:48]
			str = DecodeHexString(chunk)
			x = ToLittleEndian(str,8)
			eLBA, err := strconv.ParseInt(x, 16, 64)
			fmt.Printf("Ending LBA: %d\n", eLBA)
			println("---------------------------------------------")
// Error checking, increment counter, and loop back to get the next partition
			if err != nil {
				panic(err)	
			}
			entry++
		}
// Start of FAT analysis //
// Set the entry size to 128
	} else if mode == "-fat" { 
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
// Setting the buffer to read 32 byte entries, and reading the first chunk
		buffer = make([]byte, 32)
		f.Read(buffer)
// Recording the last entry address to jump back after getting the file info and actual file. 
// Every time we read, we add the number of bytes read
		lastEntryAddress := currSectAddress + 32
// Itirate until the end of the file to look for all posible files with the filter provided in the global variables
		for {
// Make a 32 byte buffer to read each entry. It has to be 32 bytes, that is why it's called FAT32 
			buffer = make([]byte, 32)
			_, err := f.Read(buffer)
			if err != nil {
				panic(err)	
			}
			lastEntryAddress += 32
// Debugging, printing the entier 32 bytes
			if DEBUG {
				println("----------------------------------------------------")
				fmt.Printf("32 byte chunk: %x\n", buffer)
				println("----------------------------------------------------")
			}
// Getting the name of the file and extension in the first slice of the buffer
			name := fmt.Sprintf("%s", buffer[0:8])
			extension := fmt.Sprintf("%s", buffer[8:11])
			nameOfFile := name+"."+extension
// Getting the attributes that gives specifics of the entry (i.e file, subdirectory, ...)
			attributes := buffer[11:12]
// Getting the file cluster address. It is placed in 2 slices of the entry. They are in little endian.
			FirstPartOfCluster := buffer[20:22]
			str = DecodeHexString(FirstPartOfCluster) 	
			x = ToLittleEndian(str,2)
			SecondPartOfCluster := buffer[26:28]
			str = DecodeHexString(SecondPartOfCluster) 	
			x += ToLittleEndian(str,2)
			fileClusterAddress, err := strconv.ParseInt(x, 16, 64)
// Getting the size of the file in the last slice of the entry
			chunk = buffer[28:32]
			str = DecodeHexString(chunk) 	
			x = ToLittleEndian(str,4)
			sizeOfFile, err := strconv.ParseInt(x, 16, 64)
// This is done to filter the entries wanted for extracting.
			if extension != EXTENSIONFILTER {
				continue				
			}
			if nameOfFile[0:len(NAMEFILTER)] != NAMEFILTER {
				continue				
			}
// Create list of clusters for FAT chain, It is empty, but will be incremented dynamically
// Add the first cluster found in the file entry. Keep track of the current cluster 
			var chain []int64 = make([]int64,0)
			chain = append(chain, fileClusterAddress)
			currentCluster := fileClusterAddress
// Debugging, Printing the address in the FAT1 in bytes. 32*512=the start of the FAT1, currentcluster*4=the start of the chain, +offset=to get the address  
			if DEBUG {
				println("Address Jump to first cluster in chain of file: ",((32*bPERs)+(currentCluster*4))+int64(offset))
			}
// Seek to start of the cluster chain 
// We iterate throught the FAT section of each file to append to the chains list
// If we reach an EOF or zeros, break
			for {
				f.Seek(((32*bPERs)+(currentCluster*4))+int64(offset),0)
				
				buffer = make([]byte, 4)
				f.Read(buffer)
				str = DecodeHexString(buffer)
				x = ToLittleEndian(str,4)
				currentCluster, err = strconv.ParseInt(x, 16, 64)
				if err != nil {
					panic(err)	
				}
				chain = append(chain, currentCluster)
				if x == "0FFFFFFF" {
					break
				}
				if x == "00000000" {
					break
				}
				
			}
			if DEBUG {
				fmt.Println("Chain array :", chain)
			}
// Create a file to start appending the metadata of each file in the corresponding cluster		
			recoveredFile, err1 := os.OpenFile(nameOfFile,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err1 != nil {
	//			println(err1)
			}
			defer f.Close()
// Some debugging
			if DEBUG {
				println("Address Jump to first cluster in chain of file: ",(512*((fileClusterAddress-clusterRootDir)+startSecAddRootDir))+int64(offset))
			}
// Where I look for metadata seeking each cluster of the file according to the chain list
			for index := range chain {
				nextSectAddress := (chain[index]-clusterRootDir)+startSecAddRootDir
				nextSectAddressInBytes := (nextSectAddress*512)+int64(offset)
				f.Seek(nextSectAddressInBytes,0)
				buffer = make([]byte, 512)
				_, err = f.Read(buffer)
				if err != nil {
					break	
				}	
// Print to the file appending
				if _, err := recoveredFile.Write(buffer); err != nil {
// LOG				println(err)
				}

			}
///	Printing the info of each file entry ////////////
			println("----------------------------------------------------")
			fmt.Printf("Name: %s\n", nameOfFile)
			fmt.Printf("Attributes: %x\n", attributes)
			fmt.Printf("Start of File Cluster: %d\n", fileClusterAddress)
			println("Size: ", sizeOfFile)
			println("Ending Cluster Address of File: ", chain[len(chain)-2]+1)
// Go back to the previous address for the next file
			lastEntryAddress, err = f.Seek(lastEntryAddress,0)
			if err != nil {
				panic(err)	
			}
		}


		//check error
		if err != nil {
			panic(err)	
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
