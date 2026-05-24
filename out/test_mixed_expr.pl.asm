; Generated_Assembly

STACK SEGMENT STACK
    DW 100 DUP(?)
STACK ENDS

DATA SEGMENT
    ; user variables
x DW ?
r1 DW ?
C DW 10
r2 DW ?
y DW ?

    ; compiler-generated temporaries
T0 DW ?
T1 DW ?
T2 DW ?
T3 DW ?
T4 DW ?

    ; string constants
str0 DB 'r1 = ', '$'
str1 DB 'r2 = ', '$'
DATA ENDS

CODE SEGMENT
    ASSUME CS:CODE, DS:DATA, SS:STACK

START:
    MOV AX, DATA
    MOV DS, AX

L0:
    MOV x, 5

L1:
    MOV y, 3

L2:
    MOV AX, 5
    MOV CX, C
    IMUL CX
    MOV T0, AX

L3:
    MOV AX, T0
    ADD AX, 3
    MOV T1, AX

L4:
    MOV AX, T1
    MOV r1, AX

L5:
    MOV T2, 8

L6:
    MOV T3, 2

L7:
    MOV AX, T2
    MOV CX, T3
    IMUL CX
    MOV T4, AX

L8:
    MOV AX, T4
    MOV r2, AX

L9:
    LEA DX, str0
    MOV AH, 09h
    INT 21h

L10:
    MOV AX, T1
    CALL PRINT_INT

L11:
    LEA DX, str1
    MOV AH, 09h
    INT 21h

L12:
    MOV AX, T4
    CALL PRINT_INT

    MOV AH, 4Ch
    INT 21h

; Print integer in AX
PRINT_INT PROC
    PUSH AX
    PUSH BX
    PUSH CX
    PUSH DX

    CMP AX, 0
    JGE print_pos
    PUSH AX
    MOV DL, '-'
    MOV AH, 02h
    INT 21h
    POP AX
    NEG AX
print_pos:
    MOV CX, 0
    MOV BX, 10
print_div:
    MOV DX, 0
    DIV BX
    PUSH DX
    INC CX
    CMP AX, 0
    JNE print_div
print_digits:
    POP DX
    ADD DL, '0'
    MOV AH, 02h
    INT 21h
    LOOP print_digits

    POP DX
    POP CX
    POP BX
    POP AX
    RET
PRINT_INT ENDP
CODE ENDS
END START
