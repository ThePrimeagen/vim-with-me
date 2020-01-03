#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <time.h>
#include <errno.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <sys/timerfd.h>
#include <sys/select.h>

#include "sys-commands.h"
#include "json-c/json.h"
#include "hashmap.h"
#include "twitch-irc.h"

const int ASDF_TIME = 3;
const int XRANDR_TIME = 4;
const int COMMAND_THROTTLE_MS = 60000;

int MAX_SOCKET = 0;

// TODO: Refactor:
// System commands to be better encupsulated with the system-commands interface
struct SysCommandList {
    struct syscommand_t* command;
    bool (*callBack)(struct syscommand_t*);
    struct SysCommandList* next;
    struct SysCommandList* prev;
} SysCommandList;

struct SysCommandList* head;
struct syscommand_t* asdf;
struct syscommand_t* lightMeSilly;
struct hashmap_table* recentUsers;

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

char* get_config(char *name) {
    char* value = (char*)malloc(1024);
    FILE *configfile = fopen(".config", "r");
    value[0] = '\0';

    if (configfile != NULL) {
        while (1) {
            char configname[1024];
            char tempvalue[1024];

            int status = fscanf(configfile, " %1023[^= ] = %s ", configname, tempvalue); //Parse key=value

            if (status == EOF) {
                break;
            }

            if (strcmp(configname, name) == 0) {
                strncpy(value, tempvalue, strlen(tempvalue)+1);
                break;
            }
        }
        fclose(configfile);
    }

    return value;
}

char* get_command(char line[]) {
    char* command = (char*)malloc(512);
    char clone[512];
    strncpy(clone, line, strlen(line)+1);
    char* splitted = strtok(clone, " ");
    if (splitted != NULL) {
        if (splitted[0] == ':') {
            splitted = strtok(NULL, " ");
        }
        if (splitted != NULL) {
            strncpy(command, splitted, strlen(splitted)+1);
        } else{
            command[0] = '\0';
        }
    } else{
        command[0] = '\0';
    }
    return command;
}

char* get_last_argument(char line[]) {
    char* argument = (char*)malloc(512);
    char clone[512];
    strncpy(clone, line, strlen(line)+1);
    char* splitted = strstr(clone, " :");
    if (splitted != NULL) {
        strncpy(argument, splitted+2, strlen(splitted)+1);
    } else{
        argument[0] = '\0';
    }
    return argument;
}

void addCommandToFDSelect(struct syscommand_t* command) {
    struct SysCommandList* node = createFDNode();
    node->command = command;
    node->callBack = &sysCommandOff;

    printf("fileDescriptor %d\n", node->command->fd);

    if (MAX_SOCKET < node->command->fd) {
        MAX_SOCKET = node->command->fd;
    }
}

void handleHighlightedMessage(char* lineOffset, char* twitchName, int twitchId) {
    struct timespec now;

    if (twitchId && clock_gettime(CLOCK_REALTIME, &now) != -1) {
        size_t currentMils = now.tv_sec * 1000;
        size_t lastMillis = hashmap_lookup(recentUsers, twitchId);

        size_t diff = currentMils - lastMillis;
        if (diff < COMMAND_THROTTLE_MS) {
            printf("THROTTLED %s \n", twitchName);
            return;
        }
    }

    int len = strlen("PRIVMSG #theprimeagen :");
    char* msg = lineOffset + len;

    bool addToFD = false;
    struct syscommand_t* cmd = NULL;

    if (isASDFCommand(msg)) {
        addToFD = sysCommandOn(asdf, ASDF_TIME);
        cmd = asdf;
    }
    else if (isXrandrCommand(msg)) {
        addToFD = sysCommandOn(lightMeSilly, XRANDR_TIME);
        cmd = lightMeSilly;
    }

    if (addToFD) {
        addCommandToFDSelect(cmd);
    }

    if (cmd != NULL && twitchId) {
        size_t currentMils = now.tv_sec * 1000;
        hashmap_insert(recentUsers, twitchId, currentMils);
    }
}
void handleIRC(int socket_desc, char* line) {
    int bytesRead = read_line(socket_desc, line);
    (void)bytesRead;

    printf("XXXX - Incoming line %.*s\n", bytesRead, line);

    char* lineOffset = line;
    int count = 0;
    if (line[0] == '@') {
        do {
            lineOffset++;
            count++;
        } while (!isPRIVMSG(lineOffset));
    }

    char* command = get_command(lineOffset);

    if (strcmp(command, "PING") == 0) {
        char* argument = get_last_argument(lineOffset);
        send_pong(socket_desc, argument);
        free(argument);
    }

    if (isHighlightedMessage(line)) {

        int twitchId = twitchReadUserId(line);
        char* twitchName = twitchReadNameFromIRC(line, lineOffset);

        handleHighlightedMessage(lineOffset, twitchName, twitchId);
    }

    free(command);
}

