from sys import *

filename = argv[1]
OFFSET = 63488*512
END = 131071*512

EXTENSION_FILTER = "JPG"
NAME_FILTER = "NORWAY"

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



    # seek to Root, skip folders, read everything after that point
    current_address = (s_sector_address_RD + 1) * bytesPERsector + OFFSET
    iso_file.seek(current_address, 0)
    while current_address < END:
        buffer = iso_file.read(32)
        current_address = current_address + 32

        name = ""
        for byte in buffer[0:8]:
            name += chr(byte)
        extension = ""
        for byte in buffer[8:11]:   
            extension += chr(byte)
        
        name_of_file = name.replace(" ","") + "." + extension

        first_part_cluster = int.from_bytes(buffer[20:22], byteorder='little')
        second_part_cluster = int.from_bytes(buffer[26:28], byteorder='little')

        file_cluster_address = first_part_cluster + second_part_cluster

        size_of_file = int.from_bytes(buffer[28:32], byteorder='little')

        # checks if we are looking at the correct file
        if extension != EXTENSION_FILTER or name_of_file[0:len(NAME_FILTER)] != NAME_FILTER:
            continue	

        # seek to start of cluster chain, then iterate through FAT table until EOF
        # saves cluster chain for file from FAT table
        cluster_chain = []
        cluster_chain.append(file_cluster_address)
        current_cluster = file_cluster_address
        while True:
            iso_file.seek((reserved_area_in_sectors * bytesPERsector) + (current_cluster * 4) + OFFSET)
            buffer = iso_file.read(4)
            current_cluster = int.from_bytes(buffer, byteorder='little')
            cluster_chain.append(current_cluster)
            # checks for EOF/NULL bytes
            temp = 0
            for byte in buffer:
                temp += byte
            if temp == 0 or temp == 16777215:
                break

        # create file to save the recovered data
        recovered_file = open(name_of_file, "wb+")
        for cluster in cluster_chain:
            next_section_address = ((cluster - cluster_address_RD) + s_sector_address_RD)
            next_section_address_bytes = next_section_address * bytesPERsector + OFFSET
            iso_file.seek(next_section_address_bytes, 0)
            buffer = iso_file.read(bytesPERsector)
            recovered_file.write(buffer)
        
        # print out data on file
        print()
        print("Name:", name_of_file)
        print("Start of File Cluster: ", file_cluster_address)
        print("Size: ", size_of_file)
        print("Ending Address of File: ", cluster_chain[len(cluster_chain)-3]+1)

        iso_file.seek(current_address)