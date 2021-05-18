#include "irc.h"
#include "twitch-irc.h"
#include "sys-commands.h"
#include "hashmap.h"

const int ASDF_TIME = 3;
const int XRANDR_TIME = 4;

struct syscommand_t* asdf;
struct syscommand_t* lightMeSilly;

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


void sendMessage(int sock, const char* line, int len) {

    char buffer[len + 25];
    snprintf(buffer, len + 25, "PRIVMSG theprimeagen :%s\r\n", line);

    send(sock, buffer, len + 25, 0);
}

void set_nick(int sock, char nick[]) {
    char nick_packet[512];
    sprintf(nick_packet, "NICK %s\r\n", nick);
    send(sock, nick_packet, strlen(nick_packet), 0);
}

void set_tags(int sock) {
    char pass_packet[512];

    // TODO: When I am a real man
    sprintf(pass_packet, "CAP REQ :twitch.tv/tags\r\n");

    send(sock, pass_packet, strlen(pass_packet), 0);
}


void SET_PASS(int sock, char pass[]) {
    char pass_packet[512];
    // TODO: When I am a real man
    sprintf(pass_packet, "PASS %s\r\n", pass);

    printf("PASS %.*s", 8, pass_packet);

    send(sock, pass_packet, strlen(pass_packet), 0);
}

void send_user_packet(int sock, char nick[]) {
    char user_packet[512];
    sprintf(user_packet, "USER %s 0 * :%s\r\n", nick, nick);
    send(sock, user_packet, strlen(user_packet), 0);
}

void join_channel(int sock, char channel[]) {
    char join_packet[512];
    sprintf(join_packet, "JOIN %s\r\n", channel);
    send(sock, join_packet, strlen(join_packet), 0);
}

void send_pong(int sock, char argument[]) {
    char pong_packet[512];
    sprintf(pong_packet, "PONG :%s\r\n", argument);
    send(sock, pong_packet, strlen(pong_packet), 0);
}

void send_message(int sock, char to[], char message[]) {
    char message_packet[512];
    sprintf(message_packet, "PRIVMSG %s :%s\r\n", to, message);
    send(sock, message_packet, strlen(message_packet), 0);
}

bool isPRIVMSG(char* lineOffset) {
    return strncmp(lineOffset, "PRIVMSG", 7) == 0;
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

    printf("IRC: Just read %.*s\n", (int)length, buffer);
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

void handleHighlightedMessage(char* lineOffset, char* twitchName, int twitchId) {

    if (sysCommandIsThrottled(twitchName, twitchId)) {
        return;
    }

    int len = strlen("PRIVMSG #theprimeagen :");
    char* msg = lineOffset + len;

    if (isASDFCommand(msg)) {
        sysCommandOn(twitchId, asdf, ASDF_TIME);
    }
    else if (isXrandrCommand(msg)) {
        sysCommandOn(twitchId, lightMeSilly, XRANDR_TIME);
    }
    else if (isVimCommand(msg)) {
        struct vimcommand_t cmd;
        parseVimCommand(msg, &cmd);

        vimCommandRun(twitchId, &cmd);
    }
}

bool isStatusLine(char* line) {
    return strstr(line, "!vimstatus") != NULL;
}

void handleIRC(int socket_desc, char* line, int bytesRead) {
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

    /*
    if (isStatusLine(lineOffset)) {
        sendMessage(socket_desc, "vim-with-me is up", 17);
    }
    */

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
    lightMeSilly->on = "xrandr --output HDMI-0 --brightness 0.05";
    lightMeSilly->off = "xrandr --output HDMI-0 --brightness 1";
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



