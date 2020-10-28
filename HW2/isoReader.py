
from sys import *
from csv import *


class Partition():
    type_description = "Partition Type: "
    s_sector = "Partition Address (LBA): "
    n_sectors = "Number of Sectors in Partition: "
    # GPT only
    name = "Partition Name: "
    guid = "Partition GUID: "
    start = "Partition Starting Address: "
    end = "Partition Ending Address: "

def readMBR(partition_headers):
    # creates table for partition types
    MBR_P_TYPE_BYTE = []
    MBR_P_TYPE_DESCRIPTION = []
    with open('mbr_partition_types.csv', newline='') as data_dump:
        file_reader = reader(data_dump, dialect='excel')
        for row in file_reader:
            MBR_P_TYPE_BYTE.append(int(row[0],16))
            MBR_P_TYPE_DESCRIPTION.append(row[1])
        
    # assign values to each partition
    partitions = []
    for partition_header in partition_headers:
        partition = Partition()
        # assigns description
        for i in range(len(MBR_P_TYPE_BYTE)):
            if(MBR_P_TYPE_BYTE[i] == partition_header[0]):
                partition.type_description += MBR_P_TYPE_DESCRIPTION[i]

        # calculates size and amount of sectors 
        partition.s_sector += str(int(hex(partition_header[6])+"0000",16) + int(hex(partition_header[5])+"00",16) + int(hex(partition_header[4]),16)) # bytes 6+5+4
        partition.n_sectors += str(int(hex(partition_header[10])+"0000",16) + int(hex(partition_header[9])+"00",16) + int(hex(partition_header[8]),16)) # bytes 10+9+8
        partitions.append(partition)

    # Print results
    print("Number of Partitions:", number_of_partitions)
    print()
    for i in range(len(partitions)):
        print("Partition {} Details:".format(i+1))
        print(partitions[i].type_description)
        print(partitions[i].s_sector)
        print(partitions[i].n_sectors)
        print()


        

def readGPT(partition_headers):
    GPT_P_GUID = []
    GPT_P_TYPE_DESCRIPTION = []
    with open('gpt_partition_guids.csv', newline='') as data_dump:
        file_reader = reader(data_dump, dialect='excel')
        for row in file_reader:
            GPT_P_GUID.append(row[0].replace("'", ""))
            GPT_P_TYPE_DESCRIPTION.append(row[2] + " - " + row[1])

    partitionsG = []
    for partition_header in partition_headers:
        i = 0
        partitionG = Partition()
        partitionG.start += str(int.from_bytes(partition_header[32:40], byteorder='little'))
        partitionG.end += str(int.from_bytes(partition_header[40:48], byteorder='little'))
        start = []
        end = []
        g = []
        u = []
        t = []
        d = ""
        f = ""
        while i < len(partition_header):
            # type GUID
            if i < 4:
                if(len(hex(partition_header[i]).replace("0x","")) == 1):
                    g.append(hex(partition_header[i]).replace("0x","0"))
                else:
                    g.append(hex(partition_header[i]).replace("0x",""))
                
            elif i >= 4 and i < 6:
                u.append(hex(partition_header[i]).replace("0x",""))
            elif i >= 6 and i < 8:
                t.append(hex(partition_header[i]).replace("0x",""))
            elif i >= 8 and i < 10:
                d += hex(partition_header[i]).replace("0x","")
            elif i >= 10 and i < 16:
                f += hex(partition_header[i]).replace("0x","")
            # name
            elif i >= 56:
                partitionG.name += chr(partition_header[i])
            i += 1

        # creates GUID
        for x in g[::-1]:
            partitionG.guid += x.upper()
        partitionG.guid += "-"
        for x in u[::-1]:
            partitionG.guid += x.upper()
        partitionG.guid += "-"
        for x in t[::-1]:
            partitionG.guid += x.upper()
        partitionG.guid = partitionG.guid + "-" + d.upper() + "-" + f.upper()

        # assigns description
        for i in range(len(GPT_P_GUID)):
            if(GPT_P_GUID[i] == partitionG.guid.replace("Partition GUID: ","")):
                partitionG.type_description += GPT_P_TYPE_DESCRIPTION[i]
        
        partitionsG.append(partitionG)

    print("Number of Partitions:",number_of_partitions)
    for i in range(len(partitionsG)):
        print()
        print("Partition {} Details:".format(i+1))
        print(partitionsG[i].name)
        print(partitionsG[i].guid)
        print(partitionsG[i].type_description)
        print(partitionsG[i].start)
        print(partitionsG[i].end)
        



# -m (MBR), -g (GPT)
mode = argv[1]
filename = argv[2]

# MBR
if(mode == '-m'):
    iso_file = open(filename, "rb").read(512)
    number_of_partitions = 4
    partition_data = []
    partitions = []

    i = 434 
    j = 0
    #Partition headers run from [434,497] in intervals of 16 bytes
    while(i < 498): 
        # if first byte of partition header is empty then partition is empty
        if(j == 0 and iso_file[i] == 0):
            i += 16
            number_of_partitions -= 1
        else:
            partition_data.append(iso_file[i])
            i += 1
            j += 1
        # adds the partition header to the partition list and reset variables
        if(j == 16):
            partitions.append(partition_data)
            j = 0
            partition_data = []

    readMBR(partitions)

# GPT
elif(mode == '-g'):
    iso_file = open(filename, "rb").read(4096)
    number_of_partitions = 0
    partition_data = []
    partitions = []

    i = 1024 
    j = 0
    #Partition headers run from [1024,3072] in intervals of 128 bytes
    while(i < len(iso_file)): 
        # if first byte of partition header is empty then partition is empty
        if(j == 0 and iso_file[i] == 0):
            i += 128
        else:
            partition_data.append(iso_file[i])
            i += 1
            j += 1
        # adds the partition header to the partition list and reset variables
        if(j == 128):
            partitions.append(partition_data)
            number_of_partitions += 1
            j = 0
            partition_data = []

    readGPT(partitions)
else:
    print("incorrect mode, (-g or -m only)")