#include <ctype.h>

#include "irc.h"
#include "twitch-irc.h"
#include "sys-commands.h"
#include "hashmap.h"
#include "primepoints.h"

#define DEBUG true

const int ASDF_TIME = 3;
const int XRANDR_TIME = 4;

struct syscommand_t* asdf;
struct syscommand_t* lightMeSilly;

const int STOP_HURTING_MY_FEELINGS = 1024 * 10;

bool isPrime(char* twitchName, char* tags) {
    return strncmp(twitchName, "theprimeagen", 12) == 0 &&
       strstr(tags, "badges=broadcaster/1") != NULL;
}

char* get_command(char line[]) {
    char* command = (char*)malloc(STOP_HURTING_MY_FEELINGS);
    char clone[STOP_HURTING_MY_FEELINGS];
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


void set_nick(int sock, char nick[]) {
    char nick_packet[STOP_HURTING_MY_FEELINGS];
    sprintf(nick_packet, "NICK %s\r\n", nick);
    send(sock, nick_packet, strlen(nick_packet), 0);
}

void set_tags(int sock) {
    char pass_packet[STOP_HURTING_MY_FEELINGS];

    // TODO: When I am a real man
    sprintf(pass_packet, "CAP REQ :twitch.tv/tags\r\n");

    send(sock, pass_packet, strlen(pass_packet), 0);
}


void SET_PASS(int sock, char pass[]) {
    char pass_packet[STOP_HURTING_MY_FEELINGS];
    // TODO: When I am a real man
    sprintf(pass_packet, "PASS %s\r\n", pass);

    printf("PASS %.*s", 8, pass_packet);

    send(sock, pass_packet, strlen(pass_packet), 0);
}

void send_user_packet(int sock, char nick[]) {
    char user_packet[STOP_HURTING_MY_FEELINGS];
    sprintf(user_packet, "USER %s 0 * :%s\r\n", nick, nick);
    send(sock, user_packet, strlen(user_packet), 0);
}

void join_channel(int sock, char channel[]) {
    char join_packet[STOP_HURTING_MY_FEELINGS];
    sprintf(join_packet, "JOIN %s\r\n", channel);
    send(sock, join_packet, strlen(join_packet), 0);
}

void send_pong(int sock, char argument[]) {
    char pong_packet[STOP_HURTING_MY_FEELINGS];
    sprintf(pong_packet, "PONG :%s\r\n", argument);
    send(sock, pong_packet, strlen(pong_packet), 0);
}

void send_message(int sock, char to[], char message[]) {
    char message_packet[STOP_HURTING_MY_FEELINGS];
    sprintf(message_packet, "PRIVMSG %s :%s\r\n", to, message);
    send(sock, message_packet, strlen(message_packet), 0);
}

bool isPRIVMSG(char* lineOffset) {
    return strncmp(lineOffset, "PRIVMSG", 7) == 0;
}

char* get_last_argument(char line[]) {
    char* argument = (char*)malloc(STOP_HURTING_MY_FEELINGS);
    char clone[STOP_HURTING_MY_FEELINGS];

    strncpy(clone, line, strlen(line)+1);
    char* splitted = strstr(clone, " :");
    if (splitted != NULL) {
        strncpy(argument, splitted+2, strlen(splitted)+1);
    } else{
        argument[0] = '\0';
    }
    return argument;
}

// PTHREAD MAGIC
int read_line(int sock, char* buffer) {
    size_t length = 0;

    while (1) {
        char data;
        int result = recv(sock, &data, 1, 0);

        if ((result <= 0) || (data == EOF)) {
            perror("Connection closed");
            exit(1);
        }

        buffer[length] = data;
        length++;

        if (length >= 2 && buffer[length-2] == '\r' && buffer[length-1] == '\n') {
            buffer[length-2] = '\0';
            break;
        }
    }

    return length;
}

void parseVimCommand(char* msg, struct vimcommand_t* command) {
    // "!vim "
    msg += 5;
    char* navCommand = msg;

    command->navCommand = navCommand[0];

    if (strlen(navCommand) == 1) {
        return;
    }

    command->times = atoi(navCommand + 2);
}

bool isPointBuyMessage(char* msg) {
    return strncmp(msg, "!buy", 4) == 0;
}

const char* PRIVATE_MSG = "PRIVMSG #theprimeagen :%s\r\n";
void sendMessage(int sockDesc, char* msg) {
    int len = strlen(msg) + strlen(PRIVATE_MSG) - 2;
    char buf[len];

    int out = sprintf(buf, PRIVATE_MSG, msg);

    if (out != len) {
        printf("Lengths don't match on strings... %d %d\n", len, out);
    }

    out = send(sockDesc, buf, out, 0);
    printf("sendMessage#send: %d\n", out);
}

void handleSystemRequest(int socketDesc, char* msg, char* twitchName, int twitchId) {
    printf("Handling System Request: %s - %s\n", twitchName, msg);

    if (isPointCheck(msg)) {
        char buf[strlen(twitchName) + 4];
        sprintf(buf, "%s: %d", twitchName, getPrimePoints(twitchName));

        sendMessage(socketDesc, buf);
        return;
    }

    if (sysCommandIsThrottled(twitchName, twitchId)) {
#if DEBUG
        printf("handleSystemRequest#!sysCommandIsThrottled(%s) \n", twitchName);
#endif
        return;
    }

    if (!hasEnoughPrimePoints(twitchName, 100)) {
#if DEBUG
        printf("handleSystemRequest#!hasEnoughPrimePoints(%s, 100) \n", twitchName);
#endif
        return;
    }

    if (isASDFCommand(msg) && sysCommandOn(twitchId, asdf, ASDF_TIME)) {
        removePrimePoints(twitchName, 100);
    }

    else if (isXrandrCommand(msg) && sysCommandOn(twitchId, lightMeSilly, XRANDR_TIME)) {
        removePrimePoints(twitchName, 100);
    }

    else if (isVimCommand(msg)) {
        struct vimcommand_t cmd;
        parseVimCommand(msg, &cmd);

        if (vimCommandRun(twitchId, &cmd)) {
            removePrimePoints(twitchName, 100);
        }
    }
}

bool isSystemCommandRequest(char* msg) {
#ifdef DEBUG
    printf("isSystemCommandRequest(%s): %d %d %d %d\n", msg, isASDFCommand(msg) ,
            isXrandrCommand(msg) , isPointCheck(msg) , isVimCommand(msg));
#endif
    return isASDFCommand(msg) ||
        isXrandrCommand(msg) ||
        isPointCheck(msg) ||
        isVimCommand(msg);
}

void handleHighlightedMessage(char* msg, char* twitchName) {
    if (isPointBuyMessage(msg)) {
        addPrimePoints(twitchName, 100);
    }
}

bool isPrimeAdd(char* msg) {
    return strncmp(msg, "!add", 4) == 0;
}

bool isPrimeGW(char* msg) {
    return strncmp(msg, "!gw", 3) == 0;
}

bool isSpecialCommandMessage(char* msg) {
    return  isPrimeAdd(msg) || isPrimeGW(msg);
}

char* parseUserFromPrimeCommand(char* cmd) {
    char* firstSpace = strstr(cmd, " ");
    if (firstSpace == NULL) {
        return NULL;
    }

    firstSpace++;

    char* secondSpace = strstr(firstSpace, " ");
    if (secondSpace == NULL) {
        return NULL;
    }

    secondSpace[0] = 0;


    char* p = firstSpace;
    for ( ; *p; ++p) *p = tolower(*p);

    return firstSpace;
}

void handlePrimeCommand(int socket_desc, char* msg) {
    (void)socket_desc;

    printf("PrimeCommandTime: Lets command together %s\n", msg);

    if (isPrimeAdd(msg)) {
        int msgLen = strlen(msg);
        char* user = parseUserFromPrimeCommand(msg);

        if (user == NULL) {
            return;
        }

        printf("strlen: %d + 1 + 4 - %d = \n", (int)strlen(user), (int)msgLen);

        if (strlen(user) + 1 + 4 - msgLen <= 0) {
            return;
        }

        int arg = atoi(user + strlen(user) + 1);
        printf("Here is the primeCommand %s %d\n", user, arg);

        addPrimePoints(user, arg);
    }
}

void handleIRC(int socket_desc, char* line, int bytesRead) {
    (void)bytesRead;

    if (DEBUG) {
        printf("XXXX - Incoming line %.*s\n", bytesRead, line);
    }

    char* lineOffset = line;
    int count = 0;
    if (line[0] == '@') {
        do {
            lineOffset++;
            count++;
        } while (!isPRIVMSG(lineOffset));
    }

    char* command = get_command(lineOffset);

    lineOffset--;
    lineOffset[0] = 0;
    lineOffset++;

    if (strcmp(command, "PING") == 0) {
        char* argument = get_last_argument(lineOffset);
        send_pong(socket_desc, argument);
        free(argument);
    }

    if (line == lineOffset) {
        printf("Ignoring this line %s\n", line);
        return;
    }

    int twitchId = twitchReadUserId(line);
    char* twitchName = twitchReadNameFromIRC(line, lineOffset);

    int len = strlen("PRIVMSG #theprimeagen :");
    char* msg = lineOffset + len;

    printf("isSystemCommandRequest(%s): %d\n", msg, isSystemCommandRequest(msg));

    if (isPrime(twitchName, line) && isSpecialCommandMessage(msg)) {
        handlePrimeCommand(socket_desc, msg);
    }

    else if (isHighlightedMessage(line)) {
        handleHighlightedMessage(msg, twitchName);
    }

    else if (isSystemCommandRequest(msg)) {
        handleSystemRequest(socket_desc, msg, twitchName, twitchId);
    }

    free(command);
}

void* ircRun(void* dat) {
    struct ircConfig* config = (struct ircConfig*)dat;
    asdf = (struct syscommand_t*)malloc(sizeof(syscommand_t));
    memset(asdf, 0, sizeof(syscommand_t));
    asdf->on = "setxkbmap us";
    asdf->off = "setxkbmap us real-prog-dvorak";
    asdf->mutex = config->system;

    lightMeSilly = (struct syscommand_t*)malloc(sizeof(syscommand_t));
    memset(lightMeSilly, 0, sizeof(syscommand_t));
    lightMeSilly->on = "xrandr --output HDMI-1 --brightness 0.05";
    lightMeSilly->off = "xrandr --output HDMI-1 --brightness 1";
    lightMeSilly->mutex = config->system;

    printf("IP: %s -- Port: %s\n", config->ip, config->port);
    int socket_desc = socket(AF_INET, SOCK_STREAM, 0);
    if (socket_desc == -1) {
        perror("Could not create socket");
        exit(1);
    }

    struct sockaddr_in server;
    server.sin_addr.s_addr = inet_addr(config->ip);
    server.sin_family = AF_INET;
    server.sin_port = htons(atoi(config->port));

    if (connect(socket_desc, (struct sockaddr *) &server, sizeof(server)) < 0) {
        perror("Could not connect");
        exit(1);
    }

    // TODO: Don't forget to do that one thing you are suppose to do, but you
    // are clearly to lazy to do anything.
    SET_PASS(socket_desc, config->pass);
    set_nick(socket_desc, config->nick);
    send_user_packet(socket_desc, config->nick);
    join_channel(socket_desc, config->channels);
    set_tags(socket_desc);

    char line[1024 * 10];

    while (1) {
        handleIRC(socket_desc, line, read_line(socket_desc, line));
    }

    return NULL;
}



