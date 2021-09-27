section .text
global binary_convert

binary_convert:
    xor rax, rax               ; result = 0
    xor rsi, rsi               ; i = 0

    jmp test

grow_number:
    shl rax, 1                 ; result *= 2

loop:
    movzx rdx, byte [rdi + rsi]; c = *rdi + i - QUESTION: "I notice that I have to specify
                               ; that I only want a single 'byte' here. An ascii character can be specified in a single byte.
                               ; Why is there extra junk when I move a 'word' or a 'qword' into rdx that makes the comparsion fail?"
    cmp rdx, 0x31              ; c != '1' ?

    jne increment

    add rax, 0x1               ; set lowest order bit to 1

increment:
    inc rsi                    ; i++

test:
    movzx rdx, byte [rdi + rsi]; c = *rdi + i
    cmp rdx, 0x00              ; c != '\0' ?
    jne grow_number

    ret
