
dotproduct.o:	file format mach-o 64-bit x86-64


Disassembly of section __TEXT,__text:

0000000000000000 <_dotproduct>:
; {
       0: 55                           	push	rbp
       1: 48 89 e5                     	mov	rbp, rsp
       4: 48 89 7d f8                  	mov	qword ptr [rbp - 8], rdi
       8: 48 89 75 f0                  	mov	qword ptr [rbp - 16], rsi
;   data_t sum = 0, sum1 = 0, sum2 = 0, u_val, v_val, u_val1, v_val1, u_val2,
       c: 48 c7 45 e8 00 00 00 00      	mov	qword ptr [rbp - 24], 0
      14: 48 c7 45 e0 00 00 00 00      	mov	qword ptr [rbp - 32], 0
      1c: 48 c7 45 d8 00 00 00 00      	mov	qword ptr [rbp - 40], 0
;   int length = u->len;
      24: 48 8b 45 f8                  	mov	rax, qword ptr [rbp - 8]
      28: 48 8b 00                     	mov	rax, qword ptr [rax]
      2b: 89 45 a4                     	mov	dword ptr [rbp - 92], eax
;   for (i = 0; i < length;
      2e: 48 c7 45 98 00 00 00 00      	mov	qword ptr [rbp - 104], 0
      36: 48 8b 45 98                  	mov	rax, qword ptr [rbp - 104]
      3a: 48 63 4d a4                  	movsxd	rcx, dword ptr [rbp - 92]
      3e: 48 39 c8                     	cmp	rax, rcx
      41: 0f 8d c0 00 00 00            	jge	0x107 <_dotproduct+0x107>
;     u_val = u->data[i];
      47: 48 8b 45 f8                  	mov	rax, qword ptr [rbp - 8]
      4b: 48 8b 40 08                  	mov	rax, qword ptr [rax + 8]
      4f: 48 8b 4d 98                  	mov	rcx, qword ptr [rbp - 104]
      53: 48 8b 04 c8                  	mov	rax, qword ptr [rax + 8*rcx]
      57: 48 89 45 d0                  	mov	qword ptr [rbp - 48], rax
;     v_val = v->data[i];
      5b: 48 8b 45 f0                  	mov	rax, qword ptr [rbp - 16]
      5f: 48 8b 40 08                  	mov	rax, qword ptr [rax + 8]
      63: 48 8b 4d 98                  	mov	rcx, qword ptr [rbp - 104]
      67: 48 8b 04 c8                  	mov	rax, qword ptr [rax + 8*rcx]
      6b: 48 89 45 c8                  	mov	qword ptr [rbp - 56], rax
;     u_val1 = u->data[i + 1];
      6f: 48 8b 45 f8                  	mov	rax, qword ptr [rbp - 8]
      73: 48 8b 40 08                  	mov	rax, qword ptr [rax + 8]
      77: 48 8b 4d 98                  	mov	rcx, qword ptr [rbp - 104]
      7b: 48 8b 44 c8 08               	mov	rax, qword ptr [rax + 8*rcx + 8]
      80: 48 89 45 c0                  	mov	qword ptr [rbp - 64], rax
;     v_val1 = v->data[i + 1];
      84: 48 8b 45 f0                  	mov	rax, qword ptr [rbp - 16]
      88: 48 8b 40 08                  	mov	rax, qword ptr [rax + 8]
      8c: 48 8b 4d 98                  	mov	rcx, qword ptr [rbp - 104]
      90: 48 8b 44 c8 08               	mov	rax, qword ptr [rax + 8*rcx + 8]
      95: 48 89 45 b8                  	mov	qword ptr [rbp - 72], rax
;     u_val2 = u->data[i + 2];
      99: 48 8b 45 f8                  	mov	rax, qword ptr [rbp - 8]
      9d: 48 8b 40 08                  	mov	rax, qword ptr [rax + 8]
      a1: 48 8b 4d 98                  	mov	rcx, qword ptr [rbp - 104]
      a5: 48 8b 44 c8 10               	mov	rax, qword ptr [rax + 8*rcx + 16]
      aa: 48 89 45 b0                  	mov	qword ptr [rbp - 80], rax
;     v_val2 = v->data[i + 2];
      ae: 48 8b 45 f0                  	mov	rax, qword ptr [rbp - 16]
      b2: 48 8b 40 08                  	mov	rax, qword ptr [rax + 8]
      b6: 48 8b 4d 98                  	mov	rcx, qword ptr [rbp - 104]
      ba: 48 8b 44 c8 10               	mov	rax, qword ptr [rax + 8*rcx + 16]
      bf: 48 89 45 a8                  	mov	qword ptr [rbp - 88], rax
;     sum += u_val * v_val;
      c3: 48 8b 45 d0                  	mov	rax, qword ptr [rbp - 48]
      c7: 48 0f af 45 c8               	imul	rax, qword ptr [rbp - 56]
      cc: 48 03 45 e8                  	add	rax, qword ptr [rbp - 24]
      d0: 48 89 45 e8                  	mov	qword ptr [rbp - 24], rax
;     sum1 += u_val1 * v_val1;
      d4: 48 8b 45 c0                  	mov	rax, qword ptr [rbp - 64]
      d8: 48 0f af 45 b8               	imul	rax, qword ptr [rbp - 72]
      dd: 48 03 45 e0                  	add	rax, qword ptr [rbp - 32]
      e1: 48 89 45 e0                  	mov	qword ptr [rbp - 32], rax
;     sum2 += u_val2 * v_val2;
      e5: 48 8b 45 b0                  	mov	rax, qword ptr [rbp - 80]
      e9: 48 0f af 45 a8               	imul	rax, qword ptr [rbp - 88]
      ee: 48 03 45 d8                  	add	rax, qword ptr [rbp - 40]
      f2: 48 89 45 d8                  	mov	qword ptr [rbp - 40], rax
;        i += 3) { // we can assume both vectors are same length
      f6: 48 8b 45 98                  	mov	rax, qword ptr [rbp - 104]
      fa: 48 83 c0 03                  	add	rax, 3
      fe: 48 89 45 98                  	mov	qword ptr [rbp - 104], rax
;   for (i = 0; i < length;
     102: e9 2f ff ff ff               	jmp	0x36 <_dotproduct+0x36>
;   for (; i < length; i++) {
     107: e9 00 00 00 00               	jmp	0x10c <_dotproduct+0x10c>
     10c: 48 8b 45 98                  	mov	rax, qword ptr [rbp - 104]
     110: 48 63 4d a4                  	movsxd	rcx, dword ptr [rbp - 92]
     114: 48 39 c8                     	cmp	rax, rcx
     117: 0f 8d 4a 00 00 00            	jge	0x167 <_dotproduct+0x167>
;     u_val = u->data[i];
     11d: 48 8b 45 f8                  	mov	rax, qword ptr [rbp - 8]
     121: 48 8b 40 08                  	mov	rax, qword ptr [rax + 8]
     125: 48 8b 4d 98                  	mov	rcx, qword ptr [rbp - 104]
     129: 48 8b 04 c8                  	mov	rax, qword ptr [rax + 8*rcx]
     12d: 48 89 45 d0                  	mov	qword ptr [rbp - 48], rax
;     v_val = v->data[i];
     131: 48 8b 45 f0                  	mov	rax, qword ptr [rbp - 16]
     135: 48 8b 40 08                  	mov	rax, qword ptr [rax + 8]
     139: 48 8b 4d 98                  	mov	rcx, qword ptr [rbp - 104]
     13d: 48 8b 04 c8                  	mov	rax, qword ptr [rax + 8*rcx]
     141: 48 89 45 c8                  	mov	qword ptr [rbp - 56], rax
;     sum += u_val * u_val;
     145: 48 8b 45 d0                  	mov	rax, qword ptr [rbp - 48]
     149: 48 0f af 45 d0               	imul	rax, qword ptr [rbp - 48]
     14e: 48 03 45 e8                  	add	rax, qword ptr [rbp - 24]
     152: 48 89 45 e8                  	mov	qword ptr [rbp - 24], rax
;   for (; i < length; i++) {
     156: 48 8b 45 98                  	mov	rax, qword ptr [rbp - 104]
     15a: 48 83 c0 01                  	add	rax, 1
     15e: 48 89 45 98                  	mov	qword ptr [rbp - 104], rax
     162: e9 a5 ff ff ff               	jmp	0x10c <_dotproduct+0x10c>
;   return sum + sum1 + sum2;
     167: 48 8b 45 e8                  	mov	rax, qword ptr [rbp - 24]
     16b: 48 03 45 e0                  	add	rax, qword ptr [rbp - 32]
     16f: 48 03 45 d8                  	add	rax, qword ptr [rbp - 40]
     173: 5d                           	pop	rbp
     174: c3                           	ret
