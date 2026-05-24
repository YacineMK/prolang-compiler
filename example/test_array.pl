BeginProject TestArray;

Setup:
define i : integer;
define Tab : [integer; 10];
define sum : integer;

Run:
{
    i <-- 0;
    sum <-- 0;

    loop while (i < 5)
    {
        Tab[i] <-- i * 2;
        sum <-- sum + Tab[i];
        i <-- i + 1;
    } endloop;

    out("sum = ", sum);
}
EndProject;
