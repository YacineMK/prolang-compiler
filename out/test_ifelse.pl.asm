; Generated_Assembly

STACK SEGMENT STACK
    DW 100 DUP(?)
STACK ENDS

DATA SEGMENT
    ; user variables
a DW ?
b DW ?
result DW ?

    ; compiler-generated temporaries
T0 DW ?
T1 DW ?
T2 DW ?
T3 DW ?

    ; string constants
str0 DB 'result = ', '$'
str1 DB 'result2 = ', '$'
DATA ENDS

CODE SEGMENT
    ASSUME CS:CODE, DS:DATA, SS:STACK

START:
    MOV AX, DATA
    MOV DS, AX

L0:
    MOV a, 15

L1:
    MOV b, 10

L2:
    MOV T0, 1

L3:
    MOV AX, T0
    CMP AX, 0
    JZ L7

L4:
    MOV T1, 5

L5:
    MOV AX, T1
    MOV result, AX

L6:
    JMP L9

L7:
    MOV T2, -5

L8:
    MOV AX, T2
    MOV result, AX

L9:
    LEA DX, str0
    MOV AH, 09h
    INT 21h

L10:
    MOV AX, result
    CALL PRINT_INT

L11:
    MOV T3, 0

L12:
    MOV AX, T3
    CMP AX, 0
    JZ L15

L13:
    MOV result, 100

L14:
    JMP L16

L15:
    MOV result, 200

L16:
    LEA DX, str1
    MOV AH, 09h
    INT 21h

L17:
    MOV AX, result
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
