#include "huffman.h"

#include <filesystem>
#include <fstream>
#include <ios>
#include <iostream>
#include <queue>
#include <string>
#include <vector>

#include "node.h"

namespace fs = std::filesystem;

namespace huffman {

constexpr auto write_mode =
    std::ios_base::out | std::ios_base::trunc | std::ios_base::binary;

constexpr auto read_mode = std::ios_base::in | std::ios_base::binary;

void huffman::compress(const fs::path& path, const fs::path& zip_path) {
  std::ifstream ifile(path, read_mode);
  long long input_size = 0;
  while (ifile) {
    char token;
    ifile.get(token);
    ++freq_map[token];
    ++input_size;
  }
  ifile.close();
  std::ofstream file(zip_path, write_mode);
  auto root = gen_tree();
  dfs(root, "");
  size_t enc_map_len = enc_map.size();
  file.write((char*)(&enc_map_len), sizeof(size_t));
  for (auto [token, enc] : enc_map) {
    file.write((char*)(&token), sizeof(char));
    size_t enc_len = enc.size();
    file.write((char*)(&enc_len), sizeof(size_t));
    file.write((char*)(enc.c_str()), enc.size());
  }
  ifile.open(path);
  file.write((char*)(&input_size), sizeof(long long));
  while (ifile) {
    char token;
    ifile.get(token);
    auto enc = enc_map[token];
    size_t enc_len = enc.size();
    file.write((char*)(&enc_len), sizeof(size_t));
    file.write((char*)(enc.c_str()), enc.size());
  }
  ifile.close();
  file.close();
}

node* huffman::gen_tree() {
  for (auto [token, freq] : freq_map) {
    queue.push(new node(token, freq));
  }
  while (queue.size() > 1) {
    auto left = queue.top();
    queue.pop();
    auto right = queue.top();
    queue.pop();
    queue.push(new node(left, right));
  }
  return queue.top();
}

void huffman::dfs(node* node, std::string enc) {
  if (node == nullptr) {
    return;
  }
  bool is_leaf = node->left == nullptr && node->right == nullptr;
  if (is_leaf) {
    enc_map[node->token] = enc;
    return;
  }
  dfs(node->left, enc + "0");
  dfs(node->right, enc + "1");
}

void huffman::decompress(const fs::path& zip_path, const fs::path& path) {
  std::ifstream ifile(zip_path, read_mode);
  size_t enc_map_len;
  ifile.read((char*)(&enc_map_len), sizeof(size_t));
  std::unordered_map<std::string, char> inv_enc_map;
  for (size_t i = 0; i < enc_map_len; ++i) {
    char token;
    size_t enc_arr_len;
    ifile.read((char*)(&token), sizeof(char));
    ifile.read((char*)(&enc_arr_len), sizeof(size_t));
    char* enc_arr = new char[enc_arr_len];
    ifile.read(enc_arr, enc_arr_len * sizeof(char));
    std::string enc(enc_arr, enc_arr + enc_arr_len);
    delete[] enc_arr;
    inv_enc_map[enc] = token;
  }
  long long input_size;
  ifile.read((char*)(&input_size), sizeof(long long));
  std::ofstream ofile(path, write_mode);
  for (int i = 0; i < input_size; ++i) {
    size_t enc_arr_len;
    ifile.read((char*)(&enc_arr_len), sizeof(size_t));
    char* enc_arr = new char[enc_arr_len];
    ifile.read(enc_arr, enc_arr_len * sizeof(char));
    std::string enc(enc_arr, enc_arr + enc_arr_len);
    char token = inv_enc_map[enc];
    ofile.write((char*)(&token), sizeof(char));
    delete[] enc_arr;
  }
  ifile.close();
  ofile.close();
}

}  // namespace huffman
