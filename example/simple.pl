BeginProject simple_test;

Setup :
  define x : integer = 5;
  define y : integer = 10;
  const Pi : float = 3.14;

Run :
{
  x <- 5;
  y <- x + 10;
  out("Result: ", y);
}

EndProject ;
