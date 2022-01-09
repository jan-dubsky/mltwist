#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

// Source: https://gist.github.com/foobaz/3287f153d125277eefea
uint32_t gapsqrt64(uint64_t a) {
	uint64_t rem = 0, root = 0;

	for (int i = 64 / 2; i > 0; i--) {
		root <<= 1;
		rem = (rem << 2) | (a >> (64 - 2));
		a <<= 2;
		if (root < rem) {
			rem -= root | 1;
			root += 2;
		}
	}
	return (uint32_t)(root >> 1);
}

bool is_prime(const uint64_t num) {
	for (uint64_t i = 2; i < gapsqrt64(num); i++) {
		if ((num % i) == 0) return false;
	}
	return true;
}

void find_primes(const uint64_t cnt) {
	printf("%llu\n", 2);

	for (uint64_t i = 3, found = 1; found < cnt; i++) {
		if (is_prime(i)) {
			found++;
			printf("%llu\n", i);
		}
	}
}

int main(const int argc, char** const argv) {
	if (argc != 2) {
		fprintf(stderr, "Invalid number of parameters: %d\n", argc);
		return 1;
	}

	char* end;
	const uint64_t cnt = strtoull(argv[1], &end, 10);
	if (*end) {
		fprintf(stderr, "Invalid numeric parameter: %s\n", argv[1]);
		return 1;
	}

	find_primes(cnt);

	return 0;
}
