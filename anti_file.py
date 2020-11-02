#Anna Wolf and Katie Hay
#Anti-File-Hiding Homework
#To run: anti_file.py <iso file>         ex: anti_file.py FAT_Imposter.iso

# Libraries
import sys, struct, math
from sys import *

# Info
NEEDED = "JPG"
FILTERS = ["JPG exif", "JPEG jfif", "PNG", "GIF", "DOC", "DOCX", "HTML", "ODT", "PDF"]
BEGINNING = [b'\xff\xd8\xff\xe1', b'\xff\xd8\xff\xe0\x00', b'\x89\x50\x4e', b'\x47\x49\46', b'\xD0\xCF\x11', b'\x50\x4b\x03', b'\x3c\x21\x44', b'\x50\x4B\x03', b'\x25\x50\x44']
END = [b'\xff\xd9\x00\x00', b'\xff\xd9\x00\x00', b'\xff\xd9\x00\x00', b'\xff\xd9\x00\x00',  b'\xff\xd9\x00\x00',  b'\xff\xd9\x00\x00', b'\x74\x6d\x6c\x3e', b'\xff\xd9\x00\x00', b'\xff\xd9\x00\x00']
FAT_info = {}
FAT_desc = {}
OFFSET = 0

# Reads boot sector and gets needed info
def openSesame(file):
    with open(file, "rb") as f:
        #seeking to partition
        f.seek(OFFSET, 0)
        
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

        #find the original among the imposters
        findOriginal(file)
            

def findOriginal(file):
    with open(file, "rb") as f:
        #first seek to root
        length_file = len(f.read())
        current_add = (FAT_info.get("ssad",0) * FAT_info.get("bps", 0) + OFFSET)
        f.seek(current_add, 0)
        while current_add < length_file:
            data = f.read(32)
            current_add += 32

            #get file information
            n = ""
            ex = ""
            for byte in data[0:8]:
                n += chr(byte)
            for byte in data[8:11]:
                ex += chr(byte)
            file_name = n.replace(" ", "") + "." + ex.replace(" ","")

            #grab the file cluster
            first_part = int.from_bytes(data[20:22], "little")
            second_part = int.from_bytes(data[26:28], "little")
            file_cluster = first_part + second_part
            size_of_file = int.from_bytes(data[28:32], "little")

            #check the file ext
            if ex != NEEDED:
                continue
            
            #seek to start of cluster chain, then iterate through FAT table until EOF
            cluster_chain = []
            cluster_chain.append(file_cluster)
            current_cluster = file_cluster
            while True:
                #seek to FAT
                f.seek((FAT_info.get("sar",0) * FAT_info.get("bps",0)) + (current_cluster * 4) + OFFSET)
                data = f.read(4)
                current_cluster = int.from_bytes(data, "little")
                cluster_chain.append(current_cluster)
                #checking for EOF
                check = data.hex()
                if check == "0fffffff" or check == "ffffffff" or check == "00000000" or check == "ffffff0f":
                    break

            #get the beginning signature of the file
            next_section_address = ((cluster_chain[0] - FAT_info.get("card",0)) + FAT_info.get("ssad",0))
            next_section_address_bytes = next_section_address * FAT_info.get("bps",0) + OFFSET
            f.seek(next_section_address_bytes, 0)
            data = f.read(FAT_info.get("bps",0))
            starting_bytes = data[0:5]


            #see if it's the file we need by comparing its signature to our list of beginning
            #set a flag if it matches
            flag = 0
            if starting_bytes == BEGINNING[0]:
                filetype = "EXIF"
                flag = 1
            elif starting_bytes == BEGINNING[1]:
                filetype = "JFIF"
                flag = 1

            #if flag was set to 1, we found the original, display its data
            if flag == 1:
                # create file to save the recovered data
                recovered_file = open(file_name, "wb+")
                # iterate through cluster chain
                for cluster in cluster_chain:
                    next_section_address = (cluster - FAT_info.get("card",0) + FAT_info.get("ssad",0))
                    next_section_address_bytes = next_section_address * FAT_info.get("bps",0) + OFFSET
                    f.seek(next_section_address_bytes, 0)
                    data = f.read(FAT_info.get("bps",0))
                    recovered_file.write(data)

                # print out data on file
                print()
                print("Name:", file_name)
                print("EXIF or JFIF: ", filetype)
                print("Start of File Cluster: ", file_cluster)
                print("Size: ", size_of_file)
                print("Ending Cluster of File: ", cluster_chain[len(cluster_chain)-2])

            f.seek(current_add)

# Implementation    
openSesame(sys.argv[1])

# Note:
# Whenever we got stuck, we were able to look to our teammates codes for clarification.
# We understand the code and commented accordingly, but there were ideas that we had to
# bounce off of others code.
