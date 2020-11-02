from difflib import SequenceMatcher

EXTENSION_FILTER = "JPG"
SIGNATURES = [[b'\x47\x49\x46\x38\x39\x61',"GIF"],[b'\x25\x50\x44\x46\x2d',"PDF"],[b'\x89\x50\x4e\x47\x0d',"PNG"],[b'\xd0\xcf\x11\xe0\xa1',"DOC"],[b'\x50\x4b\x03\x04\x14\x00',"JAR"],[b'\x3c\x21\x44\x4f\x43\x54\x59\x50\x45',"HTML"],[b'\x46\x49\x46\x00',"JIF"],[b'\x45\x78\x69\x66',"EXIF"],[b'\x53\x50\x49\x46\x46',"SPIFF"]]
DEBUG = True

def CompareSignatures( header ):
    signatureOptions = [0,""]
    signaturePrediction = [b'\x00\x00\x00',"Unknown"]
    for signature in SIGNATURES:
        seq = SequenceMatcher(a=header, b=signature[0].hex())
        if seq.ratio() > 0.20:
            if seq.ratio() > signatureOptions[0]:
                if DEBUG:
                    print(header)
                    print(signature)
                    print(seq.ratio())
                signatureOptions = [seq.ratio(),signature[1]]
                signaturePrediction = signature  
    
    print("\n************************************")
    print("Predicting file signature: {}".format(signaturePrediction))
    print("************************************")

    return signaturePrediction


file = open("FAT_Imposter.iso", "rb")

buffer = file.read(512)
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

length = len(file.read())
current_byte = s_sector_address_RD * bytesPERsector
file.seek(current_byte)

while current_byte < length:
    buffer = file.read(32)
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
        
        if extension != EXTENSION_FILTER:
            continue

        cluster_chain = []
        cluster_chain.append(file_cluster_address)
        current_cluster = file_cluster_address
       
        while True:
            file.seek((sizeReserved * bytesPERsector) + (current_cluster * 4))
            buffer = file.read(4)
            current_cluster = int.from_bytes(buffer, byteorder='little')
            cluster_chain.append(current_cluster)
            if buffer.hex() == "00000000" or buffer.hex() == "ffffff0f":
                break

        recovered_file = open(nameOfFile, "wb+")
        for i in range(len(cluster_chain)):
            next_section_address = ((cluster_chain[i] - cluster_address_RD) + s_sector_address_RD)
            next_section_address_bytes = next_section_address * bytesPERsector
            file.seek(next_section_address_bytes, 0)
            buffer = file.read(bytesPERsector)
            if i == 0:
                FileHeader = buffer[:16].hex()
                if "4a46494600" in FileHeader:
                    segmentTag = "JIF"
                elif "4578696600" in FileHeader:
                    segmentTag = "EXIF"
                elif "535049464600" in FileHeader:
                    segmentTag = "SPIFF"
                else:
                    SignaturePrediction = CompareSignatures(FileHeader[:16])
                    segmentTag = SignaturePrediction[1]
                    recovered_file.write(SignaturePrediction[0])
                    buffer = buffer[len(SignaturePrediction[0]):]
                
                

            recovered_file.write(buffer)


        print("Name: {}".format(nameOfFile) )
        print("Extension: {}".format(extension))
        print("Type: {}".format(segmentTag))
        print("Attributes: {}".format(attributes))
        print("Cluster Address: {}".format(file_cluster_address))
        print("Size: {}".format(size_of_file))
        print("-----------------------------------------")

        file.seek(current_byte)