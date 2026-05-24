; Generated_Assembly

STACK SEGMENT STACK
    DW 100 DUP(?)
STACK ENDS

DATA SEGMENT
    ; user variables
Tabfloat DW 30 DUP(?)
x DW ?
k DW ?
somme DW 0
Tabint DW 50 DUP(?)
Pi DW 3
a DW ?
i DW ?
b DW ?
j DW ?
y DW ?
z DW ?
Max DW 100
moyenne DW 0

    ; compiler-generated temporaries
T0 DW ?
T1 DW ?
T2 DW ?
T3 DW ?
T4 DW ?
T5 DW ?
T6 DW ?
T7 DW ?
T8 DW ?
T9 DW ?
T10 DW ?
T11 DW ?
T12 DW ?
T13 DW ?
T14 DW ?
T15 DW ?
T16 DW ?
T17 DW ?
T18 DW ?
T19 DW ?
T20 DW ?
T21 DW ?
T22 DW ?
T23 DW ?
T24 DW ?
T25 DW ?
T26 DW ?
T27 DW ?
T28 DW ?
T29 DW ?
T30 DW ?
T31 DW ?
T32 DW ?
T33 DW ?
T34 DW ?
T35 DW ?
T36 DW ?
T37 DW ?
T38 DW ?
T39 DW ?
T40 DW ?
T41 DW ?
T42 DW ?
T43 DW ?
T44 DW ?
T45 DW ?
T46 DW ?
T47 DW ?
T48 DW ?
T49 DW ?
T50 DW ?
T51 DW ?

    ; string constants
str0 DB 'Valeur finale de x: ', '$'
str1 DB 'Somme: ', '$'
str2 DB 'Moyenne: ', '$'
DATA ENDS

CODE SEGMENT
    ASSUME CS:CODE, DS:DATA, SS:STACK

START:
    MOV AX, DATA
    MOV DS, AX

L0:
    MOV x, 10

L1:
    MOV y, 5

L2:
    MOV z, 2

L3:
    MOV a, 2

L4:
    MOV AX, 2
    MOV CX, Pi
    ADD AX, CX
    MOV T0, AX

L5:
    MOV AX, T0
    MOV CX, 2
    IMUL CX
    MOV T1, AX

L6:
    MOV AX, T1
    MOV b, AX

L7:
    MOV T2, 10

L8:
    MOV AX, 10
    MOV CX, T2
    ADD AX, CX
    MOV T3, AX

L9:
    MOV SI, 0
    ADD SI, SI
    MOV AX, T3
    MOV Tabint[SI], AX

L10:
    MOV AX, T1
    ADD AX, 3
    MOV T4, AX

L11:
    MOV AX, T4
    MOV CX, 2
    CWD
    IDIV CX
    MOV T5, AX

L12:
    MOV SI, 1
    ADD SI, SI
    MOV AX, T5
    MOV Tabfloat[SI], AX

L13:
    MOV T6, 1

L14:
    MOV T7, 15

L15:
    MOV AX, 2
    CMP AX, T7
    MOV AX, 0
    JGE S15_end
    MOV AX, 1
S15_end:
    MOV T8, AX

L16:
    MOV AX, T6
    AND AX, T8
    MOV T9, AX

L17:
    MOV T10, 0

L18:
    MOV AX, T10
    XOR AX, 1
    MOV T11, AX

L19:
    MOV AX, T9
    OR AX, T11
    MOV T12, AX

L20:
    MOV AX, T12
    CMP AX, 0
    JZ L46

L21:
    MOV T13, 15

L22:
    MOV AX, T13
    ADD AX, 2
    MOV T14, AX

L23:
    MOV AX, T14
    MOV somme, AX

L24:
    MOV i, 0

L25:
    MOV T15, 1

L26:
    MOV AX, T15
    CMP AX, 0
    JZ L45

L27:
    MOV SI, 0
    ADD SI, SI
    MOV AX, Tabint[SI]
    MOV T16, AX

L28:
    MOV AX, T16
    MOV T17, AX

L29:
    MOV SI, 0
    ADD SI, SI
    MOV AX, T17
    MOV Tabint[SI], AX

L30:
    MOV T18, 1

L31:
    MOV SI, 0
    ADD SI, SI
    MOV AX, Tabint[SI]
    MOV T19, AX

L32:
    MOV AX, T19
    CMP AX, 10
    MOV AX, 0
    JLE S32_end
    MOV AX, 1
S32_end:
    MOV T20, AX

L33:
    MOV AX, T18
    AND AX, T20
    MOV T21, AX

L34:
    MOV AX, T21
    CMP AX, 0
    JZ L39

