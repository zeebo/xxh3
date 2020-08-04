#include "upstream/xxhash.h"
#include <stdio.h>

int main() {
    unsigned char buf[4096];
    for (int i = 0; i < 4096; i++) {
        buf[i] = (unsigned char)((i+1)%251);

        uint64_t h = XXH3_64bits(buf, (size_t)i);
        printf("\t0x%llx, ", h);

        if (i % 4 == 3) {
            printf("\n");
        }
    }
}
