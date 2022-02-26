OUTPUT_X64=primes_x86_64
OUTPUT_RISC_V=primes_risc_v

build: $(OUTPUT_X64) $(OUTPUT_RISC_V)

$(OUTPUT_X64): *.c
	gcc -o $@ $^

$(OUTPUT_RISC_V): *.c
	../../riscv_gcc/run.sh riscv64-unknown-linux-gnu-gcc -o $@ -misa-spec=2.2 -mabi=lp64 -march=rv64im $^

clean:
	rm -f $(OUTPUT_X64) $(OUTPUT_RISC_V)
