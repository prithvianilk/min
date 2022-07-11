#pragma once

namespace huffman {
struct node {
  char token;
  int freq;
  node *left, *right;

  node(char, int);
  node(node *, node *);
  bool operator>(node *rhs) { return freq > rhs->freq; }
};

struct node_cmp {
  constexpr bool operator()(node *left, node *right) {
    return left->freq > right->freq;
  }
};

}  // namespace huffman
