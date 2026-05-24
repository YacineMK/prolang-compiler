BeginProject TestMixedExpr;

Setup:
define x : integer;
define y : integer;
define r1 : integer;
define r2 : integer;
const C : integer = 10;

Run:
{
    x <-- 5;
    y <-- 3;

    r1 <-- x * C + y;
    r2 <-- (x + y) * (x - y);

    out("r1 = ", r1);
    out("r2 = ", r2);
}
EndProject;
