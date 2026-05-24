; Generated_Assembly

STACK SEGMENT STACK
    DW 100 DUP(?)
STACK ENDS

DATA SEGMENT
    ; user variables
i DW ?
count DW ?
j DW ?

    ; compiler-generated temporaries
T0 DW ?
T1 DW ?
T2 DW ?
T3 DW ?
T4 DW ?

    ; string constants
str0 DB 'count = ', '$'
DATA ENDS

CODE SEGMENT
    ASSUME CS:CODE, DS:DATA, SS:STACK

START:
    MOV AX, DATA
    MOV DS, AX

L0:
    MOV count, 0

L1:
    MOV i, 0

L2:
    MOV AX, i
    CMP AX, 3
    MOV AX, 0
    JGE S2_end
    MOV AX, 1
S2_end:
    MOV T0, AX

L3:
    MOV AX, T0
    CMP AX, 0
    JZ L15

L4:
    MOV j, 0

L5:
    MOV AX, j
    CMP AX, 2
    MOV AX, 0
    JGE S5_end
    MOV AX, 1
S5_end:
    MOV T1, AX

L6:
    MOV AX, T1
    CMP AX, 0
    JZ L12

L7:
    MOV AX, count
    ADD AX, 1
    MOV T2, AX

L8:
    MOV AX, T2
    MOV count, AX

L9:
    MOV AX, j
    ADD AX, 1
    MOV T3, AX

L10:
    MOV AX, T3
    MOV j, AX

L11:
    JMP L5

L12:
    MOV AX, i
    ADD AX, 1
    MOV T4, AX

L13:
    MOV AX, T4
    MOV i, AX

L14:
    JMP L2

L15:
    LEA DX, str0
    MOV AH, 09h
    INT 21h

L16:
    MOV AX, count
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
