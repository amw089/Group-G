# Anna Wolf and Katie Hay
# File System Analysis Homework
# To run: FAT.py <iso file>          ex: FAT.py FAT_FS.iso

import sys, struct, math

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

# Searches for "/Photos/homework.jpg"
def findfile(file):
    with open(file, "rb") as f:
        bootdata = f.read()
        offset = (FAT_info.get("ssad",0) * 512) #548864
        dd = bootdata[offset+32:offset+40]
        #print("offset = ",offset) #end/sep
        found = "false"
        while found == "false":
            # If string photos is found in current 8 bytes, look for highbytes and lowbytes to find cluster address
            if "photos" in str(dd).lower():
                found = "true"
                highbytes = int.from_bytes(bootdata[offset+32+20:offset+32+22], "little")
                lowbytes = int.from_bytes(bootdata[offset+32+26:offset+32+28], "little")
                print("Cluster Address of Directory Entry: ", str(lowbytes + highbytes))
                
                # Photos folder found, search for homework file
                onset = (FAT_info.get("ssad",0) * 512 + 512)
                found2 = "false"
                ee = bootdata[onset + 96: onset + 104]
                while found2 == "false":
                    if "homework" in str(ee).lower():
                        found2 = "true"
                        highbytes2 = int.from_bytes(bootdata[onset+96+20:onset+96+22], "little")
                        lowbytes2 = int.from_bytes(bootdata[onset+96+26:onset+96+28], "little")
                        print("Cluster Address of File Data: ", str(lowbytes2 + highbytes2))
                        sizeoffile = int.from_bytes(bootdata[onset+96+28:onset+96+32], "little")
                        print("Size of File in Bytes: ", str(sizeoffile))
                        # Gets end-cluster based off sizeoffile, plus 4 since starts at cluster 4
                        filecluster = math.ceil(sizeoffile/512) + 4 
                        print("Ending Cluster Address of File: ", str(filecluster))
                    else:
                        onset += 32
            else:
                offset += 32
                
#to get to entry in FAT: cluster address * 4 (offset from beginning of FAT)
                #ex: 32*512 offset --> beginning of fat
                #4 * 4 = 16
                
        
# Implementation
opensesame(sys.argv[1])
findfile(sys.argv[1])
