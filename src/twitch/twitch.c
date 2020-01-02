#include "twitch.h"

char nameParts[3] = {'@', '!', ':'};

const char* listenCommand = "{\
    \"type\": \"LISTEN\",\
    \"data\": {\
        \"topics\": [\"%s\"],\
        \"auth_token\": \"%s\"\
    }\
}";

/**
 * don't forget to free....
 */
char* twitchListenCommand(const char* topic, const char* authToken) {
    size_t listenCommandLen = strlen(listenCommand) - 4;
    size_t size = strlen(topic) + strlen(authToken) + listenCommandLen + 1;
    char* buf = (char*)malloc(size);

    snprintf(buf, size, listenCommand, topic, authToken);
    return buf;
}

/*
  @badge-info=subscriber/3;badges=vip/1,subscriber/3,partner/1;color=#000000;
  display-name=Heph;emotes=1:44-45;flags=;id=7265f6cd-87fe-46e1-b69c-5481502cd9f7;
  mod=0;room-id=167160215;subscriber=1;tmi-sent-ts=1577909246349;turbo=0;user-id=29181352;user-type=
  :heph!heph@heph.tmi.twitch.tv PRIVMSG #theprimeagen :ill probably end up

  doing something similar :)
 * */
char* twitchReadNameFromIRC(const char* ircTagInfo, char* msgStart) {
    char* nameStart = msgStart;
    int idx = 0;
    do {
        --nameStart;
        if (nameStart[0] == nameParts[idx]) {
            nameStart[0] = '\0';
            idx++;
        }

    } while (idx != 3 && nameStart != ircTagInfo);
    nameStart++;

    if (nameStart == ircTagInfo + 1) {
        return NULL;
    }

    return nameStart;
}

/*
  @badge-info=subscriber/3;badges=vip/1,subscriber/3,partner/1;color=#000000;
  display-name=Heph;emotes=1:44-45;flags=;id=7265f6cd-87fe-46e1-b69c-5481502cd9f7;
  mod=0;room-id=167160215;subscriber=1;tmi-sent-ts=1577909246349;turbo=0;user-id=29181352;user-type=
  :heph!heph@heph.tmi.twitch.tv PRIVMSG #theprimeagen :ill probably end up

  doing something similar :)
 * */
int twitchReadUserId(const char* ircTagInfo) {
    char* userId = strstr(ircTagInfo, "user-id=");
    printf("twitchReadUserId: %.*s\n%.*s\n", 50, ircTagInfo, 25, userId);

    if (userId == NULL) {
        return 0;
    }

    do {
        ++userId;
    } while ((*userId) != '=');
    ++userId;

    char* end = userId;
    do {
        --end;
    } while ((*end) != ' ' && (*end) != ';');

    end[0] = '\0';
    printf("final results: %d - %s\n", atoi(userId), userId);

    return atoi(userId);
}
