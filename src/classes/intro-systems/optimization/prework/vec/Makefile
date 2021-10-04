
CFLAGS = -g -Wall -Wextra -pedantic -Werror

ifeq ($(shell uname -s),Darwin)
	ALL_LDFLAGS = -Wl,-pie -Wl,-fatal_warnings
	ALL_ASFLAGS = -f macho64 --prefix _
else
	ALL_LDFLAGS = -pie -Wl,--fatal-warnings
	ALL_ASFLAGS = -f elf64
endif

ALL_CFLAGS = -std=c99 -fPIE -m64 $(CFLAGS)

C_OBJS = $(patsubst %.c,%.o,$(wildcard *.c))
ALL_OBJS = $(C_OBJS) vendor/unity.o

CC_CMD = $(CC) $(ALL_CFLAGS) -g3 -c -o $@ $<

all: tests
	@./$<

tests: $(ALL_OBJS)
	@$(CC) $(ALL_CFLAGS) $(ALL_LDFLAGS) -o $@ $(ALL_OBJS)

assembly: $(ALL_OBJS)
    @objdump -x86-asm-syntax=intel -S dotproduct.o > ./dotproduct.s

%.o: %.c
	@$(CC_CMD)

vendor/unity.o: vendor/unity.c vendor/unity.h vendor/unity_internals.h
	@$(CC_CMD)

clean:
	@rm -f *.o vendor/*.o tests

.PHONY: assembly all clean
