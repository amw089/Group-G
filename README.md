Matt's Branch

I have all of my homework and challenge files uploaded here


Homework 2:

iso_reader.py runs with either -g or -m (python3 iso_reader.iso -g gpt_dump.iso)
FATreader.py (python3 FATreader.py FAT_FS.iso) (preloaded offset and filters for challenge 2)

Homework 3:

fileCarver.py runs with either -f or -d 
-d : searches data section and finds 44 files (1 fragmented only first sector) files found denoted with file_at_-d.png
-f: searches the FAT and finds 36 files (1 corrupted) files found denoted file_at_-f.png

Homework 4:

python3 fileCarver.py FAT_imposter.iso : finds the single JFIF file and saves its data in a file

Challenge 3:

Finds the files iterativly, I believe conflicting ending signaturs results in slightly corrupted text documents
