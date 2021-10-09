	.file	"benchmark.c"
	.text
	.globl	matrix_alloc
	.type	matrix_alloc, @function
matrix_alloc:
.LFB6:
	.cfi_startproc
	endbr64
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	pushq	%rbx
	subq	$40, %rsp
	.cfi_offset 3, -24
	movl	%edi, -36(%rbp)
	movl	%esi, -40(%rbp)
	movl	-36(%rbp), %eax
	cltq
	salq	$3, %rax
	movq	%rax, %rdi
	call	malloc@PLT
	movq	%rax, -24(%rbp)
	movl	$0, -28(%rbp)
	jmp	.L2
.L3:
	movl	-40(%rbp), %eax
	cltq
	movl	-28(%rbp), %edx
	movslq	%edx, %rdx
	leaq	0(,%rdx,8), %rcx
	movq	-24(%rbp), %rdx
	leaq	(%rcx,%rdx), %rbx
	movl	$8, %esi
	movq	%rax, %rdi
	call	calloc@PLT
	movq	%rax, (%rbx)
	addl	$1, -28(%rbp)
.L2:
	movl	-28(%rbp), %eax
	cmpl	-36(%rbp), %eax
	jl	.L3
	movq	-24(%rbp), %rax
	addq	$40, %rsp
	popq	%rbx
	popq	%rbp
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE6:
	.size	matrix_alloc, .-matrix_alloc
	.globl	matrix_free
	.type	matrix_free, @function
matrix_free:
.LFB7:
	.cfi_startproc
	endbr64
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	subq	$32, %rsp
	movq	%rdi, -24(%rbp)
	movl	%esi, -28(%rbp)
	movl	$0, -4(%rbp)
	jmp	.L6
.L7:
	movl	-4(%rbp), %eax
	cltq
	leaq	0(,%rax,8), %rdx
	movq	-24(%rbp), %rax
	addq	%rdx, %rax
	movq	(%rax), %rax
	movq	%rax, %rdi
	call	free@PLT
	addl	$1, -4(%rbp)
.L6:
	movl	-4(%rbp), %eax
	cmpl	-28(%rbp), %eax
	jl	.L7
	movq	-24(%rbp), %rax
	movq	%rax, %rdi
	call	free@PLT
	nop
	leave
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE7:
	.size	matrix_free, .-matrix_free
	.globl	matrix_fill_random
	.type	matrix_fill_random, @function
matrix_fill_random:
.LFB8:
	.cfi_startproc
	endbr64
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	subq	$32, %rsp
	movq	%rdi, -24(%rbp)
	movl	%esi, -28(%rbp)
	movl	%edx, -32(%rbp)
	movl	$0, -8(%rbp)
	jmp	.L9
.L12:
	movl	$0, -4(%rbp)
	jmp	.L10
.L11:
	call	rand@PLT
	cvtsi2sdl	%eax, %xmm0
	movl	-8(%rbp), %eax
	cltq
	leaq	0(,%rax,8), %rdx
	movq	-24(%rbp), %rax
	addq	%rdx, %rax
	movq	(%rax), %rax
	movl	-4(%rbp), %edx
	movslq	%edx, %rdx
	salq	$3, %rdx
	addq	%rdx, %rax
	movsd	.LC0(%rip), %xmm1
	divsd	%xmm1, %xmm0
	movsd	%xmm0, (%rax)
	addl	$1, -4(%rbp)
.L10:
	movl	-4(%rbp), %eax
	cmpl	-32(%rbp), %eax
	jl	.L11
	addl	$1, -8(%rbp)
.L9:
	movl	-8(%rbp), %eax
	cmpl	-28(%rbp), %eax
	jl	.L12
	nop
	nop
	leave
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE8:
	.size	matrix_fill_random, .-matrix_fill_random
	.globl	matrix_equal
	.type	matrix_equal, @function
matrix_equal:
.LFB9:
	.cfi_startproc
	endbr64
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	movq	%rdi, -24(%rbp)
	movq	%rsi, -32(%rbp)
	movl	%edx, -36(%rbp)
	movl	%ecx, -40(%rbp)
	movl	$0, -8(%rbp)
	jmp	.L14
