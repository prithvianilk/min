#pragma once

#include <filesystem>
#include <functional>
#include <queue>
#include <string>
#include <unordered_map>
#include <vector>

#include "node.h"

namespace fs = std::filesystem;

namespace huffman {
class huffman {
 private:
  std::string input;
  std::unordered_map<char, int> freq_map;
  std::priority_queue<node*, std::vector<node*>, node_cmp> queue;
  std::unordered_map<char, std::string> enc_map;

  node* gen_tree();
  void dfs(node*, std::string);

 public:
  void compress(const fs::path&, const fs::path&);
  void decompress(const fs::path&, const fs::path&);
};
}  // namespace huffman
