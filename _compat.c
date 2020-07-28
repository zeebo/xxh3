#include "upstream/xxhash.h"
#include <stdio.h>

int main() {
    unsigned char buf[1024];
    for (int i = 0; i < 1024; i++) {
        buf[i] = (unsigned char)i;

        uint64_t h = XXH3_64bits(buf, (size_t)i);
        printf("\t%llu,\n", h);
    }
}
