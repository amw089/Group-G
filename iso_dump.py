# Katie Hay and Anna Wolf
# Partition analysis
# To run: iso_dump.py -<mbr/gpt> <iso file>

import sys, struct, binascii
from binascii import hexlify

# Function to get number of partitions in MBR file
def mbr_number_of_partitions(file):
    with open(file, "rb") as f:
        mbr_data = f.read(512)
        start = 446
        end = 462
        num_partitions = 0
        while(num_partitions <= 3):
            # Get number of partitions
            partition = struct.unpack("<BBBBBBBBBBBBBBBB", mbr_data[start:end])
            if(partition[1] != 0) and (partition[2] != 0) and (partition[3] != 0) and (partition[4] != 0):
                num_partitions += 1
            else:
                break
            start = end
            end = end + 16
        return num_partitions

# Function to parse MBR data
def mbr_parse(partiton):
    partition_type = mbr_partition_type(partition[4])
    partition_address = mbr_partition_size(partition[8:12])
    num_sectors = mbr_partition_size(partition[12:16])
    return partition_type, partition_address, num_sectors
        
# Function to determine size of MBR partition
def mbr_partition_size(length):
    b1 = struct.pack("<B", length[0])
    b2 = struct.pack("<B", length[1])
    b3 = struct.pack("<B", length[2])
    b4 = struct.pack("<B", length[3])
    size = b1 + b2 + b3 + b4
    size = struct.unpack("<L", size)[0]
    return size

# Function to identify MBR partition type
def mbr_partition_type(part_type):
    part_type = "0x{:02x}".format(part_type).upper()
    with open ("mbr_partition_types.csv", "r", encoding="UTF8") as f:
        f.readline()
        for row in f:
            if (row[0:2] == part_type[2:]):
                return (row[3:-1])

# Function to get number of partitions in GPT file
def gpt_number_of_partitions(file):
    with open(file, "rb") as f:
        gpt_data = f.read(10240)
        start = 1024
        end = 1152
        num_partitions = 0
        while (num_partitions <= 127):
            try:
                partition = struct.unpack("<QQQQQQQQQQQQQQQQ", gpt_data[start:end])
                num_partitions += 1
            except struct.error:
                break
            start = end
            end = end + 128
        return num_partitions

# Function to parse GPT data
def gpt_parse(partiton):
    partition_name = gpt_partition_name(partition)
    partition_guid = gpt_partition_guid(hex(partition[0]).upper(), hex(partition[1]).upper())
    partition_type = gpt_partition_type(partition_guid)
    return partition_guid, partition_type, partition_name
    

# Function to identify GPT partition type
def gpt_partition_name(name_data):
    final = ""
    name = hex(name_data[8])
    name = name[2:]
    final = final + bytes.fromhex('{}'.format(name)).decode('utf-8')
    name = hex(name_data[7])
    name = name[2:]
    final = final + bytes.fromhex('{}'.format(name)).decode('utf-8')
    name = hex(name_data[9])
    name = name[2:]
##    final = final + bytes.fromhex('{}'.format(name)).decode('utf-8')
    return final

def gpt_partition_guid(data, data2):
    half = data[-8:] + "-" + data[-12:-8] + "-" + data[-16:-12] + "-"
    guid = half + data2[-2:] + data2[-4:-2] + "-" + data2[-6:-4] + data2[-8:-6]  + data2[-10:-8]  + data2[-12:-10]  + data2[-14:-12]  + data2[-16:-14]  
    return guid

def gpt_partition_type(guid):
    with open ("gpt_partition_guids.csv", "r") as f:
        f.readline()
        for row in f:
            if(row[:36] == guid):
                return row[37:-1]







# Checks arguments
if (len(sys.argv) != 3):
    print("Usage: iso_dump.py -mbr/gpt <.iso file>")
else:
    # Opens .iso file and initiates parsing based on mode (MBR/GPT)    
    if (sys.argv[1] == "-mbr"):
        # Print out number of partitions
        num = mbr_number_of_partitions(sys.argv[2])
        print("Number of partitions: {}".format(num))
        # Print out details of each partition
        with open(sys.argv[2], "rb") as f:
            mbr_data = f.read(512)
            start = 446
            end = 462
            partition_num = 1
            while(partition_num <= num):
                partition = struct.unpack("<BBBBBBBBBBBBBBBB", mbr_data[start:end])
                partition_data = mbr_parse(partition)
                print("\nPartition {} Details:".format(partition_num))
                print("Partition Type: {}".format(partition_data[0]))
                print("Partition Address (LBA): {}".format(partition_data[1]))
                print("Number of Sectors in Partition: {}".format(partition_data[2]))
                partition_num += 1
                start = end
                end = end + 16

    if (sys.argv[1] == "-gpt"):
        # Print out number of partitions
        num = gpt_number_of_partitions(sys.argv[2])
        print("Number of partitions: {}".format(num))
        # Print out details of each partition
        with open(sys.argv[2], "rb") as f:
            gpt_data = f.read(10240)
            start = 1024
            end = 1152
            partition_num = 1
            while (partition_num <= num):
                partition = struct.unpack("<QQQQQQQQQQQQQQQQ", gpt_data[start:end])
                partition_data = gpt_parse(partition)
                print("\nPartition {} Details:".format(partition_num))
                print("Partition Name: {}".format(partition_data[2]))
                print("Partition GUID: {}".format(partition_data[0]))
                print("Partition Type: {}".format(partition_data[1]))
                print("Partition Starting Address: {}".format(partition[4]))
                print("Partition Ending Address: {}".format(partition[5]))
                partition_num += 1
                start = end
                end = end + 128
                
                
                          
              
