; Generated_Assembly

STACK SEGMENT STACK
    DW 100 DUP(?)
STACK ENDS

DATA SEGMENT
    ; user variables
Tab DW 10 DUP(?)
i DW ?
sum DW ?

    ; compiler-generated temporaries
T0 DW ?
T1 DW ?
T2 DW ?
T3 DW ?
T4 DW ?

    ; string constants
str0 DB 'sum = ', '$'
DATA ENDS

CODE SEGMENT
    ASSUME CS:CODE, DS:DATA, SS:STACK

START:
    MOV AX, DATA
    MOV DS, AX

L0:
    MOV i, 0

L1:
    MOV sum, 0

L2:
    MOV AX, i
    CMP AX, 5
    MOV AX, 0
    JGE S2_end
    MOV AX, 1
S2_end:
    MOV T0, AX

L3:
    MOV AX, T0
    CMP AX, 0
    JZ L12

L4:
    MOV AX, i
    MOV CX, 2
    IMUL CX
    MOV T1, AX

L5:
    MOV SI, i
    ADD SI, SI
    MOV AX, T1
    MOV Tab[SI], AX

L6:
    MOV SI, i
    ADD SI, SI
    MOV AX, Tab[SI]
    MOV T2, AX

L7:
    MOV AX, sum
    MOV CX, T2
    ADD AX, CX
    MOV T3, AX

L8:
    MOV AX, T3
    MOV sum, AX

L9:
    MOV AX, i
    ADD AX, 1
    MOV T4, AX

L10:
    MOV AX, T4
    MOV i, AX

L11:
    JMP L2

L12:
    LEA DX, str0
    MOV AH, 09h
    INT 21h

L13:
    MOV AX, sum
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
