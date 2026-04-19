BeginProject test_full_compiler;

Setup :
    %% ============================================================
    %% SECTION 1: Declaration de variables simples (sans init)
    %% ============================================================

    %% Variable entiere simple
    define x : integer;

    %% Variable flottante simple
    define moyenne : float;

    %% Plusieurs variables du meme type (separees par |)
    define a | b | c : integer;

    %% Plusieurs variables flottantes
    define f1 | f2 | f3 : float;

    %% ============================================================
    %% SECTION 2: Declaration de variables avec initialisation
    %% ============================================================

    define compteur : integer = 0;
    define total : integer = 100;
    define pi_approx : float = 3.14;
    define zero_f : float = 0.0;

    %% Valeurs signees (entre parentheses)
    define neg_int : integer = (-5);
    define pos_int : integer = (+5);
    define neg_flt : float = (-5.76);
    define pos_flt : float = (+ 568.77);

    %% ============================================================
    %% SECTION 3: Declaration de constantes
    %% ============================================================

    const Pi : float = 3.14159;
    const MAX_VAL : integer = 32767;
    const MIN_VAL : integer = (-32768);
    const ZERO : integer = 0;
    const RATE : float = 0.5;

    %% ============================================================
    %% SECTION 4: Declaration de tableaux
    %% ============================================================

    define Tab : [integer ; 20];
    define notes : [float ; 10];
    define data : [integer ; 1];

    //* 
        Identificateurs valides utilises dans ce programme :
        - x, a, b, c, f1, f2, f3 : courts, simples
        - compteur, total, moyenne : longueur <= 14
        - neg_int, pos_int : avec tiret bas unique (non consecutif, non terminal)
        - pi_approx : avec tiret bas au milieu
    *//