.L20:
	movl	$0, -4(%rbp)
	jmp	.L15
.L19:
	movl	-8(%rbp), %eax
	cltq
	leaq	0(,%rax,8), %rdx
	movq	-24(%rbp), %rax
	addq	%rdx, %rax
	movq	(%rax), %rax
	movl	-4(%rbp), %edx
	movslq	%edx, %rdx
	salq	$3, %rdx
	addq	%rdx, %rax
	movsd	(%rax), %xmm0
	movl	-8(%rbp), %eax
	cltq
	leaq	0(,%rax,8), %rdx
	movq	-32(%rbp), %rax
	addq	%rdx, %rax
	movq	(%rax), %rax
	movl	-4(%rbp), %edx
	movslq	%edx, %rdx
	salq	$3, %rdx
	addq	%rdx, %rax
	movsd	(%rax), %xmm1
	ucomisd	%xmm1, %xmm0
	jp	.L21
	ucomisd	%xmm1, %xmm0
	je	.L22
.L21:
	movl	$0, %eax
	jmp	.L18
.L22:
	addl	$1, -4(%rbp)
.L15:
	movl	-4(%rbp), %eax
	cmpl	-40(%rbp), %eax
	jl	.L19
	addl	$1, -8(%rbp)
.L14:
	movl	-8(%rbp), %eax
	cmpl	-36(%rbp), %eax
	jl	.L20
	movl	$1, %eax
.L18:
	popq	%rbp
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE9:
	.size	matrix_equal, .-matrix_equal
	.globl	flush_cache
	.type	flush_cache, @function
flush_cache:
.LFB10:
	.cfi_startproc
	endbr64
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	subq	$16, %rsp
	movl	$4194304, -12(%rbp)
	movl	-12(%rbp), %eax
	cltq
	movq	%rax, %rdi
	call	malloc@PLT
	movq	%rax, -8(%rbp)
	movl	$0, -16(%rbp)
	jmp	.L24
.L25:
	movl	-16(%rbp), %eax
	movslq	%eax, %rdx
	movq	-8(%rbp), %rax
	addq	%rdx, %rax
	movl	-16(%rbp), %edx
	movb	%dl, (%rax)
	addl	$1, -16(%rbp)
.L24:
	movl	-16(%rbp), %eax
	cmpl	-12(%rbp), %eax
	jl	.L25
	movq	-8(%rbp), %rax
	movq	%rax, %rdi
	call	free@PLT
	nop
	leave
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE10:
	.size	flush_cache, .-flush_cache
	.section	.rodata
.LC1:
	.string	"Usage: ./benchmark [n]"
	.align 8
.LC3:
	.string	"Naive: %.3fs\nFast: %.3fs\n%0.2fx speedup\n"
	.align 8
.LC4:
	.string	"\nHowever, matrix results did not match!"
	.text
	.globl	main
	.type	main, @function
main:
.LFB11:
	.cfi_startproc
	endbr64
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	subq	$96, %rsp
	movl	%edi, -84(%rbp)
	movq	%rsi, -96(%rbp)
	cmpl	$2, -84(%rbp)
	je	.L27
	leaq	.LC1(%rip), %rdi
	call	puts@PLT
	movl	$1, %edi
	call	exit@PLT
