#include "sys-commands.h"

struct timespec now;
bool inSysCommandMode(struct syscommand_t* command) {
    if (clock_gettime(CLOCK_REALTIME, &now) == -1) {
        return false;
    }

    return command->timeToBeDone.it_value.tv_sec >= now.tv_sec;
}

void updateTimer(struct syscommand_t* command, int timeToAdd) {
    printf("updateTimer\n");
    command->timeToBeDone.it_value.tv_sec += timeToAdd;

    if (timerfd_settime(command->fd, TFD_TIMER_ABSTIME, &command->timeToBeDone, NULL) == -1) {
        printf("Timer update failed horribly\n");
    }

    if (clock_gettime(CLOCK_REALTIME, &now) == -1) {
        printf("clock_gettime has seemed to fail\n");
    }

    printf("New Duration is %lu\n", command->timeToBeDone.it_value.tv_sec - now.tv_sec);
}

bool sysCommandOn(struct syscommand_t* command, int timeToAdd) {
    if (inSysCommandMode(command)) {
        updateTimer(command, timeToAdd);
        return false;
    }

    printf("System Command On: %s\n", command->on);
    system(command->on);
    if (clock_gettime(CLOCK_REALTIME, &now) == -1) {
        printf("clock_gettime has seemed to fail\n");
    }

    /* Create a CLOCK_REALTIME absolute timer with initial
       expiration and interval as specified in command line */

    command->timeToBeDone.it_value.tv_sec = now.tv_sec + timeToAdd;
    command->timeToBeDone.it_value.tv_nsec = now.tv_nsec;

    command->fd = timerfd_create(CLOCK_REALTIME, 0);

    if (command->fd == -1) {
        printf("timerfd_create has failed to create.\n");
    }

    if (timerfd_settime(command->fd, TFD_TIMER_ABSTIME, &command->timeToBeDone, NULL) == -1) {
        printf("SH_T timerfd_settime failed\n");
    }

    return true;
}

//system("setxkbmap us real-prog-dvorak");
bool sysCommandOff(struct syscommand_t* command) {
    printf("System Command Off: %s\n", command->on);
    system(command->off);
    return false;
}

bool isASDFCommand(char* ptr) {
    return strncmp(ptr, "!asdf", 5) == 0;
}

bool isXrandrCommand(char* ptr) {
    return strncmp(ptr, "!xrandr", 7) == 0;
}

bool isSystemCommand(char* ptr) {
    return isASDFCommand(ptr) || isXrandrCommand(ptr);
}
