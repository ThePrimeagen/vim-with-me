#ifndef __HASH_MAP__H_
#define __HASH_MAP__H_

#include <stdio.h>
#include <stdlib.h>

struct hashmap_node {
    int key;
    size_t val;
    struct hashmap_node *next;
};

struct hashmap_table {
    int size;
    struct hashmap_node **list;
};

struct hashmap_table* hashmap_createTable(int size);
void hashmap_insert(struct hashmap_table *t,int key, size_t val);
size_t hashmap_lookup(struct hashmap_table *t,int key);

#endif
