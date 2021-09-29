section .text
global fib
fib:
    push  rbx     ; backup 'rbx'
    push  rbp     ; backup 'rbp'

    mov   ebx, edi; push 'n' to saved ebx register

    cmp ebx, 1    ; n <= 1?
    jle done

    lea edi, [ ebx - 1 ] ; prepare 'n -1' arg for first rec call
    call fib ; fib(n-1)
    mov ebp, eax ; save result fib(n-1)

    lea edi, [ ebx - 2 ];  prepare 'n - 2' arg for second rec call
    call fib ; fib(n-2)
    mov ebx, eax; store fib(n-2) answer

    add ebx, ebp; fib(n-2) answer + fib(n-1) answer

done:
    mov eax, ebx    ; put answer in eax

    pop rbp ; restore 'rbp'
    pop rbx ; restore 'rbx'

	ret
