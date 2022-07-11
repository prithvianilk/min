#include "node.h"

namespace huffman {
node::node(char token, int freq)
    : token(token), freq(freq), left(nullptr), right(nullptr) {}

node::node(node* left, node* right) : left(left), right(right) {
  freq = left->freq + right->freq;
}
}  // namespace huffman
