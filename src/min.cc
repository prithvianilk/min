#include <string.h>

#include <iostream>

#include "huffman/huffman.h"

void print_usage() {
  std::cout << "Usage:\n";
  std::cout << "\tmin zip <source-path> <zip-path>\n";
  std::cout << "\tmin unzip <zip-path> <dest-path>\n";
}

int main(int argc, const char** argv) {
  huffman::huffman h;
  if (argc != 4) {
    print_usage();
    return 1;
  }
  if (strcmp(argv[1], "zip") == 0) {
    h.compress(argv[2], argv[3]);
  } else if (strcmp(argv[1], "unzip") == 0) {
    h.decompress(argv[2], argv[3]);
  } else {
    print_usage();
    return 1;
  }
  return 0;
}