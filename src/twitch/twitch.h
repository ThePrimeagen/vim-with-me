#ifndef __TWITCH_COMMANDS_H_
#define __TWITCH_COMMANDS_H_

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stddef.h>

char* twitchListenCommand(const char* topic, const char* auth);
char* twitchReadNameFromIRC(const char* irc, char* msgStart);
int twitchReadUserId(const char* irc);

#endif



