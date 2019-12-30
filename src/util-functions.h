#ifndef __UTIL_FUNCTIONS_H_
#define __UTIL_FUNCTIONS_H_

#include <sys/time.h>
#include <stdlib.h>

// Straight up stole from SO
// https://stackoverflow.com/questions/3756323/how-to-get-the-current-time-in-milliseconds-from-c-in-linux
long long currentTimestamp() {

    struct timeval te;

    // get current time
    gettimeofday(&te, NULL);

    // calculate milliseconds
    long long milliseconds = te.tv_sec * 1000LL + te.tv_usec / 1000;

    return milliseconds;
}

#endif

