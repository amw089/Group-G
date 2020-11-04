//////////////////////////////////////////////
// Anti-File-Hiding                       ////
// This is a stripped down version of hw2 ////
// New code in line 157                   ////
//////////////////////////////////////////////

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"encoding/hex"
)
// Global Variables //
var DEBUG bool = false
var EXTENSIONFILTER string = "JPG"
var NAMEFILTER string = ""

func main() {
    // Print Usage at the begining
	PrintUsage()
    // Parsing the command line for methods, -mbr -gpt -fat
	if len(os.Args) < 2 {
		os.Exit(0)
	} 
	file := os.Args[1]

    // Open File for reading
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

		buffer := make([]byte, 512)
		f.Read(buffer)

		bPERs,clusterRootDir,startSecAddRootDir := BootSectorInfo(buffer)
		// Skiping to root directory and skiping the folders
		f.Seek(((startSecAddRootDir+1)*512),0)
		currSectAddress := ((startSecAddRootDir+1)*512)
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
				println("\n-------------------------EOF------------------------")
				os.Exit(0)
			}
			lastEntryAddress += 32

            // Getting the name of the file and extension in the first slice of the buffer
			nameOfFile := fmt.Sprintf("%s", buffer[0:8])
			extension := fmt.Sprintf("%s", buffer[8:11])
			segmentTag := ""
            // Getting the attributes that gives specifics of the entry (i.e file, subdirectory, ...)
			attributes := buffer[11:12]
            // Getting the file cluster address. It is placed in 2 slices of the entry. They are in little endian.
			FirstPartOfCluster := buffer[20:22]
			str := DecodeHexString(FirstPartOfCluster) 	
			x := ToLittleEndian(str,2)
			SecondPartOfCluster := buffer[26:28]
			str = DecodeHexString(SecondPartOfCluster) 	
			x += ToLittleEndian(str,2)
			fileClusterAddress, err := strconv.ParseInt(x, 16, 64)
            // Getting the size of the file in the last slice of the entry
			chunk := buffer[28:32]
			str = DecodeHexString(chunk) 	
			x = ToLittleEndian(str,4)
			sizeOfFile, err := strconv.ParseInt(x, 16, 64)

            // This is done to filter the entries wanted for extracting.
			if DecodeHexString(attributes) != "20" {
				continue
			}
			if EXTENSIONFILTER != "" {
				if extension != EXTENSIONFILTER {
					continue
			} 				
			}
			if NAMEFILTER != "" {
				if nameOfFile[0:len(NAMEFILTER)] != NAMEFILTER {
					continue				
				}
			}
			
            // Create list of clusters for FAT chain, It is empty, but will be incremented dynamically
            // Add the first cluster found in the file entry. Keep track of the current cluster 
			var chain []int64 = make([]int64,0)
			chain = append(chain, fileClusterAddress)
			currentCluster := fileClusterAddress
            // Debugging, Printing the address in the FAT1 in bytes. 32*512=the start of the FAT1, currentcluster*4=the start of the chain, +offset=to get the address  
			if DEBUG {
				println("Address Jump to first cluster in chain of file: ",((32*bPERs)+(currentCluster*4)))
			}
            // Seek to start of the cluster chain 
            // We iterate throught the FAT section of each file to append to the chains list
            // If we reach an EOF or zeros, break
			for {
				f.Seek(((32*bPERs)+(currentCluster*4)),0)
				
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
				println("----------------------------------------------------")
				println("\nBuilding "+nameOfFile)
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
				println("----------------------------------------------------")
				println("Address Jump to first cluster in chain of file: ",(512*((fileClusterAddress-clusterRootDir)+startSecAddRootDir)))
			}
            // Where I look for metadata seeking each cluster of the file according to the chain list
			for index := range chain {
				nextSectAddress := (chain[index]-clusterRootDir)+startSecAddRootDir
				nextSectAddressInBytes := (nextSectAddress*512)
				f.Seek(nextSectAddressInBytes,0)
				buffer = make([]byte, 512)
				_, err = f.Read(buffer)
				if err != nil {
					break	
				}
				// Evaluate headers for problems and types
				if index == 0 {
					FileHeader := DecodeHexString(buffer[:16])
					if DEBUG {
						println("----------------------------------------------------")
						println("File Header",FileHeader)
					}
					
					if strings.Index(FileHeader,"383961") == 6 {
						segmentTag = "GIF"
						CorrectedHeader, err := hex.DecodeString("474946")
						if _, err = recoveredFile.Write(CorrectedHeader); err != nil {
							        // LOG				println(err)				
						}
						buffer = buffer[3:]	
					} else if strings.Index(FileHeader,"462D") == 6 {
						segmentTag = "PDF"
						CorrectedHeader, err := hex.DecodeString("255044")
						if _, err = recoveredFile.Write(CorrectedHeader); err != nil {
							    // LOG				println(err)				
						}
						buffer = buffer[3:]	
						
					} else if strings.Index(FileHeader,"470D") == 6 {
						segmentTag = "PNG"
						CorrectedHeader, err := hex.DecodeString("89504E")
						if _, err = recoveredFile.Write(CorrectedHeader); err != nil {
							    // LOG				println(err)				
						}
						buffer = buffer[3:]	
						
					} else if strings.Index(FileHeader,"E0A1") == 6 {
						segmentTag = "DOC"
						CorrectedHeader, err := hex.DecodeString("D0CF11")
						if _, err = recoveredFile.Write(CorrectedHeader); err != nil {
							    // LOG				println(err)				
						}
						buffer = buffer[3:]	
						
					} else if strings.Index(FileHeader,"041400") == 6 {
						segmentTag = "JAR"
						CorrectedHeader, err := hex.DecodeString("504B03")
						if _, err = recoveredFile.Write(CorrectedHeader); err != nil {
							    // LOG				println(err)				
						}
						buffer = buffer[3:]	
						
					} else if strings.Index(FileHeader,"4F4354595045") == 6 {
						segmentTag = "HTML"
						CorrectedHeader, err := hex.DecodeString("3C2144")
						if _, err = recoveredFile.Write(CorrectedHeader); err != nil {
							    // LOG				println(err)				
						}
						buffer = buffer[3:]	
						
					} else if strings.Contains(FileHeader,"4A46494600") {
						segmentTag = "JIF"
					} else if strings.Contains(FileHeader,"4578696600") {
						segmentTag = "EXIF"
					} else if strings.Contains(FileHeader,"535049464600") {
						segmentTag = "EXIF"
					}
				}	
                // Print to the file appending
				if _, err := recoveredFile.Write(buffer); err != nil {
                // LOG				println(err)
				}
			}
            ///	Printing the info of each file entry         ////////////
			println("----------------------------------------------------")
			
			fmt.Printf("Name: %s\n",nameOfFile)
			fmt.Printf("Extension: %s\n",extension)
			fmt.Printf("Type: %s\n",segmentTag)
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
	}

func BootSectorInfo(buffer []byte) (int64,int64, int64) {
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
	if err != nil {
		panic(err)	
	}

	return bPERs,clusterRootDir,startSecAddRootDir
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

func PrintUsage() {
	println("---------------------------------------------")
	println("Usage: go run hw2-1.go <filename>")
	println("---------------------------------------------")
}


