/*

Two different ways to loop over an array of arrays.

Spotted at:
http://stackoverflow.com/questions/9936132/why-does-the-order-of-the-loops-affect-performance-when-iterating-over-a-2d-arra

*/

#include "loop-order.h"

void
option_one()
{
  int i, j;
  static int x[ROWS][COLS];
  for (i = 0; i < ROWS; i++) {
    for (j = 0; j < COLS; j++) {
      x[i][j] = i + j;
    }
  }
}

void
option_two()
{
  int i, j;
  static int x[ROWS][COLS];
  for (i = 0; i < COLS; i++) {
    for (j = 0; j < ROWS; j++) {
      x[j][i] = i + j;
    }
  }
}

int
main()
{

  // option_one();
  option_two();
  return 0;
}
