#include "irc.h"
#include "twitch-irc.h"

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

