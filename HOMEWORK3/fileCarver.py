from sys import *

filename = argv[1]
OFFSET = 0

BEGINING_SIGNATURE = "FFD8FF"
END_SIGNATURE = "FFD9"

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

    # skip to data section
    current_byte = s_sector_address_RD * bytesPERsector
    iso_file.seek(current_byte)
    counter = 1
    length = len(iso_file.read())
    previous_buffer = "000000"
    # search for file signatures
    while current_byte < length:
        print(current_byte, '/', length, '\t',round(float(current_byte)/float(length) * 100,3),"% searched")
        buffer = iso_file.read(3).hex().upper()
        # file has been found
        if(buffer == BEGINING_SIGNATURE and previous_buffer == "000000"):
            # read bytes after the starting signature open file and write all bytes
            current_byte += 3
            iso_file.seek(current_byte)
            buffer1 = iso_file.read(2)
            recovered_file = open("file{}.JPG".format(counter), "wb+")
            recovered_file.write(b'\xFF\xD8\xFF')
            # read 2 bytes at a time until EOF indicator
            while buffer1.hex().upper() != END_SIGNATURE:
                recovered_file.write(buffer1)
                current_byte += 2
                iso_file.seek(current_byte)
                buffer1 = iso_file.read(2)
            recovered_file.write(b'\xFF\xD9')
            recovered_file.close()
            counter += 1
        else:
            current_byte += 3
        previous_buffer = buffer
        iso_file.seek(current_byte)