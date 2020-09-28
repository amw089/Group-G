import datetime
import csv

with open('database_dump.csv', newline='') as data_dump:
    reader = csv.reader(data_dump, dialect='excel')
    for row in reader:
        time_stamp = row[2]
        if(time_stamp != "last_access"):
            time_stamp = datetime.datetime.fromtimestamp(int(time_stamp), datetime.datetime.timezone.UTC)
            print(time_stamp)