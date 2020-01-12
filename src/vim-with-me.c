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
#include <semaphore.h>
#include <pthread.h>

#include "irc.h"
#include "sys-commands.h"
#include "primepoints.h"

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

int main(int argc, char* argv[]) {

    if (argc != 2) {
        printf("COME ON PRIME.... Provide the points file brahahahahahha\n");
        exit(1);
    }

    struct ircConfig* config = (struct ircConfig*)malloc(sizeof(struct ircConfig));
    config->ip = get_config((char*)"server");
    config->port = get_config((char*)"port");
    config->nick = get_config((char*)"nick");
    config->pass = get_config((char*)"pass");
    config->channels = get_config((char*)"channels");

    printf("Yeah, ip: %s\n", config->ip);
    printf("Yeah, port: %s\n", config->port);
    printf("Yeah, nick: %s\n", config->nick);
    printf("Yeah, channels: %s\n", config->channels);

    printf("Getting userdata %s\n", argv[1]);
    FILE* f = fopen(argv[1], "r");

    fseek(f, 0, SEEK_END);
    long fsize = ftell(f);
    fseek(f, 0, SEEK_SET);  /* same as rewind(f); */

    char* string = (char*)malloc(fsize + 1);
    fread(string, 1, fsize, f);
    fclose(f);

    printf("Look at this file %d - %s\n", (int)strlen(string), string);
    primePointInit(string, argv[1]);

    sem_t s;
    sem_init(&s, 0, 0);

    // TODO: Pthread that mumbo jumbo
    pthread_t ircThread;
    pthread_t sysThread;

    config->system = &s;

    int err = pthread_create(&sysThread, NULL, runSysCommands, &s);
    if (err) {
        printf("Could not create the sys thread %d-%s\n", err, strerror(err));
        return -1;
    }

    err = pthread_create(&ircThread, NULL, ircRun, config);
    if (err) {
        printf("Could not create the irc thread %d-%s\n", err, strerror(err));
        return -1;
    }

    pthread_join(sysThread, NULL);
    pthread_join(ircThread, NULL);

    // GOOD GUY MAIN
    free(config->ip);
    free(config->port);
    free(config->nick);
    free(config->pass);
    free(config->channels);
    free(config);

    return 0;
}

