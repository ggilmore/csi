	.file	"pagecount_orig.c"
	.intel_syntax noprefix
	.text
	.globl	pagecount
	.type	pagecount, @function
pagecount:
.LFB23:
	.cfi_startproc
	endbr64
	mov	rax, rdi
	mov	edx, 0
	div	rsi
	ret
	.cfi_endproc
.LFE23:
	.size	pagecount, .-pagecount
	.section	.rodata.str1.8,"aMS",@progbits,1
	.align 8
.LC3:
	.string	"%.2fs to run %d tests (%.2fns per test)\n"
	.text
	.globl	main
	.type	main, @function
main:
.LFB24:
	.cfi_startproc
	endbr64
	push	r13
	.cfi_def_cfa_offset 16
	.cfi_offset 13, -16
	push	r12
	.cfi_def_cfa_offset 24
	.cfi_offset 12, -24
	push	rbp
	.cfi_def_cfa_offset 32
	.cfi_offset 6, -32
	push	rbx
	.cfi_def_cfa_offset 40
	.cfi_offset 3, -40
	sub	rsp, 72
	.cfi_def_cfa_offset 112
	mov	rax, QWORD PTR fs:40
	mov	QWORD PTR 56[rsp], rax
	xor	eax, eax
	movabs	rax, 4294967296
	mov	QWORD PTR [rsp], rax
	movabs	rsi, 1099511627776
	mov	QWORD PTR 8[rsp], rsi
	movabs	rbx, 4503599627370496
	mov	QWORD PTR 16[rsp], rbx
	mov	QWORD PTR 32[rsp], 4096
	mov	QWORD PTR 40[rsp], 65536
	mov	QWORD PTR 48[rsp], rax
	call	clock@PLT
	mov	r12, rax
	mov	ebx, 0
	mov	edx, 0
.L3:
	movsx	rax, edx
	imul	rax, rax, 1431655766
	shr	rax, 32
	mov	ecx, edx
	sar	ecx, 31
	sub	eax, ecx
	lea	eax, [rax+rax*2]
	mov	edi, edx
	sub	edi, eax
	mov	eax, edi
	cdqe
	add	ebx, 1
	add	ebx, DWORD PTR [rsp+rax*8]
	add	ebx, DWORD PTR 32[rsp+rax*8]
	add	edx, 1
	cmp	edx, 222222222
	jne	.L3
	call	clock@PLT
	mov	rbp, rax
	call	clock@PLT
	mov	r13, rax
	mov	ecx, 0
.L4:
	movsx	rax, ecx
	imul	rax, rax, 1431655766
	shr	rax, 32
	mov	edx, ecx
	sar	edx, 31
	sub	eax, edx
	lea	eax, [rax+rax*2]
	mov	edi, ecx
	sub	edi, eax
	mov	eax, edi
	cdqe
	mov	rsi, QWORD PTR [rsp+rax*8]
	mov	rdi, QWORD PTR 32[rsp+rax*8]
	mov	rax, rsi
	mov	edx, 0
	div	rdi
	add	rsi, rdi
	add	rsi, rax
	add	ebx, esi
	add	ecx, 1
	cmp	ecx, 222222222
	jne	.L4
	call	clock@PLT
	sub	rax, r13
	sub	rbp, r12
	sub	rax, rbp
	pxor	xmm0, xmm0
	cvtsi2sd	xmm0, rax
	divsd	xmm0, QWORD PTR .LC0[rip]
	movapd	xmm1, xmm0
	mulsd	xmm1, QWORD PTR .LC1[rip]
	divsd	xmm1, QWORD PTR .LC2[rip]
	mov	edx, 222222222
	lea	rsi, .LC3[rip]
	mov	edi, 1
	mov	eax, 2
	call	__printf_chk@PLT
	mov	rax, QWORD PTR 56[rsp]
	xor	rax, QWORD PTR fs:40
	jne	.L9
	mov	eax, ebx
	add	rsp, 72
	.cfi_remember_state
	.cfi_def_cfa_offset 40
	pop	rbx
	.cfi_def_cfa_offset 32
	pop	rbp
	.cfi_def_cfa_offset 24
	pop	r12
	.cfi_def_cfa_offset 16
	pop	r13
	.cfi_def_cfa_offset 8
	ret
.L9:
	.cfi_restore_state
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE24:
	.size	main, .-main
	.section	.rodata.cst8,"aM",@progbits,8
	.align 8
.LC0:
	.long	0
	.long	1093567616
	.align 8
.LC1:
	.long	0
	.long	1104006501
	.align 8
.LC2:
	.long	469762048
	.long	1101692335
	.ident	"GCC: (Ubuntu 9.3.0-17ubuntu1~20.04) 9.3.0"
	.section	.note.GNU-stack,"",@progbits
	.section	.note.gnu.property,"a"
	.align 8
	.long	 1f - 0f
	.long	 4f - 1f
	.long	 5
0:
	.string	 "GNU"
1:
	.align 8
	.long	 0xc0000002
	.long	 3f - 2f
2:
	.long	 0x3
3:
	.align 8
4:
