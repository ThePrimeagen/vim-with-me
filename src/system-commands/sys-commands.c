#include "hashmap.h"
#include "sys-commands.h"
#include "hashmap.h"
#include <semaphore.h>
#include <pthread.h>

int MAX_SOCKET = 0;
struct timespec now;
pthread_mutex_t isRunningMutex;
bool isRunning = false;
struct hashmap_table* recentUsers;

void sysCommandThrottleUser(int twitchId) {
    if (twitchId && clock_gettime(CLOCK_REALTIME, &now) != -1) {
        hashmap_insert(recentUsers, twitchId, now.tv_sec * 1000);
    }
}

bool inSysCommandMode(struct syscommand_t* command) {
    if (clock_gettime(CLOCK_REALTIME, &now) == -1) {
        return false;
    }

    printf("isSystemCommand %zu >= %zu\n", command->timeToBeDone.it_value.tv_sec, now.tv_sec);
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

bool sysCommandOn(int twitchId, struct syscommand_t* command, int timeToAdd) {
    if (inSysCommandMode(command)) {
        updateTimer(command, timeToAdd);
        return false;
    }

    if (clock_gettime(CLOCK_REALTIME, &now) == -1) {
        printf("clock_gettime has seemed to fail\n");
        return false;
    }

    /* Create a CLOCK_REALTIME absolute timer with initial
       expiration and interval as specified in command line */

    command->timeToBeDone.it_value.tv_sec = now.tv_sec + timeToAdd;
    command->timeToBeDone.it_value.tv_nsec = now.tv_nsec;

    command->fd = timerfd_create(CLOCK_REALTIME, 0);
    MAX_SOCKET = MAX_SOCKET < command->fd ? command->fd : MAX_SOCKET;

    if (command->fd == -1) {
        printf("timerfd_create has failed to create.\n");
        return false;
    }

    if (timerfd_settime(command->fd, TFD_TIMER_ABSTIME, &command->timeToBeDone, NULL) == -1) {
        printf("SH_T timerfd_settime failed %d - %s\n", errno, strerror(errno));
        return false;
    }

    printf("System Command On: %s\n", command->on);
    system(command->on);
    addCommandToFDSelect(command);

    pthread_mutex_lock(&isRunningMutex);
    if (!isRunning) {
        isRunning = true;
        sem_post(command->mutex);
    }
    pthread_mutex_unlock(&isRunningMutex);

    sysCommandThrottleUser(twitchId);
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

bool isVimCommand(char* ptr) {
    return strncmp(ptr, "!vim", 4) == 0 && strlen(ptr) > 5;
}

bool isSystemCommand(char* ptr) {
    return isVimCommand(ptr) ||
        isASDFCommand(ptr) ||
        isXrandrCommand(ptr);
}

const int COMMAND_THROTTLE_MS = 60000;
struct SysCommandList {
    struct syscommand_t* command;
    bool (*callBack)(struct syscommand_t*);
    struct SysCommandList* next;
    struct SysCommandList* prev;
} SysCommandList;

// TODONE: Refactor completed (shortly from now).
// System commands to be better encupsulated with the system-commands interface
struct SysCommandList* head;

bool sysCommandIsThrottled(char* twitchName, int twitchId) {
    printf("sysCommandIsThrottled %s\n", twitchName);
    if (twitchId && clock_gettime(CLOCK_REALTIME, &now) != -1) {
        size_t currentMils = now.tv_sec * 1000;
        size_t lastMillis = hashmap_lookup(recentUsers, twitchId);

        size_t diff = currentMils - lastMillis;
        if (diff < COMMAND_THROTTLE_MS) {
            printf("THROTTLED %s \n", twitchName);
            return true;
        }
    }

    return false;
}

void removeFDNode(struct SysCommandList* node) {
    printf("remove Node\n");

    struct SysCommandList* prev = node->prev;
    struct SysCommandList* next = node->next;

    if (next != NULL) {
        next->prev = prev;
    }

    if (prev != NULL) {
        prev->next = next;
    }

    printf("Checking %p == %p\n", node, head);
    if (node == head) {
        head = head->next;
    }

    free(node);
}

struct SysCommandList* createFDNode() {
    printf("createFDNode\n");

    struct SysCommandList* next = (struct SysCommandList*)malloc(sizeof(struct SysCommandList));
    if (head == NULL) {
        head = next;
    }
    else {
        struct SysCommandList* curr = head;
        while (curr->next != NULL) {
            curr = curr->next;
        }

        curr->next = next;
        next->prev = curr;
    }

    return next;
}

int addCommandToFDSelect(struct syscommand_t* command) {
    printf("Adding command %s\n", command->on);

    struct SysCommandList* node = createFDNode();
    node->command = command;
    node->callBack = &sysCommandOff;

    printf("fileDescriptor %d\n", node->command->fd);

    return node->command->fd;
}

bool hasRepeatCommand(struct vimcommand_t* command) {
    bool hasRepeat = false;
    switch (command->navCommand) {
        case 'j': 
        case 'k':
        case 'l':
        case 'h':
            hasRepeat = true;
            break;
    }

    return hasRepeat;
}

bool isValidVimCommand(struct vimcommand_t* command) {
    bool isValid = false;
    bool isLineColNav = false;

    printf("isValidVimCommand %c %d\n", command->navCommand, command->times);
    switch (command->navCommand) {
        case 'j': // Down
        case 'k': // Up
        case 'l': // Right
        case 'h': // Left
            isLineColNav = true;
            isValid = true;
            break;
        case 'g': // Top of page : I will assume g == gg
        case 'G': // Bottom of page
        case 'v': // Start Visual mode per character.
        case 'V': // Start Visual mode linewise.
        case 'A': // Append text at the end of the line [count] times.
        case 'I': // Beginning of line and insert
        case '~': // Switch case of the character under the cursor and move the cursor to the right [count] times.
        case '%': // Jumps to closing of expression
        case '$': // End of the line and [count - 1] lines downward.
            isValid = true;
            break;

        // Do we allow the D?
        // and can D be equivalent, type wise, to 8?
        // 8===D
        // case 'x':
        // case 's':
    }

    isValid = isValid && (
            (isLineColNav && command->times < 50) || !isLineColNav);

    return isValid;
}

void vimCommandRun(int twitchId, struct vimcommand_t* command) {
    if (!isValidVimCommand(command)) {
        printf("This is not a valid vim command %c\n", command->navCommand);
        return;
    }

    // Sync execution
    char buf[100];
    char navBuf[5];
    navBuf[4] = '\0';
    char* navPtr = navBuf;

    int commandLen = 0;
    if (hasRepeatCommand(command)) {
        commandLen += snprintf(navBuf, 3, "%d", command->times);
    }

    navPtr = navBuf + commandLen;
    navPtr[0] = 'j'; // remove the null terminator

    printf("navBuf (repeat only): %.*s\n", commandLen, navBuf);

    commandLen += snprintf(navPtr, 4 - commandLen, "%c", command->navCommand);

    if (command->navCommand == 'g') {
        (navPtr + 1)[0] = 'g';
        (navPtr + 2)[0] = '\0';
        commandLen += 2;
    }

    printf("navBuf (after repeat only): %.*s\n", commandLen, navBuf);

    int len = snprintf(buf, 100, "vim --remote-send \"<C-c>%s\"", navBuf);
    printf("executing command: %.*s\n", len, buf);

    system(buf);
    sysCommandThrottleUser(twitchId);
}

// Pthread thing here...
// Assumptions.
// 1.  That data is a semaphore.
// 2.  That the semaphore has an initial value of 0
void* runSysCommands(void *dat) {

    recentUsers = hashmap_createTable(1381);
    sem_t* mutex = (sem_t*)dat;

    fd_set commands;

    // 1.  We wait for something to come it.
    // 2.  We
    while (1) {

        printf("I am about to wait like a BA (system command)\n");
        // This really can only be called once... until there are no more system
        // commands ready to be called
        sem_wait(mutex);
        printf("I AM DONE WAITING\n");

        while (head != NULL) {
            struct SysCommandList* curr = head;

            FD_ZERO(&commands);

            while (curr != NULL) {
                FD_SET(curr->command->fd, &commands);
                curr = curr->next;
            }

            // Wait for one of the time FDs to go off
            if (select(MAX_SOCKET + 1, &commands, 0, 0, 0) == -1) {
                printf("WHAT JUST HAPPENED????\n");
                printf("My select statement has completely failed, therefore there is one and only thing left to do.\n");
                exit(-1);
            }

            curr = head;

            while (curr != NULL) {

                struct SysCommandList* next = curr->next;

                // if this is the one up then we need to remove it from the
                // list.
                if (FD_ISSET(curr->command->fd, &commands)) {
                    printf("We have something to callback\n");
                    if (!curr->callBack(curr->command)) {
                        removeFDNode(curr);
                    }
                }

                curr = next;
            }
        }

        pthread_mutex_lock(&isRunningMutex);
        isRunning = false;
        pthread_mutex_unlock(&isRunningMutex);
    }

    return NULL;
}

