BeginProject TestLoop;

Setup:
define i : integer;
define x : integer;

Run:
{
    x <-- 0;
    i <-- 0;

    loop while (i < 10)
    {
        x <-- x + 3;
        i <-- i + 1;
    } endloop;

    out("x = ", x);
    out("i = ", i);
}
EndProject;
