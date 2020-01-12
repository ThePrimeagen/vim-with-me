#pragma once

#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>
#include <sys/timerfd.h>
#include <sys/select.h>
#include <semaphore.h>

struct vimcommand_t {
    char navCommand;
    int times;
    //
    // This later............... (15 minutes).
    struct itimerspec timeToBeExecuted;
} vimcommand_t;

struct syscommand_t {
    struct itimerspec timeToBeDone;
    int fd;
    const char* on;
    const char* off;
    sem_t* mutex;
} syscommand_t;

bool inSysCommandMode(struct syscommand_t* command);
bool sysCommandOn(int twitchId, struct syscommand_t* command, int secondsToAdd);
bool sysCommandOff(struct syscommand_t* command);
bool isASDFCommand(char* ptr);
bool isVimCommand(char* ptr);
bool isXrandrCommand(char* ptr);
bool isSystemCommand(char* ptr);
bool isPointCheck(char* ptr);
int addCommandToFDSelect(struct syscommand_t* command);
bool sysCommandIsThrottled(char* twitchName, int twitchId);
void* runSysCommands(void *dat);
bool vimCommandRun(int twitchId, struct vimcommand_t* command);