L35:
    MOV SI, 0
    ADD SI, SI
    MOV AX, Tabint[SI]
    MOV T22, AX

L36:
    MOV AX, T22
    MOV CX, 1
    IMUL CX
    MOV T23, AX

L37:
    MOV SI, 0
    ADD SI, SI
    MOV AX, T23
    MOV Tabfloat[SI], AX

L38:
    JMP L42

L39:
    MOV SI, 0
    ADD SI, SI
    MOV AX, Tabint[SI]
    MOV T24, AX

L40:
    MOV AX, T24
    MOV CX, 2
    CWD
    IDIV CX
    MOV T25, AX

L41:
    MOV SI, 0
    ADD SI, SI
    MOV AX, T25
    MOV Tabfloat[SI], AX

L42:
    MOV T26, 1

L43:
    MOV AX, T26
    MOV i, AX

L44:
    JMP L25

L45:
    JMP L47

L46:
    MOV somme, 0

L47:
    MOV AX, 10
    CMP AX, Max
    MOV AX, 0
    JG S47_end
    MOV AX, 1
S47_end:
    MOV T27, AX

L48:
    MOV T28, 1

L49:
    MOV T29, 1

L50:
    MOV AX, T28
    OR AX, T29
    MOV T30, AX

L51:
    MOV AX, T27
    AND AX, T30
    MOV T31, AX

L52:
    MOV AX, T31
    CMP AX, 0
    JZ L69

L53:
    MOV T32, 11

L54:
    MOV AX, T32
    MOV x, AX

L55:
    MOV AX, T32
    CMP AX, 5
    MOV AX, 0
    JNE S55_end
    MOV AX, 1
S55_end:
    MOV T33, AX

L56:
    MOV AX, T33
    XOR AX, 1
    MOV T34, AX

L57:
    MOV AX, T34
    CMP AX, 0
    JZ L68

L58:
    MOV T35, 6

L59:
    MOV AX, T35
    MOV y, AX

L60:
    MOV AX, T32
    ADD AX, 1
    MOV T36, AX

L61:
    MOV SI, 0
    ADD SI, SI
    MOV AX, Tabint[SI]
    MOV T37, AX

L62:
    MOV SI, 1
    ADD SI, SI
    MOV AX, Tabint[SI]
    MOV T38, AX

L63:
    MOV AX, T37
    MOV CX, T38
    ADD AX, CX
    MOV T39, AX

L64:
    MOV AX, T32
    MOV CX, T35
    SUB AX, CX
    MOV T40, AX

L65:
    MOV AX, T39
    MOV CX, T40
    IMUL CX
    MOV T41, AX

L66:
    MOV SI, T36
    ADD SI, SI
    MOV AX, T41
    MOV Tabint[SI], AX

L67:
    JMP L55

L68:
    JMP L47

L69:
    MOV j, 1

L70:
    MOV T42, 1

L71:
    MOV AX, T42
    CMP AX, 0
    JZ L84

L72:
    MOV T43, 0

L73:
    MOV SI, T43
    ADD SI, SI
    MOV AX, Tabfloat[SI]
    MOV T44, AX

L74:
    MOV SI, 1
    ADD SI, SI
    MOV AX, Tabfloat[SI]
    MOV T45, AX

L75:
    MOV AX, T44
    MOV CX, T45
    ADD AX, CX
    MOV T46, AX

L76:
    MOV AX, T46
    MOV CX, 2
    CWD
    IDIV CX
    MOV T47, AX

L77:
    MOV SI, 1
    ADD SI, SI
    MOV AX, T47
    MOV Tabfloat[SI], AX

L78:
    MOV SI, 1
    ADD SI, SI
    MOV AX, Tabfloat[SI]
    MOV T48, AX

L79:
    MOV AX, moyenne
    MOV CX, T48
    ADD AX, CX
    MOV T49, AX

L80:
    MOV AX, T49
    MOV moyenne, AX

L81:
    MOV T50, 2

L82:
    MOV AX, T50
    MOV j, AX

L83:
    JMP L70

L84:
    MOV AX, T49
    MOV CX, 20
    CWD
    IDIV CX
    MOV T51, AX

L85:
    MOV AX, T51
    MOV moyenne, AX

L86:
    MOV AH, 01h
    INT 21h
    SUB AL, '0'
    MOV AH, 0
    MOV x, AX

L87:
    LEA DX, str0
    MOV AH, 09h
    INT 21h

L88:
    MOV AX, x
    CALL PRINT_INT

L89:
    LEA DX, str1
    MOV AH, 09h
    INT 21h

L90:
    MOV AX, 0
    CALL PRINT_INT

L91:
    LEA DX, str2
    MOV AH, 09h
    INT 21h

L92:
    MOV AX, T51
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
