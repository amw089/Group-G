import uuid
import csv
import hashlib
from datetime import datetime

#DEFINITIONS

#iterating through excel
def reading_csv(n):
    f = open(n)
    csv_f = csv.reader(f)
    for row in csv_f:
        return (row[0])
def sha256hash(n):
    m = hashlib.sha256("b"+n).hexdigest()
    print(m)

    
#creates a uuid based on the hardcoded namespace from input "n"
def create(n):
    x = uuid.UUID('d9b2d63d-a233-4123-847a-76838bf2413a')
    return (str((uuid.uuid5(x, n))))

#opens text file, goes line by line, and turns names into uuids     
def make(filename):
    qbfile = open(filename, "r")
    for line in qbfile:
        stripped_line = line.strip()
        create(stripped_line)

#looks through dictionary, turns each line into uuid, finds matches in cvs file
def scan_dictionary(dictfilename, csvfilename):
    qbfile = open(dictfilename, "r")
    for line in qbfile:
        stripped_line = line.strip()
        t = create(stripped_line)
        f = open(csvfilename)
        csv_f = csv.reader(f)
        for row in csv_f:
            if(t == str(row[0])):
                print(line + " = " + t)

#looks through dictionary, hashes each line, finds matches for passwords in cvs file
def pscan_dictionary(dictfilename, csvfilename):
    qbfile = open(dictfilename, "rb")
    for line in qbfile:
        stripped_line = line.strip()
        m = hashlib.sha256(stripped_line).hexdigest()
        f = open(csvfilename)
        csv_f = csv.reader(f)
        for row in csv_f:
            if(str.upper(m) == (str(row[1]))):
                print(str(line) + " = " + str.upper(m))

def pscan_dictionary2(dictfilename, passhash):
    qbfile = open(dictfilename, "rb")
    for line in qbfile:
        stripped_line = line.strip()
        m = hashlib.sha256(stripped_line).hexdigest()
        if (str.upper(m) == str(passhash)):
            print(str(line) + " = " + str.upper(m))

def timestamp(n):
    dt_object = datetime.fromtimestamp(n)
    print(dt_object)

#translates timestamps from cvs file into readable
def timestamp2(csvfilename):
    f = open(csvfilename)
    csv_f = csv.reader(f)
    next(csv_f)
    for row in csv_f:
        #print (str(row[2]))
        st = int(row[2])
        #print(st)
        dt_object = datetime.fromtimestamp(st)
        print(dt_object)

def timestamp3(csvfilename, num):
    f = open(csvfilename)
    csv_f = csv.reader(f)
    st = int(row[2])
    dt_object = datetime.fromtimestamp(st)
    print(dt_object)
        

#writing csv file
##def write_csv(csvfilename, newfilename):
##    with open(csvfilename, mode = 'w') as csv file:
##        fieldnames = ['username', 'password', 'last_access']
##        writer = csv.DictWriter(csv_file, fieldnames=fieldnames)

#combine all three
def combination(dictfilename, csvfilename):
    qbfile = open(dictfilename, "r")
    f = open(csvfilename)
    csv_f = csv.reader(f)
    next(csv_f)
    for row in csv_f:
        for line in qbfile:
            stripped_line = line.strip()
            t = create(stripped_line)
            if(t == str(row[0])):
                print(line + " = " + t)
                dang = row[1]
                pscan_dictionary2(dictfilename, dang)
                st = int(row[2])
                dt_object = datetime.fromtimestamp(st)
                print(dt_object)
                break
            else:
                print("nope")
                
            
                
                
            
                
                
combination("dictionary.txt", "database_dump.csv")
                
  
    

        
#IMPLEMENTATION
    
#pscan_dictionary2("dictionary.txt", "database_dump.csv")
#timestamp3("database_dump.csv", 3)

