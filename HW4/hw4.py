from sys import *

filename = open("FAT_Imposter.iso", "rb")

buffer = filename.read(512)
bytesPERsector = int.from_bytes(buffer[11:13], byteorder='little')
print("Bytes/Sector: {}".format(bytesPERsector))
sectorsPERcluster = int.from_bytes(buffer[13:14], byteorder='little')
print("Sectors/Cluster: {}".format(sectorsPERcluster))
sizeReserved = int.from_bytes(buffer[14:16], byteorder='little')
print("Size of Reserved Area in Sectors: {}".format(sizeReserved))
print("Start Address of 1st FAT: {}".format(sizeReserved))
numFATs = int.from_bytes(buffer[16:17], byteorder='little')
print("# of FATs: {}".format(numFATs))
secPerFAT = int.from_bytes(buffer[36:40], byteorder='little')
print("Sectors/FAT: {}".format(secPerFAT))
cluster_address_RD = int.from_bytes(buffer[44:48], byteorder='little')
print("Cluster Address of Root Directory: {}".format(cluster_address_RD))
s_sector_address_RD = numFATs * secPerFAT + sizeReserved
print("Starting Sector Address of the Data Section: {}".format(s_sector_address_RD))

length = len(filename.read())
current_byte = s_sector_address_RD * bytesPERsector
filename.seek(current_byte)

while current_byte < length:
    buffer = filename.read(32)
    current_byte += 32
    attributes = buffer[11:12].hex()
    if attributes == "20":
        try:
            nameOfFile = buffer[0:8].decode("utf-8")
            extension = buffer[8:11].decode("utf-8")
            segmentTag = ""
            first_part_cluster = int.from_bytes(buffer[20:22], byteorder='little')
            second_part_cluster = int.from_bytes(buffer[26:28], byteorder='little')
            file_cluster_address = first_part_cluster + second_part_cluster
            size_of_file = int.from_bytes(buffer[28:32], byteorder='little')
        except:
            continue
        
        if extension != "JPG":
            continue

        cluster_chain = []
        cluster_chain.append(file_cluster_address)
        current_cluster = file_cluster_address
       
        while True:
            filename.seek((sizeReserved * bytesPERsector) + (current_cluster * 4))
            buffer = filename.read(4)
            current_cluster = int.from_bytes(buffer, byteorder='little')
            cluster_chain.append(current_cluster)
            if buffer.hex() == "00000000" or buffer.hex() == "ffffff0f":
                break

        recovered_file = open(nameOfFile, "wb+")
        for i in range(len(cluster_chain)):
            next_section_address = ((cluster_chain[i] - cluster_address_RD) + s_sector_address_RD)
            next_section_address_bytes = next_section_address * bytesPERsector
            filename.seek(next_section_address_bytes, 0)
            buffer = filename.read(bytesPERsector)
            if i == 0:
                FileHeader = buffer[:16].hex()
                if FileHeader.find("383961") == 6:
                    segmentTag = "GIF"
                    recovered_file.write(b'\x47\x49\x46')
                    buffer = buffer[3:]
                elif FileHeader.find("462d") == 6:
                    segmentTag = "PDF"
                    recovered_file.write(b'\x25\x50\x44')
                    buffer = buffer[3:]
                elif FileHeader.find("470d") == 6:
                    segmentTag = "PNG"
                    recovered_file.write(b'\x89\x50\x4e')
                    buffer = buffer[3:]
                elif FileHeader.find("e0a1") == 6:
                    segmentTag = "DOC"
                    recovered_file.write(b'\xd0\xcf\x11')
                    buffer = buffer[3:]
                elif FileHeader.find("041400") == 6:
                    segmentTag = "JAR"
                    recovered_file.write(b'\x50\x4b\x03')
                    buffer = buffer[3:]
                elif FileHeader.find("4f4354595045") == 6:
                    segmentTag = "HTML"
                    recovered_file.write(b'\x3c\x21\x44')
                    buffer = buffer[3:]
                elif "4a46494600" in FileHeader:
                    segmentTag = "JIF"
                elif "4578696600" in FileHeader:
                    segmentTag = "EXIF"
                elif "535049464600" in FileHeader:
                    segmentTag = "SPIFF"
                else:
                    segmentTag = "UNKNOWN"

            recovered_file.write(buffer)


        print("-----------------------------------------")
        print("Name: {}".format(nameOfFile) )
        print("Extension: {}".format(extension))
        print("Type: {}".format(segmentTag))
        print("Attributes: {}".format(attributes))
        print("Cluster Address: {}".format(file_cluster_address))
        print("Size: {}".format(size_of_file))
       
        filename.seek(current_byte)