.L27:
	movq	-96(%rbp), %rax
	addq	$8, %rax
	movq	(%rax), %rax
	movq	%rax, %rdi
	call	atoi@PLT
	movl	%eax, -68(%rbp)
	movl	-68(%rbp), %edx
	movl	-68(%rbp), %eax
	movl	%edx, %esi
	movl	%eax, %edi
	call	matrix_alloc
	movq	%rax, -64(%rbp)
	movl	-68(%rbp), %edx
	movl	-68(%rbp), %eax
	movl	%edx, %esi
	movl	%eax, %edi
	call	matrix_alloc
	movq	%rax, -56(%rbp)
	movl	-68(%rbp), %edx
	movl	-68(%rbp), %eax
	movl	%edx, %esi
	movl	%eax, %edi
	call	matrix_alloc
	movq	%rax, -48(%rbp)
	movl	-68(%rbp), %edx
	movl	-68(%rbp), %eax
	movl	%edx, %esi
	movl	%eax, %edi
	call	matrix_alloc
	movq	%rax, -40(%rbp)
	movl	-68(%rbp), %edx
	movl	-68(%rbp), %ecx
	movq	-64(%rbp), %rax
	movl	%ecx, %esi
	movq	%rax, %rdi
	call	matrix_fill_random
	movl	-68(%rbp), %edx
	movl	-68(%rbp), %ecx
	movq	-56(%rbp), %rax
	movl	%ecx, %esi
	movq	%rax, %rdi
	call	matrix_fill_random
	movl	$0, %eax
	call	flush_cache
	call	clock@PLT
	movq	%rax, -32(%rbp)
	movl	-68(%rbp), %r8d
	movl	-68(%rbp), %edi
	movl	-68(%rbp), %ecx
	movq	-56(%rbp), %rdx
	movq	-64(%rbp), %rsi
	movq	-48(%rbp), %rax
	movl	%r8d, %r9d
	movl	%edi, %r8d
	movq	%rax, %rdi
	call	matrix_multiply@PLT
	call	clock@PLT
	movq	%rax, -24(%rbp)
	movq	-24(%rbp), %rax
	subq	-32(%rbp), %rax
	cvtsi2sdq	%rax, %xmm0
	movsd	.LC2(%rip), %xmm1
	divsd	%xmm1, %xmm0
	movsd	%xmm0, -16(%rbp)
	movl	$0, %eax
	call	flush_cache
	call	clock@PLT
	movq	%rax, -32(%rbp)
	movl	-68(%rbp), %r8d
	movl	-68(%rbp), %edi
	movl	-68(%rbp), %ecx
	movq	-56(%rbp), %rdx
	movq	-64(%rbp), %rsi
	movq	-40(%rbp), %rax
	movl	%r8d, %r9d
	movl	%edi, %r8d
	movq	%rax, %rdi
	call	fast_matrix_multiply@PLT
	call	clock@PLT
	movq	%rax, -24(%rbp)
	movq	-24(%rbp), %rax
	subq	-32(%rbp), %rax
	cvtsi2sdq	%rax, %xmm0
	movsd	.LC2(%rip), %xmm1
	divsd	%xmm1, %xmm0
	movsd	%xmm0, -8(%rbp)
	movsd	-16(%rbp), %xmm0
	movapd	%xmm0, %xmm1
	divsd	-8(%rbp), %xmm1
	movsd	-8(%rbp), %xmm0
	movq	-16(%rbp), %rax
	movapd	%xmm1, %xmm2
	movapd	%xmm0, %xmm1
	movq	%rax, %xmm0
	leaq	.LC3(%rip), %rdi
	movl	$3, %eax
	call	printf@PLT
	movl	-68(%rbp), %ecx
	movl	-68(%rbp), %edx
	movq	-40(%rbp), %rsi
	movq	-48(%rbp), %rax
	movq	%rax, %rdi
	call	matrix_equal
	xorl	$1, %eax
	testb	%al, %al
	je	.L28
	leaq	.LC4(%rip), %rdi
	call	puts@PLT
.L28:
	movl	-68(%rbp), %edx
	movq	-64(%rbp), %rax
	movl	%edx, %esi
	movq	%rax, %rdi
	call	matrix_free
	movl	-68(%rbp), %edx
	movq	-56(%rbp), %rax
	movl	%edx, %esi
	movq	%rax, %rdi
	call	matrix_free
	movl	-68(%rbp), %edx
	movq	-48(%rbp), %rax
	movl	%edx, %esi
	movq	%rax, %rdi
	call	matrix_free
	movl	-68(%rbp), %edx
	movq	-40(%rbp), %rax
	movl	%edx, %esi
	movq	%rax, %rdi
	call	matrix_free
	movl	$0, %eax
	leave
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE11:
	.size	main, .-main
	.section	.rodata
	.align 8
.LC0:
	.long	4290772992
	.long	1105199103
	.align 8
.LC2:
	.long	0
	.long	1093567616
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
