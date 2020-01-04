#pragma once

#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

// Networking side of the things.
#include <errno.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <sys/timerfd.h>
#include <sys/select.h>
#include <semaphore.h>

struct ircConfig {
    char* ip;
    char* port;
    char* nick;
    char* pass;
    char* channels;
    sem_t* system;
};

/**
 * runs.  Thats it.
 */
void* ircRun(void* dat);

