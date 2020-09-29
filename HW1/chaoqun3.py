import hashlib
import sys
# It is a good idea to store the filename into a variable.
# The variable can later become a function argument when the
# code is converted to a function body.
filename = ('D1.txt')

# Using the newer with construct to close the file automatically.
with open(filename) as f:
    data = f.readlines()

# Or using the older approach and closing the filea explicitly.
# Here the data is re-read again, do not use both ;)
f = open(filename)
data = f.readlines()
f.close()


# The data is of the list type.  The Python list type is actually
# a dynamic array. The lines contain also the \n; hence the .rstrip()
for n, line in enumerate(data, 1):
    print ('{:2}.'.format(n), line.rstrip())

print ('-----------------')

# You can later iterate through the list for other purpose, for
# example to read them via the csv.reader.
import csv
result = []
reader = csv.reader(data)
for row in reader:
    result[row] = hashlib.sha256(data.encode('utf-8')).hexdigest()
    print(result)
