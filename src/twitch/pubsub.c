#include "twitch-pubsub.h"

const char* listenCommand = "{\
    \"type\": \"LISTEN\",\
    \"data\": {\
        \"topics\": [\"%s\"],\
        \"auth_token\": \"%s\"\
    }\
}";

int twitchListenCommand(const char* topic, const char* authToken, char* buf) {
    size_t listenCommandLen = strlen(listenCommand) - 4;
    size_t size = strlen(topic) + strlen(authToken) + listenCommandLen + 1;

    return snprintf(buf, size, listenCommand, topic, authToken);
}

