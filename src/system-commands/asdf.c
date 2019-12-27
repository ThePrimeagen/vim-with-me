#include "asdf.h"

void addTime(int timeToAdd) {
    timeLeft += timeToAdd;
}

bool isInASDF() {
    return timeLeft > 0;
}

void asdf(void* ptr) {
    system("setxkbmap us");

    int timeToAdd = *((int*)ptr);
    timeLeft += timeToAdd;

    while (true) {
        sleep(1);
        if (--timeLeft <= 0) {
            break;
        }
    }

    system("setxkbmap us real-prog-dvorak");
}



bool isASDFCommand(char* ptr) {
    return strncmp(ptr, "!asdf", 5) == 0;
}