Run :
{
    //* ============================================================
        SECTION 5: Instructions d'affectation - cas simples
       ============================================================ *//

    %% Affectation entiere
    x <- 23;
    a <- 0;
    b <- (-10);
    c <- (+99);

    %% Affectation flottante
    moyenne <- 0.0;
    f1 <- 3.14;
    f2 <- (-5.76);
    f3 <- (+ 568.77);

    //* ============================================================
        SECTION 6: Affectation avec expressions arithmetiques
       ============================================================ *//

    %% Addition et soustraction
    x <- a + b;
    x <- a - b;
    x <- a + b - c;

    %% Multiplication et division
    x <- a * b;
    x <- a / b;

    %% Expressions complexes (associativite gauche)
    x <- a + b * c;
    x <- (a + b) * c;
    x <- a * b + c * x;
    x <- (a + 2) / (5 - b + (-3));

    %% Expressions avec flottants
    moyenne <- f1 + f2;
    moyenne <- f1 * f2 - f3;
    pi_approx <- Pi * f1;

    //* ============================================================
        SECTION 7: Affectation sur tableaux
       ============================================================ *//

    Tab[0] <- 10;
    Tab[1] <- (-5);
    Tab[5] <- a + b;
    Tab[19] <- x * 2;

    notes[0] <- 15.5;
    notes[9] <- (+ 18.75);

    %% Lecture d'element de tableau dans expression
    x <- Tab[0] + Tab[1];
    moyenne <- notes[0] + notes[9];
    b <- Tab[5] * c;

    %% Affectation tableau avec expression complexe
    Tab[2] <- (a + 2) / (5 - b + (-3));

    //* ============================================================
        SECTION 8: Instructions d'entree / sortie
       ============================================================ *//

    %% Lecture d'une variable
    in(x);
    in(moyenne);
    in(a);

    %% Ecriture d'une variable
    out("Valeur de x : ", x);
    out("Moyenne : ", moyenne);

    %% Ecriture d'un seul argument
    out("Debut du programme");

    %% Lecture puis affichage
    in(b);
    out("Vous avez saisi: ", b);

    //* ============================================================
        SECTION 9: Instruction conditionnelle IF - cas simples
       ============================================================ *//

    %% Condition avec operateur >
    if (a > b) then:
    {
        x <- a;
    } else {
        x <- b;
    } endIf;

    %% Condition avec operateur <
    if (x < 10) then:
    {
        out("x est inferieur a 10");
    } else {
        out("x est superieur ou egal a 10");
    } endIf;

    %% Condition avec >=
    if (moyenne >= 10.0) then:
    {
        out("admis");
    } else {
        out("recale");
    } endIf;

    %% Condition avec <=
    if (compteur <= ZERO) then:
    {
        compteur <- 1;
    } else {
        compteur <- compteur + 1;
    } endIf;

    %% Condition avec ==
    if (a == 0) then:
    {
        out("a est nul");
    } else {
        out("a non nul");
    } endIf;

    %% Condition avec !=
    if (b != 0) then:
    {
        x <- a / b;
    } else {
        out("division par zero impossible");
    } endIf;

    //* ============================================================
        SECTION 10: Condition avec operateurs logiques
       ============================================================ *//

    %% Condition AND
    if ((a > 0) AND (b > 0)) then:
    {
        out("a et b positifs");
    } else {
        out("au moins un negatif");
    } endIf;

    %% Condition OR
    if ((a == 0) OR (b == 0)) then:
    {
        out("au moins un nul");
    } else {
        x <- a / b;
    } endIf;

    %% Condition NON
    if (NON(a == 0)) then:
    {
        x <- 100 / a;
    } else {
        out("a est nul, division impossible");
    } endIf;

    %% Condition composee AND + OR
    if (((a > 0) AND (b > 0)) OR (c == 0)) then:
    {
        out("condition complexe vraie");
    } else {
        out("condition complexe fausse");
    } endIf;

    //* ============================================================
        SECTION 11: Boucle WHILE simple
       ============================================================ *//

    compteur <- 0;
    loop while (compteur < 10)
    {
        compteur <- compteur + 1;
    } endloop;

    %% While avec corps plus complexe
    x <- 1;
    loop while (x <= 100)
    {
        total <- total + x;
        x <- x + 1;
        out("iteration: ", x);
    } endloop;

    %% While avec condition logique
    loop while ((a < 10) AND (b < 10))
    {
        a <- a + 1;
        b <- b + 1;
    } endloop;

    //* ============================================================
        SECTION 12: Boucle FOR
       ============================================================ *//

    %% For simple
    for compteur in 1 to 10
    {
        out("i = ", compteur);
    } endfor;

    %% For avec corps complexe
    for x in 0 to 19
    {
        Tab[x] <- x * 2;
    } endfor;

    %% For avec calcul
    total <- 0;
    for a in 1 to 100
    {
        total <- total + a;
    } endfor;
    out("Somme 1..100 = ", total);

    //* ============================================================
        SECTION 13: Boucles imbriquees
       ============================================================ *//

    %% While dans while
    a <- 0;
    loop while (a < 5)
    {
        b <- 0;
        loop while (b < 5)
        {
            x <- a * 5 + b;
            b <- b + 1;
        } endloop;
        a <- a + 1;
    } endloop;

    %% For dans for
    for a in 0 to 4
    {
        for b in 0 to 4
        {
            Tab[a] <- Tab[a] + b;
        } endfor;
    } endfor;

    %% While dans for
    for a in 0 to 9
    {
        b <- 0;
        loop while (b < a)
        {
            notes[a] <- notes[a] + 1.0;
            b <- b + 1;
        } endloop;
    } endfor;

    %% For dans while
    compteur <- 0;
    loop while (compteur < 3)
    {
        for x in 1 to 5
        {
            Tab[x] <- Tab[x] * compteur;
        } endfor;
        compteur <- compteur + 1;
    } endloop;

    //* ============================================================
        SECTION 14: IF imbrique dans boucle et vice-versa
       ============================================================ *//

    for a in 0 to 9
    {
        if (Tab[a] > 0) then:
        {
            out("Tab positif a index ", a);
        } else {
            Tab[a] <- 0;
        } endIf;
    } endfor;

    loop while (compteur < 10)
    {
        if ((compteur == 5) OR (compteur == 0)) then:
        {
            out("borne atteinte");
        } else {
            total <- total + compteur;
        } endIf;
        compteur <- compteur + 1;
    } endloop;

    //* ============================================================
        SECTION 15: Expressions arithmetiques avec constantes
       ============================================================ *//

    x <- MAX_VAL;
    b <- MIN_VAL;
    moyenne <- Pi * f1 * f1;
    f2 <- RATE * f1 + zero_f;
    c <- MAX_VAL - 1;

    //* ============================================================
        SECTION 16: Priorite et associativite des operateurs
       ============================================================ *//

    %% * et / ont priorite sur + et -
    x <- 2 + 3 * 4;        %% = 14 (non 20)
    x <- 10 - 2 * 3;       %% = 4  (non 24)
    x <- 8 / 2 + 1;        %% = 5  (non 2)

    %% Associativite gauche
    x <- 10 - 3 - 2;       %% = 5  (non 9)
    x <- 12 / 4 / 3;       %% = 1  (non 9)
    x <- 2 + 3 + 4 + 5;    %% = 14

    %% Parentheses pour forcer priorite
    x <- (2 + 3) * 4;      %% = 20
    x <- (10 - 2) * (3 + 1);

    %% Comparaisons ont priorite sur logiques
    if ((a + 1 > 0) AND (b - 1 < 10)) then:
    {
        out("ok");
    } else {
        out("non");
    } endIf;

    //* ============================================================
        SECTION 17: Cas limites sur les types
       ============================================================ *//

    %% Entier valeur max et min (avec parentheses pour signes)
    x <- 32767;
    x <- (-32768);

    %% Float avec point decimal obligatoire
    moyenne <- 0.0;
    f1 <- 1.0;
    f2 <- (-0.5);
    f3 <- (+ 999.99);

    //* ============================================================
        SECTION 18: Commentaires mixtes dans le code
       ============================================================ *//

    %% Ceci est un commentaire sur une seule ligne
    x <- a + b; %% commentaire en fin de ligne

    //* 
        Ceci est un commentaire
        sur plusieurs lignes
        il peut contenir n'importe quoi :
        define, if, loop, for, etc.
    *//

    a <- x * 2;

    //* commentaire multiligne court *//

    b <- a - 1;

    //* ============================================================
        SECTION 19: Chaines de caracteres dans out
       ============================================================ *//

    out("Hello World");
    out("Valeur de pi = ", Pi);
    out("a = ", a, " b = ", b);
    out("Fin du programme");

    //* ============================================================
        SECTION 20: Scenario complet - calcul de moyenne d'un tableau
       ============================================================ *//

    %% Saisie des notes
    for compteur in 0 to 9
    {
        in(notes[compteur]);
    } endfor;

    %% Calcul de la somme
    total <- 0;
    moyenne <- 0.0;
    for compteur in 0 to 9
    {
        moyenne <- moyenne + notes[compteur];
    } endfor;

    %% Calcul de la moyenne
    moyenne <- moyenne / 10;
    out("Moyenne = ", moyenne);

    %% Affichage selon resultat
    if (moyenne >= 10.0) then:
    {
        out("Resultat : ADMIS");
    } else {
        out("Resultat : RECALE");
    } endIf;

    //* ============================================================
        SECTION 21: Scenario complet - tri a bulles
       ============================================================ *//

    for a in 0 to 18
    {
        for b in 0 to (18 - a)
        {
            if (Tab[b] > Tab[b + 1]) then:
            {
                x <- Tab[b];
                Tab[b] <- Tab[b + 1];
                Tab[b + 1] <- x;
            } else {
                %% Rien a faire
                x <- 0;
            } endIf;
        } endfor;
    } endfor;

    out("Tableau trie :");
    for compteur in 0 to 19
    {
        out("Tab[", compteur, "] = ", Tab[compteur]);
    } endfor;

}
EndProject ;