int main() {
    recentUsers = hashmap_createTable(1381);

    asdf = (struct syscommand_t*)malloc(sizeof(syscommand_t));
    asdf->on = "setxkbmap us";
    asdf->off = "setxkbmap us real-prog-dvorak";

    lightMeSilly = (struct syscommand_t*)malloc(sizeof(syscommand_t));
    lightMeSilly->on = "xrandr --output HDMI-1 --brightness 0.05";
    lightMeSilly->off = "xrandr --output HDMI-1 --brightness 1";

    int socket_desc = socket(AF_INET, SOCK_STREAM, 0);
    if (socket_desc == -1) {
        perror("Could not create socket");
        exit(1);
    }

    char* ip = get_config((char*)"server");
    char* port = get_config((char*)"port");

    printf("IP: %s -- Port: %s\n", ip, port);

    struct sockaddr_in server;
    server.sin_addr.s_addr = inet_addr(ip);
    server.sin_family = AF_INET;
    server.sin_port = htons(atoi(port));

    free(ip);
    free(port);

    if (connect(socket_desc, (struct sockaddr *) &server, sizeof(server)) < 0) {
        perror("Could not connect");
        exit(1);
    }

    char* nick = get_config((char*)"nick");
    char* pass = get_config((char*)"pass");
    char* channels = get_config((char*)"channels");

    // TODO: Don't forget to do that one thing you are suppose to do, but you
    // are clearly to lazy to do anything.
    SET_PASS(socket_desc, pass);
    set_nick(socket_desc, nick);
    send_user_packet(socket_desc, nick);
    join_channel(socket_desc, channels);
    set_tags(socket_desc);

    free(nick);
    free(channels);

    char line[1024 * 10];
    fd_set reads;

    MAX_SOCKET = socket_desc;

    int loopCount = 0;
    while (++loopCount < 250000) {
        printf("Going on another loop %d\n", loopCount);

        FD_ZERO(&reads);
        FD_SET(socket_desc, &reads);

        printf("clearing structs %d\n", loopCount);
        struct SysCommandList* curr = head;
        while (curr != NULL) {
            FD_SET(curr->command->fd, &reads);
            curr = curr->next;
        }

        printf("select %d\n", loopCount);
        if (select(MAX_SOCKET + 1, &reads, 0, 0, 0) < 0) {
            printf("failed#select %d - %s\n", errno, strerror(errno));
            return 1;
        }

        // Read the IRC stuffs
        if (FD_ISSET(socket_desc, &reads)) {
            handleIRC(socket_desc, line);
        }

        curr = head;
        printf("while curr != null %p %d\n", head, loopCount);
        while (curr != NULL) {
            printf("checking FD %d\n", curr->command->fd);
            struct SysCommandList* next = curr->next;

            if (FD_ISSET(curr->command->fd, &reads)) {
                printf("We have something to callback\n");
                if (!curr->callBack(curr->command)) {
                    removeFDNode(curr);
                }
            }

            curr = next;
        }
    }
}

