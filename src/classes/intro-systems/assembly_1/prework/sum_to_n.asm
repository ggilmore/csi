; section .text
; global sum_to_n, loop, done
; sum_to_n:
;     ; rdi = n
;     xor rax, rax ; total = 0
;     xor rsi, rsi ; i = 0 (where do I save local vars)?

;     cmp rsi, rdi ; i > n ?
;     jg done

; loop:
;     add rax, rsi ; total = total + i;
;     inc rsi      ; i++
;     cmp rsi, rdi ; i <= n ?
;     jle loop

; done:
; 	ret

section .text
global sum_to_n
sum_to_n:
    ; rdi = n
    mov rax, rdi
    inc rax
    imul rdi
    sar rax, 1
    ret
