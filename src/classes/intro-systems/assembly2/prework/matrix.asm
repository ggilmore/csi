section .text
global index
index:
	; rdi: matrix
	; rsi: rows
	; rdx: cols
	; rcx: rindex
	; r8: cindex

    mov rax, rdx
    imul rax, rcx ; cols * rindex
    add rax, r8 ; cols * rindex + cindex = offset

    mov rax, [rdi + rax * 4]; offset * 4 bytes per int
	ret
