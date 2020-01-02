#include "hashmap.h"

struct hashmap_table* hashmap_createTable(int size) {
    struct hashmap_table *t = (struct hashmap_table*)malloc(sizeof(struct hashmap_table));
    t->size = size;
    t->list = (struct hashmap_node**)malloc(sizeof(struct hashmap_node*)*size);
    int i;
    for(i=0;i<size;i++) {
        t->list[i] = NULL;
    }
    return t;
}

int hashmap_hashCode(struct hashmap_table *t,int key) {
    if(key<0)
        return -(key%t->size);
    return key%t->size;
}

void hashmap_insert(struct hashmap_table *t,int key, size_t val) {
    int pos = hashmap_hashCode(t,key);
    struct hashmap_node *list = t->list[pos];
    struct hashmap_node *newNode = (struct hashmap_node*)malloc(sizeof(struct hashmap_node));
    struct hashmap_node *temp = list;
    while(temp){
        if(temp->key==key){
            temp->val = val;
            return;
        }
        temp = temp->next;
    }
    newNode->key = key;
    newNode->val = val;
    newNode->next = list;
    t->list[pos] = newNode;
}
size_t hashmap_lookup(struct hashmap_table *t,int key){
    int pos = hashmap_hashCode(t,key);
    struct hashmap_node *list = t->list[pos];
    struct hashmap_node *temp = list;
    while(temp){
        if(temp->key==key){
            return temp->val;
        }
        temp = temp->next;
    }
    return -1;
}
