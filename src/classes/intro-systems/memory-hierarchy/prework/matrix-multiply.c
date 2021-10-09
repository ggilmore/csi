/*
Naive code for multiplying two matrices together.

There must be a better way!
*/

#include <stdio.h>
#include <stdlib.h>
double**
my_matrix_alloc(int m, int n)
{
  double** matrix = malloc(m * sizeof(double));

  for (int i = 0; i < m; i++)
    matrix[i] = calloc(n, sizeof(double));

  return matrix;
}

// Free the entirety of an m row matrix
void
my_matrix_free(double** matrix, int m)
{
  for (int i = 0; i < m; i++)
    free(matrix[i]);
  free(matrix);
}

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void
matrix_multiply(double** C,
                double** A,
                double** B,
                int a_rows,
                int a_cols,
                int b_cols)
{
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }
}

void
fast_matrix_multiply(double** c,
                     double** a,
                     double** b,
                     int a_rows,
                     int a_cols,
                     int b_cols)
{

  // return matrix_multiply(c, a, b, a_rows, a_cols, b_cols);
  // TODO: write a faster implementation here!
  double** transpose;
  transpose = my_matrix_alloc(a_rows, a_rows);
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < a_cols; j++) {
      transpose[j][i] = b[i][j];
    }
  }

  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < a_rows; j++) {
      c[i][j] = 0;
      for (int k = 0; k < a_rows; k++)
        c[i][j] += a[i][k] * transpose[j][k];
    }
  }

  // my_matrix_free(transpose, a_rows);
}
