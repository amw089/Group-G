#Anna Wolf and Katie Hay
#File Carving Homework
#To run: file_carving.py <iso file>         ex: file_carving.py FAT_Corrupted.iso

import sys, struct, math
from sys import *

FAT_info = {}
FAT_desc = {}

# Reads boot sector and gets needed info
def opensesame(file):
    with open(file, "rb") as f:
        
        #reads all 512 bytes of boot sector
        bootdata = f.read(512)
        
        #bytes per second is bytes 11-12, use 13 because it doesn't count the last number of the list.
        #the little will reverse the two bytes since it's in little endian.
        FAT_info["bps"] = int.from_bytes(bootdata[11:13], "little")
        FAT_desc["bps"] = "Bytes/Sector"
        
        FAT_info["spc"] = int.from_bytes(bootdata[13:14], "little")
        FAT_desc["spc"] = "Sectors/Cluster"
        
        FAT_info["sar"] = int.from_bytes(bootdata[14:16], "little")
        FAT_desc["sar"] = "Size of Reserved Area in Sectors"
        
        #FAT comes right after boot sector, so starting address is the offset
        FAT_info["saf"] = FAT_info.get("sar",0)
        FAT_desc["saf"] = "Start Address of 1st FAT"
        
        FAT_info["nof"] = int.from_bytes(bootdata[16:17], "little")
        FAT_desc["nof"] = "# of FATs"
      
        FAT_info["spf"] = int.from_bytes(bootdata[36:40], "little")
        FAT_desc["spf"] = "Sectors/FAT"
        
        FAT_info["card"] = int.from_bytes(bootdata[44:48], "little")
        FAT_desc["card"] = "Cluster Address of Root Directory"
        
        #size of reserved area + (# of FATs * sectors/FAT) = Starting address for Data
        FAT_info["ssad"] = (FAT_info.get("sar",0) + (FAT_info.get("spf",0) * FAT_info.get("nof",0)))
        FAT_desc["ssad"] = "Starting Sector Address of the Data Section"
 
        for i in FAT_info:
            print (FAT_desc[i], ':', FAT_info[i])
            

def scandata(file):
    with open(file, "rb") as f:
        #file will start with FFD8FF and end with FFD9
        find_first = "FFD8FF"
        find_last = "FFD9"
        #read everything
        data = f.read()
        #offset given from previous code, go to it to get to data section
        offset = (FAT_info.get("ssad",0) * FAT_info.get("bps", 0))
        f.seek(offset)
        num_file = 1
        #keep track of where you are by comparing offset to length of FAT
        length_file = len(f.read())
        while offset < length_file:
            #read the next 3 bytes
            buffer = f.read(3)
            #compare bytes to the FFD8FF
            if (buffer.hex().upper() == find_first):
                #skip 3 bytes (the FF D8 DD)
                offset += 3
                f.seek(offset)
                #read the next two bytes
                buffer1 = f.read(2)
                #format how to open
                recovered_file = open("file{}.JPG".format(num_file), "wb+")
                recovered_file.write(b'\xFF\xD8\xFF')
                #compare offset to FFD9
                while buffer1.hex().upper() != find_last:
                    #until we reach FFD9, keep writing into the current file
                    recovered_file.write(buffer1)
                    #check every 2 bytes so we don't miss the FF D9
                    offset += 2
                    f.seek(offset)
                    buffer1 = f.read(2)
                #once we find FF D9, we can finish writing to the current file
                recovered_file.write(b'\xFF\xD9')
                recovered_file.close()
                #now we move onto the next file
                num_file += 1
            #if we don't find FF D8 DD, check the next 3 bytes
            else:
                offset += 3

            f.seek(offset)
            
# Implementation    
scandata(sys.argv[1])


