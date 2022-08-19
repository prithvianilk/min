./min zip test_1.txt 1.zip
ls -l | grep '1.zip'
ls -l | grep 'test_1.txt'
./min unzip 1.zip output_1.txt
rm 1.zip output_1.txt 