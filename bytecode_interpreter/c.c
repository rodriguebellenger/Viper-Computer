#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <time.h>
#include <stdint.h>

uint8_t RAM[1024] = {0};

int main() {
    RAM[0] = 200;
    printf("%u ", RAM[0]);
    printf("%u \n", RAM[1]);
    return 0;
}
