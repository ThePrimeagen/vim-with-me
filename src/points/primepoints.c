#include "primepoints.h"
#include <stdio.h>
#include <stdlib.h>
#include <ctype.h>

#include "json-c/json.h"

json_object* users;
char* filePath;

int getPrimePoints(char* name) {
    json_object* tmp;
    json_object_object_get_ex(users, name, &tmp);

    printf("getPrimePoints %s %d\n", name, json_object_get_int(tmp));
    return json_object_get_int(tmp);
}

bool hasEnoughPrimePoints(char* name, int points) {
    return getPrimePoints(name) >= points;
}

void removePrimePoints(char* name, int points) {
    addPrimePoints(name, -1 * points);
}

void addPrimePoints(char* name, int points) {
    if (name[0] == '@') {
        name += 1;
    }

    json_object* tmp;
    json_object_object_get_ex(users, name, &tmp);
    int currPoints = json_object_get_int(tmp);

    json_object* nextValue = json_object_new_int(currPoints + points);
    json_object_object_add(users, name, nextValue);

    FILE* f = fopen(filePath, "w");

    fprintf(f, "%s\n",
        json_object_to_json_string_ext(
            users, JSON_C_TO_STRING_SPACED | JSON_C_TO_STRING_PRETTY));

    fclose(f);
    free(tmp);
}

void primePointInit(char* contents, char* f) {
    users = json_tokener_parse(contents);
    filePath = f;
}

