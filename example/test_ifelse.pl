BeginProject TestIfElse;

Setup:
define a : integer;
define b : integer;
define result : integer;

Run:
{
    a <-- 15;
    b <-- 10;

    if (a > b) then:
    {
        result <-- a - b;
    } else {
        result <-- b - a;
    } endIf;

    out("result = ", result);

    if (a < b) then:
    {
        result <-- 100;
    } else {
        result <-- 200;
    } endIf;

    out("result2 = ", result);
}
EndProject;
