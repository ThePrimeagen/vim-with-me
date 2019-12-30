#include <sys/timerfd.h>
#include <sys/select.h>

#include "asdf.h"
#include "../util-functions.h"

struct itimerspec asdfTime;
bool inASDF = false;
int asdfFd;

bool isInASDF() {
    return inASDF;
}

int getFD() {
    return asdfFd;
}

void updateTimer(int timeToAdd) {
    printf("updateTimer\n");
    asdfTime.it_value.tv_sec += timeToAdd;
    if (timerfd_settime(asdfFd, TFD_TIMER_ABSTIME, &asdfTime, NULL) == -1) {
        printf("Timer update failed horribly\n");
    }

    struct timespec now;

    if (clock_gettime(CLOCK_REALTIME, &now) == -1) {
        printf("clock_gettime has seemed to fail\n");
    }

    printf("New Duration is %lu\n", asdfTime.it_value.tv_sec - now.tv_sec);
}

void asdfOn(int timeToAdd) {

    printf("ASDFOF\n");

    if (inASDF) {
        updateTimer(timeToAdd);
        return;
    }

    system("setxkbmap us");

    struct timespec now;

    if (clock_gettime(CLOCK_REALTIME, &now) == -1) {
        printf("clock_gettime has seemed to fail\n");
    }

    /* Create a CLOCK_REALTIME absolute timer with initial
       expiration and interval as specified in command line */

    asdfTime.it_value.tv_sec = now.tv_sec + timeToAdd;
    asdfTime.it_value.tv_nsec = now.tv_nsec;

    asdfFd = timerfd_create(CLOCK_REALTIME, 0);

    if (asdfFd == -1) {
        printf("timerfd_create has failed to create.\n");
    }

    if (timerfd_settime(asdfFd, TFD_TIMER_ABSTIME, &asdfTime, NULL) == -1) {
        printf("SH_T timerfd_settime failed\n");
    }
    inASDF = true;
}

bool asdfOff() {
    printf("Turning aoeu back on\n");

    struct timespec now;

    if (clock_gettime(CLOCK_REALTIME, &now) == -1) {
        printf("clock_gettime has seemed to fail\n");
    }

    if (now.tv_sec >= asdfTime.it_value.tv_sec) {
        inASDF = false;
        system("setxkbmap us real-prog-dvorak");
    }

    return false;
}

bool isASDFCommand(char* ptr) {
    return strncmp(ptr, "!asdf", 5) == 0;
}
