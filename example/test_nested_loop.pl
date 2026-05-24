BeginProject TestNestedLoop;

Setup:
define i : integer;
define j : integer;
define count : integer;

Run:
{
    count <-- 0;
    i <-- 0;

    loop while (i < 3)
    {
        j <-- 0;
        loop while (j < 2)
        {
            count <-- count + 1;
            j <-- j + 1;
        } endloop;
        i <-- i + 1;
    } endloop;

    out("count = ", count);
}
EndProject;
