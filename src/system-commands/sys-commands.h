#ifndef __ASDF_H_
#define __ASDF_H_

#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/timerfd.h>
#include <sys/select.h>

struct syscommand_t {
    struct itimerspec timeToBeDone;
    int fd;
    const char* on;
    const char* off;
} syscommand_t;

bool inSysCommandMode(struct syscommand_t* command);
bool sysCommandOn(struct syscommand_t* command, int secondsToAdd);
bool sysCommandOff(struct syscommand_t* command);
bool isASDFCommand(char* ptr);
bool isXrandrCommand(char* ptr);
bool isSystemCommand(char* ptr);

#endif
