#pragma once

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdbool.h>

char* twitchReadNameFromIRC(const char* irc, char* msgStart);
int twitchReadUserId(const char* irc);
bool isHighlightedMessage(char* tags);
