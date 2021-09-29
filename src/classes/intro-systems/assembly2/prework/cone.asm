default rel

section .text
global volume
volume: ; r in xmm0, h in xmm1
    mulss xmm0, xmm0 ; r^2
    mulss xmm0, [pi] ; pi*r^2
    mulss xmm0, xmm1 ; pi*r^2*h
    divss xmm0, [three]; (pi*r^2*h)/3
 	ret

section .data
pi: dd 3.14159
three: dd 3.0 ; why do I need to put "3.0" instead of just "3" here?
