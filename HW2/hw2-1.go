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
				if mode != "-fat" {
					println("Usage: go run hw2-1.go [MODES] <filename>")
					println("MODES:\n -mbr     mbr analysis mode\n -gpt     gpt analysis mode\n -fat     fat analysis mode")
					os.Exit(0)
				}
			}
		}
	} else {
		println("Usage: go run hw2-1.go [MODES] <filename>")
		println("MODES:\n -mbr     mbr analysis mode\n -gpt     gpt analysis mode\n -fat     fat analysis mode")
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

	} else if mode == "-fat" { 

		// For Entering an offset
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter and Offset if needed, if not, leave blank: ")
		text, _ := reader.ReadString('\n')	
		offset, _ := strconv.Atoi(strings.Replace(text, "\n","",1))
		fmt.Println("Offset of Partitions: ",offset)

		//Seek the offset
		f.Seek(int64(offset),0)
		buffer := make([]byte, 512)
		f.Read(buffer)

		chunk := buffer[11:13]
		str := DecodeHexString(chunk) 	
		x := ToLittleEndian(str,2)
		bPERs, err := strconv.ParseInt(x, 16, 64)
		println("Bytes/Sector: ", bPERs)

		chunk = buffer[13:14]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,1)
		temp, err := strconv.ParseInt(x, 16, 64)
		println("Sectors/Cluster: ", temp)
		
		chunk = buffer[14:16]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,2)
		sizeReserved, err := strconv.ParseInt(x, 16, 64)
		println("Size of Reserved Area in Sectors: ", sizeReserved)
		println("Start Address of 1st FAT: ", sizeReserved)

		chunk = buffer[16:17]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,1)
		numFATs, err := strconv.ParseInt(x, 16, 64)
		println("# of FATs: ", numFATs)		

		chunk = buffer[36:40]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,2)
		secPerFAT, err := strconv.ParseInt(x, 16, 64)
		println("Sectors/FAT: ", secPerFAT)

		chunk = buffer[44:48]
		str = DecodeHexString(chunk) 	
		x = ToLittleEndian(str,4)
		clusterRootDir, err := strconv.ParseInt(x, 16, 64)
		println("Cluster Address of Root Directory: ", clusterRootDir)

		startSecAddRootDir := sizeReserved + (secPerFAT*numFATs)
		println("Starting Sector Address of the Data Section: ", startSecAddRootDir)

		// Skiping to root directory and skiping the folders
		f.Seek(((startSecAddRootDir+1)*512)+int64(offset),0)
		currSectAddress := ((startSecAddRootDir+1)*512)+int64(offset)
		
		buffer = make([]byte, 32)
		f.Read(buffer)
		previousAddress := currSectAddress + 32
		
		i := 0
		for i < 512 {
			
			buffer = make([]byte, 32)
			_, err := f.Read(buffer)
			if err != nil {
				panic(err)	
			}
			previousAddress += 32
/*			println("----------------------------------------------------")
			fmt.Printf("32 byte chunk: %x\n", buffer)
			println("----------------------------------------------------")
*/			
			chunk = buffer[0:8]
			nameOfFile := fmt.Sprintf("%s", chunk)
			

			attributes := buffer[11:12]

			FirstPartOfCluster := buffer[20:22]
			str = DecodeHexString(FirstPartOfCluster) 	
			x = ToLittleEndian(str,2)
			SecondPartOfCluster := buffer[26:28]
			str = DecodeHexString(SecondPartOfCluster) 	
			x += ToLittleEndian(str,2)
			fileClusterAddress, err := strconv.ParseInt(x, 16, 64)
			

			chunk = buffer[28:32]
			str = DecodeHexString(chunk) 	
			x = ToLittleEndian(str,4)
			sizeOfFile, err := strconv.ParseInt(x, 16, 64)

			if sizeOfFile == 0 {
				continue				
			}
			if DecodeHexString(attributes)  == "00" {
				continue
			}

			//Create list for FAT chain, empty wil be incremented dynamically
			var chain []int64 = make([]int64,0)
			chain = append(chain, fileClusterAddress)
			currentCluster := fileClusterAddress

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
	//		fmt.Println("Chain array :", chain)
			recoveredFile, err1 := os.OpenFile(nameOfFile,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err1 != nil {
	//			println(err1)
			}
			defer f.Close()
			//Where I look for metadata
			for index := range chain {
			//	println("Current cluster: ",chain[index])
				
				nextSectAddress := (chain[index]-clusterRootDir)+startSecAddRootDir
				nextSectAddressInBytes := nextSectAddress*512
				f.Seek(nextSectAddressInBytes,0)
				buffer = make([]byte, 512)
				_, err = f.Read(buffer)
				if err != nil {
					break	
				}	

				if _, err := recoveredFile.Write(buffer); err != nil {
			//		println(err)
				}

			}

			//////////////
			println("----------------------------------------------------")
			fmt.Printf("Name: %s\n", nameOfFile)
			fmt.Printf("Attributes: %x\n", attributes)
			fmt.Printf("Start of File Cluster: %d\n", fileClusterAddress)
			println("Size: ", sizeOfFile)
			println("Ending Cluster Address of File: ", chain[len(chain)-2]+1)
			println("----------------------------------------------------")
			// Go back
			previousAddress, err = f.Seek(previousAddress,0)
			if err != nil {
				panic(err)	
			}

			i += 32
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