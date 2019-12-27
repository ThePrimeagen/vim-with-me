#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <time.h>
#include <sys/socket.h>
#include <arpa/inet.h>

#include "system-commands/asdf.h"

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
            return length;
        }
    }
}

void log_with_date(char* line) {
    char date[50];
    struct tm *current_time;

    time_t now = time(0);
    current_time = gmtime(&now);
    strftime(date, sizeof(date), "%Y-%m-%d %H:%M:%S", current_time);

    printf("[%s] %s\n", date, line);
}

void log_to_file(char* line, FILE *logfile) {
    char date[50];
    struct tm *current_time;

    time_t now = time(0);
    current_time = gmtime(&now);
    strftime(date, sizeof(date), "%Y-%m-%d %H:%M:%S", current_time);

    fprintf(logfile, "[%s] %s\n", date, line);
    fflush(logfile);
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

char* get_prefix(char *line) {
    char* prefix = (char*)malloc(512);
    char clone[512];

    strncpy(clone, line, strlen(line) + 1);
    if (line[0] == ':') {
        char* splitted = strtok(clone, " ");
        if (splitted != NULL) {
            strncpy(prefix, splitted+1, strlen(splitted)+1);
        } else {
            prefix[0] = '\0';
        }
    } else {
        prefix[0] = '\0';
    }
    return prefix;
}

char* get_username(char *line) {
    char* username = (char*)malloc(512);
    char clone[512];

    strncpy(clone, line, strlen(line) + 1);
    if (strchr(clone, '!') != NULL) {
        char* splitted = strtok(clone, "!");
        if (splitted != NULL) {
            strncpy(username, splitted+1, strlen(splitted)+1);
        } else {
            username[0] = '\0';
        }
    } else {
        username[0] = '\0';
    }
    return username;
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

char* get_argument(char line[], int argno) {
    char* argument = (char*)malloc(512);
    char clone[512];
    strncpy(clone, line, strlen(line)+1);

    int current_arg = 0;
    char* splitted = strtok(clone, " ");
    while (splitted != NULL) {
        if (splitted[0] != ':') {
            current_arg++;
        }
        if (current_arg == argno+1) {
            strncpy(argument, splitted, strlen(splitted)+1);
            return argument;
        }
        splitted = strtok(NULL, " ");
    }

    if (current_arg != argno) {
        argument[0] = '\0';
    }
    return argument;
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

bool isHighlightedMessage(char* tags) {
    // Tags @badge-info=subscriber/16;badges=broadcaster/1,subscriber/6,sub-gifter/50;color=#FF0000;display-name=ThePrimeagen;emotes=;flags=;id=e6b18363-a1a3-4d96-96dd-9cc288fdca39;mod=0;msg-id=highlighted-message;room-id=167160215;subscriber=1;tmi-sent-ts=1577486185539;turbo=0;user-id=167160215;user-type= :theprimeagen!theprimeagen@theprimeagen.tmi.twitch.tv

    return strstr(tags, "msg-id=highlighted-message") != NULL;
}

int main() {
    int socket_desc = socket(AF_INET, SOCK_STREAM, 0);
    if (socket_desc == -1) {
       perror("Could not create socket");
       exit(1);
    }

    char* ip = get_config((char*)"server");
    char* port = get_config((char*)"port");

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
    while (1) {

        int bytesRead = read_line(socket_desc, line);

        char* lineOffset = line;
        int count = 0;
        if (line[0] == '@') {
            do {
                lineOffset++;
                count++;
            } while (!isPRIVMSG(lineOffset));
        }

        char* prefix = get_prefix(lineOffset);
        char* username = get_username(lineOffset);
        char* command = get_command(lineOffset);
        char* argument = get_last_argument(lineOffset);

        if (strcmp(command, "PING") == 0) {
            send_pong(socket_desc, argument);
            log_with_date((char*)"Got ping. Replying with pong...");
        }

        else if (isHighlightedMessage(line)) {
            int len = strlen("PRIVMSG #theprimeagen :");
            char* msg = lineOffset + len;

            printf("Did this even work?\n");
            if (strncmp(msg, "!asdf", 5) == 0) {
                system("setxkbmap us");
            }
        }

        free(prefix);
        free(username);
        free(command);
        free(argument);
    }
}
