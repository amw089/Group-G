from sys import *

try:
    filename = argv[2]
    mode = argv[1]
except:
    print("To run this program:\npython3 fileCarver.py -d 'iso.iso' (scans the data section)\npython3 fileCarver.py -f 'iso.iso' (scans the fat section)")
    exit()

DEBUG = True

OFFSET = 0
EXTENSION_FILTER = "JPG"
BEGINING_SIGNATURE = "ffd8ff"
END_SIGNATURE = "ffd9"

with open(filename, "rb") as iso_file:

    # seek to the correct Partition
    iso_file.seek(OFFSET ,0)
    # read the header
    buffer = iso_file.read(512)

    bytesPERsector = int.from_bytes(buffer[11:13], byteorder='little')
    print("Bytes/Sector: {}".format(bytesPERsector))

    sectorsPERcluster = int.from_bytes(buffer[13:14], byteorder='little')
    print("Sectors/Cluster: {}".format(sectorsPERcluster))

    reserved_area_in_sectors = int.from_bytes(buffer[14:16], byteorder='little')
    print("Size of Reserved Area in Sectors: {}".format(reserved_area_in_sectors))

    print("Start Address of 1st FAT: {}".format(reserved_area_in_sectors))

    numFATS = int.from_bytes(buffer[16:17], byteorder='little')
    print("# of FATs: {}".format(numFATS))

    sectorsPERfat = int.from_bytes(buffer[36:40], byteorder='little')
    print("Sectors/FAT: {}".format(sectorsPERfat))

    cluster_address_RD = int.from_bytes(buffer[44:48], byteorder='little')
    print("Cluster Address of Root Directory: {}".format(cluster_address_RD))

    s_sector_address_RD = numFATS * sectorsPERfat + reserved_area_in_sectors
    print("Starting Sector Address of the Data Section: {}".format(s_sector_address_RD))

################ SCAN THE DATA SECTION #################
    if mode == "-d":
        file_count = 0
        # skip to data section
        current_byte = s_sector_address_RD * bytesPERsector
        iso_file.seek(current_byte)
        length = len(iso_file.read())
        # search for file signatures
        while current_byte < length:
            if DEBUG == True and current_byte % 100000 == 0:
                print(current_byte, '/', length, '\t', round(current_byte/length * 100,2) ,"% searched")

            buffer = iso_file.read(3).hex()
            # file has been found
            if(buffer == BEGINING_SIGNATURE):
                file_count += 1
                if DEBUG == True:
                    print("FILE FOUND at: {}".format(current_byte))

                recovered_file = open("file{}_at_{}.JPG".format(file_count, current_byte), "wb+")
                recovered_file.write(b'\xFF\xD8\xFF')

                # read bytes after the starting signature
                current_byte += 3
                iso_file.seek(current_byte)
                current_byte += 1
                previous_buffer1 = "ff"
                buffer1 = iso_file.read(1)
                counter = 0
                # read and write 1 bytes at a time until EOF indicator
                while (1):
                    # writes data to file iteratively
                    recovered_file.write(buffer1)
                    iso_file.seek(current_byte)
                    buffer1 = iso_file.read(1)
                    current_byte += 1
                    if(previous_buffer1 + buffer1.hex() == END_SIGNATURE):
                        break
                    previous_buffer1 = buffer1.hex()

                recovered_file.write(b'\xFF\xD9')
                recovered_file.close()
            else:
                current_byte += 3
            if(current_byte % 16 == 0):
                iso_file.seek(current_byte)
            else:
                while current_byte % 16 != 0:
                    current_byte += 1
                iso_file.seek(current_byte)

############### SCAN THE FAT ######################
    elif mode == "-f":
        file_count = 1
        # location of FAT table, skip header
        current_byte = reserved_area_in_sectors * bytesPERsector + 8
        current_cluster = 1
        while current_byte < (reserved_area_in_sectors + sectorsPERfat) * bytesPERsector:
            cluster_chain = []
            # build cluster chain
            while (1):
                iso_file.seek(current_byte)
                buffer = iso_file.read(4)
                if(buffer.hex() != "00000000" and buffer.hex() != "ffffffff"):
                    current_cluster = int.from_bytes(buffer, byteorder='little')
                    cluster_chain.append(current_cluster)
                else:
                    current_byte += 16
                    break
                current_byte += current_cluster * 4
                
            print(cluster_chain)
            # create file to save the recovered data
            # iterate through cluser chain and jump to data location
            recovered_file = open("file{}-f.JPG".format(file_count), "wb+")
            for cluster in cluster_chain:
                next_section_address = ((cluster - cluster_address_RD) + s_sector_address_RD)
                next_section_address_bytes = next_section_address * bytesPERsector + OFFSET
                iso_file.seek(next_section_address_bytes, 0)
                buffer = iso_file.read(bytesPERsector)
                recovered_file.write(buffer)

            
    else:
        print("To run this program:\npython3 fileCarver.py -d 'iso.iso' (scans the data section)\npython3 fileCarver.py -f 'iso.iso' (scans the fat section)")