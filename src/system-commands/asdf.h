#ifndef __ASDF_H_
#define __ASDF_H_

#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

void asdfOn(int timeToAdd);
bool asdfOff();
bool isASDFCommand(char* ptr);
bool isInASDF();
int getFD();

#endif
