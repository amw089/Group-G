import uuid
import hashlib
#store the uuid in list
result = []
i=0
#convert name to uuid  
test_uuid = uuid.UUID("d9b2d63d-a233-4123-847a-76838bf2413a")
for line in open("names.txt","r"):
    result.append(str(uuid.uuid5(test_uuid,line))+"\n")

#store those uuid in new .txt file 
with open('Unames.txt','w')as f2:
    for i in result:
        f2.writelines(str(i))


