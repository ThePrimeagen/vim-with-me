#pragma once

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

/**
 * runs.  Thats it.
 */
void primePointInit(char* contents, char* filepath);
void addPrimePoints(char* name, int points);
void removePrimePoints(char* name, int points);
bool hasEnoughPrimePoints(char* name, int points);
int getPrimePoints(char* name);

