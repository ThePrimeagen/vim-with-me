#include <ctype.h>
#include "twitch-irc.h"

char nameParts[3] = {'@', '!', ':'};

bool isHighlightedMessage(char* tags) {
    // Tags @badge-info=subscriber/16;badges=broadcaster/1,subscriber/6,sub-gifter/50;color=#FF0000;display-name=ThePrimeagen;emotes=;flags=;id=e6b18363-a1a3-4d96-96dd-9cc288fdca39;mod=0;msg-id=highlighted-message;room-id=167160215;subscriber=1;tmi-sent-ts=1577486185539;turbo=0;user-id=167160215;user-type= :theprimeagen!theprimeagen@theprimeagen.tmi.twitch.tv

    return strstr(tags, "msg-id=highlighted-message") != NULL;
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

    char* p = nameStart;
    for ( ; *p; ++p) *p = tolower(*p);

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

    return atoi(userId);
}
