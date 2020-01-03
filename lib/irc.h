#pragma once

#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Networking side of the things.
#include <errno.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <sys/timerfd.h>
#include <sys/select.h>


struct ircConfig {
    char* ip;
    char* port;
    char* nick;
    char* pass;
    char* channels;
    pthread_t thread;
};

/**
 * runs.  Thats it.
 */
int ircRun(struct ircConfig* config);

