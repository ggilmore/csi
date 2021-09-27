section .text
global pangram
pangram:
    xor rsi, rsi ; seen = 0000...
    xor rax, rax ; result = 0
    xor rdx, rdx ; i = 0

    jmp test

loop:
    movzx rcx, byte [rdi + rdx] ; c = *rdi + i
    cmp rcx, 97 ; c >= 'a'?
    jge is_char

    add rcx, 32 ; c += 32

is_char:
    cmp rcx, 122; c > 'z'?
    jg increment
    cmp rcx, 97; c < 'a'?
    jl increment

    sub rcx, 97 ; what letter of the alphabet is c?
    bts rsi, rcx

    ; xor r8, r8   ; QUESTION - Why doesn't this work, but bts does?
    ; inc r8 ;     ;
    ; shl r8, cl   ; Set 1 in the correct place for the letter (where c is in rcx)
    ; add rsi, r8

increment:
    inc rdx

test:
    movzx rcx, byte [rdi + rdx]; c = *rdi + i
    cmp rcx, 0x00              ; c != '\0' ?
    jne loop

    cmp rsi, 0x3ffffff  ; have we seen all the characters? (this is 26 in decimal - one for each letter in the alphabet)

    jne done

success:
    inc rax

done:
	ret
