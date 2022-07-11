CC = g++-11
FLAGS = -std=c++20 -Wall

.PHONY: min clean

min: src/min.cc obj/huffman.o obj/node.o
	$(CC) $(FLAGS) -o min src/min.cc obj/huffman.o obj/node.o

obj/huffman.o: src/huffman/huffman.cc src/huffman/huffman.h
	$(CC) $(FLAGS) -o obj/huffman.o -c src/huffman/huffman.cc 

obj/node.o: src/huffman/node.cc src/huffman/node.h
	$(CC) $(FLAGS) -o obj/node.o -c src/huffman/node.cc 

clean: 
	rm obj/*.o
	rm min