BeginProject TestSimple;

Setup:
define x : integer;
define i : integer;

Run:
{
    x <-- 0;
    i <-- 0;

    loop while (i < 5)
    {
        x <-- x + 2;
        i <-- i + 1;
    } endloop;

    out("Resultat: ", x);
}
EndProject;
