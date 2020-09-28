import csv
from datetime import datetime
#read csv file
with open('database_dump.csv','rt',encoding='utf8')as csvfile:
    reader = csv.DictReader(csvfile)
    username = ([row['username'] for row in reader])
for i in len(username):
    username = 
