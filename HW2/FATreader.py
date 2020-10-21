from sys import *

filename = argv[1]
offset = 63488*512

with open(filename, "rb") as iso_file:

    # seek to the correct Partition
    iso_file.seek(offset ,0)
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

    s_sector_address_DS = numFATS * sectorsPERfat + reserved_area_in_sectors
    print("Starting Sector Address of the Data Section: {}".format(s_sector_address_DS))

    print()

    #seek to Root, print contents
    i = 0
    while i < 512:
        iso_file.seek(s_sector_address_DS * bytesPERsector + offset + i, 0)
        buffer = iso_file.read(32)
        print()
        print("32 byte chunk: ", buffer)
        print("Name: ", buffer[1:10:2], buffer[16:26:2])
        print("Where data starts: ", int.from_bytes(buffer[26:28], byteorder='little'))
        print("Size: {}".format(int.from_bytes(buffer[28:32], byteorder='little')))
        i+=32
        
    
    print()
    byte_address_DE = "____"
    print("Byte Address of Directory Entry: {}".format(byte_address_DE))

    start_cluster = int.from_bytes(buffer[122:123], byteorder='little')
    print("Start Cluster Address of File Data: {}".format(start_cluster))

    size_of_file = int.from_bytes(buffer[124:127], byteorder='little')
    print("Size of File in Bytes: {}".format(size_of_file))

    print("Ending Cluster Address of File: {}".format("___"))