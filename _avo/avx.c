#include <stdint.h>
#include <immintrin.h>

#define BYTE unsigned char
#define U32 u_int32_t
#define U64 u_int64_t

#define KEYSET_DEFAULT_SIZE 48
#define STRIPE_LEN 64
#define STRIPE_ELTS (STRIPE_LEN / sizeof(U32))
#define ACC_NB (STRIPE_LEN / sizeof(U64))
#define NB_KEYS ((KEYSET_DEFAULT_SIZE - STRIPE_ELTS) / 2)
#define PRIME32_1 2654435761

void
XXH3_accumulate_512(void* acc, const void *restrict data, const void *restrict key)
{
    {             __m256i* const xacc  =       (__m256i *) acc;
        const     __m256i* const xdata = (const __m256i *) data;
        const     __m256i* const xkey  = (const __m256i *) key;

        size_t i;
        for (i=0; i < STRIPE_LEN/sizeof(__m256i); i++) {
            __m256i const d   = _mm256_loadu_si256 (xdata+i);
            __m256i const k   = _mm256_loadu_si256 (xkey+i);
            __m256i const dk  = _mm256_add_epi32 (d,k);                                  /* uint32 dk[8]  = {d0+k0, d1+k1, d2+k2, d3+k3, ...} */
            __m256i const res = _mm256_mul_epu32 (dk, _mm256_shuffle_epi32 (dk, 0x31));  /* uint64 res[4] = {dk0*dk1, dk2*dk3, ...} */
            __m256i const add = _mm256_add_epi64(d, xacc[i]);
            xacc[i]  = _mm256_add_epi64(res, add);
        }
    }
}

void
XXH3_scrambleAcc(void* acc, const void* key)
{
    {             __m256i* const xacc = (__m256i*) acc;
        const     __m256i* const xkey  = (const __m256i *) key;
        const __m256i k1 = _mm256_set1_epi32((int)PRIME32_1);

        size_t i;
        for (i=0; i < STRIPE_LEN/sizeof(__m256i); i++) {
            __m256i data = xacc[i];
            __m256i const shifted = _mm256_srli_epi64(data, 47);
            data = _mm256_xor_si256(data, shifted);

            {   __m256i const k   = _mm256_loadu_si256 (xkey+i);
                __m256i const dk  = _mm256_xor_si256   (data, k);          /* U32 dk[4]  = {d0+k0, d1+k1, d2+k2, d3+k3} */

                __m256i const dk1 = _mm256_mul_epu32 (dk, k1);

                __m256i const d2  = _mm256_shuffle_epi32 (dk, 0x31);
                __m256i const dk2 = _mm256_mul_epu32 (d2, k1);
                __m256i const dk2h= _mm256_slli_epi64 (dk2, 32);

                xacc[i] = _mm256_add_epi64(dk1, dk2h);
        }   }
    }
}

void
XXH3_accumulate(U64* acc, const void* restrict data, const U32* restrict key, size_t nbStripes)
{
    size_t n;
    # pragma clang loop unroll(enable)
    for (n = 0; n < nbStripes; n++ ) {
        XXH3_accumulate_512(acc, (const BYTE*)data + n*STRIPE_LEN, key);
        key += 2;
    }
}

void
XXH3_hashLong(U64* acc, const void* data, const void *kKey, size_t len)
{
    #define NB_KEYS ((KEYSET_DEFAULT_SIZE - STRIPE_ELTS) / 2)

    size_t const block_len = STRIPE_LEN * NB_KEYS;
    size_t const nb_blocks = len / block_len;

    size_t n;
    for (n = 0; n < nb_blocks; n++) {
        XXH3_accumulate(acc, (const BYTE*)data + n*block_len, kKey, NB_KEYS);
        XXH3_scrambleAcc(acc, kKey + (KEYSET_DEFAULT_SIZE - STRIPE_ELTS));
    }

    /* last partial block */
    {   size_t const nbStripes = (len % block_len) / STRIPE_LEN;
        XXH3_accumulate(acc, (const BYTE*)data + nb_blocks*block_len, kKey, nbStripes);

        /* last stripe */
        if (len & (STRIPE_LEN - 1)) {
            const BYTE* const p = (const BYTE*) data + len - STRIPE_LEN;
            XXH3_accumulate_512(acc, p, kKey + nbStripes*2);
    }   }
}