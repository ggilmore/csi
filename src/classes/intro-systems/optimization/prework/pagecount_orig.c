#include <math.h>
#include <stdint.h>
#include <stdio.h>
#include <time.h>

#define TEST_LOOPS 222222222

/*
Geoffrey's guess:

mov rax, rdi (move memory size to lower 64 bits)
xor rdx, rdx; 0 out rdx (high 64 bits)
div rsi; ([rdx]:[rax] <- memory size) / (rsi <- page size)

ret
*/

/*
cc -O0 -S -masm=intel pagecount.c

Geoffrey's comment: the "xor ecx, ecx" instruction is confusing
- why can't you get away with just "xor edx, edx" directly?

_pagecount:                             ## @pagecount
        .cfi_startproc
## %bb.0:
        push	rbp
        .cfi_def_cfa_offset 16
        .cfi_offset rbp, -16
        mov	rbp, rsp
        .cfi_def_cfa_register rbp
        mov	qword ptr [rbp - 8], rdi
        mov	qword ptr [rbp - 16], rsi
        mov	rax, qword ptr [rbp - 8]
        xor	ecx, ecx ; why is this needed? Can't you just xor edx, edx?
        mov	edx, ecx
        div	qword ptr [rbp - 16]
        pop	rbp
        ret
        .cfi_endproc
                                        ## -- End function
        .section	__TEXT,__literal8,8byte_literals
*/

/*
cc -O1 -S -masm=intel pagecount.c

Geoffrey's comment: minus me forgetting to push rbp,
I think this is pretty similar to mine

_pagecount:                             ## @pagecount
        .cfi_startproc
## %bb.0:
        push	rbp
        .cfi_def_cfa_offset 16
        .cfi_offset rbp, -16
        mov	rbp, rsp
        .cfi_def_cfa_register rbp
        mov	rax, rdi
        xor	edx, edx ; gcc uses a move here instead
        div	rsi
        pop	rbp
        ret
        .cfi_endproc
                                        ## -- End function
        .section	__TEXT,__literal16,16byte_literals
*/

/*
cc -O2 -S -masm=intel pagecount.c

Geoffrey's comment: overall, I think this is testing to see
whether or not it can get away with doing a int division (which
I would assume is much cheaper).

gcc at the same optimization level doesn't do anything different than -O1.

_pagecount:                             ## @pagecount
        .cfi_startproc
## %bb.0:
        push	rbp
        .cfi_def_cfa_offset 16
        .cfi_offset rbp, -16
        mov	rbp, rsp
        .cfi_def_cfa_register rbp
        mov	rax, rdi
        mov	rcx, rdi
        or	rcx, rsi
        shr	rcx, 32
        je	LBB0_1
## %bb.2:
        xor	edx, edx
        div	rsi
        pop	rbp
        ret
LBB0_1:
                                        ## kill: def $eax killed $eax killed
$rax xor	edx, edx div	esi
                                        ## kill: def $eax killed $eax def $rax
        pop	rbp
        ret
        .cfi_endproc
                                        ## -- End function
        .section	__TEXT,__literal16,16byte_literals
        .p2align	4                               ## -- Begin function
main

*/

uint64_t
pagecount(uint64_t memory_size, uint64_t page_size)
{
  return memory_size / page_size;
}

int
main(int argc, char** argv)
{
  clock_t baseline_start, baseline_end, test_start, test_end;
  uint64_t memory_size, page_size;
  double clocks_elapsed, time_elapsed;
  int i, ignore = 0;

  uint64_t msizes[] = { 1L << 32, 1L << 40, 1L << 52 };
  uint64_t psizes[] = { 1L << 12, 1L << 16, 1L << 32 };

  baseline_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) {
    memory_size = msizes[i % 3];
    page_size = psizes[i % 3];
    ignore += 1 + memory_size +
              page_size; // so that this loop isn't just optimized away
  }
  baseline_end = clock();

  test_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) {
    memory_size = msizes[i % 3];
    page_size = psizes[i % 3];
    ignore += pagecount(memory_size, page_size) + memory_size + page_size;
  }
  test_end = clock();

  clocks_elapsed = test_end - test_start - (baseline_end - baseline_start);
  time_elapsed = clocks_elapsed / CLOCKS_PER_SEC;

  printf("%.2fs to run %d tests (%.2fns per test)\n",
         time_elapsed,
         TEST_LOOPS,
         time_elapsed * 1e9 / TEST_LOOPS);
  return ignore;
}
