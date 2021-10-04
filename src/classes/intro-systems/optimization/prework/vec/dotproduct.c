#include "vec.h"

/*
I am probably just tired, but I didn't notice _any_ change
after these optimizations. I don't feel like I know if I'm timing this
accurately.

See my bench.sh script.
*/
data_t
dotproduct(vec_ptr u, vec_ptr v)
{
  data_t sum = 0, sum1 = 0, sum2 = 0, u_val, v_val, u_val1, v_val1, u_val2,
         v_val2;

  int length = u->len;
  long i;
  for (i = 0; i < length;
       i += 3) { // we can assume both vectors are same length

    u_val = u->data[i];
    v_val = v->data[i];

    u_val1 = u->data[i + 1];
    v_val1 = v->data[i + 1];

    u_val2 = u->data[i + 2];
    v_val2 = v->data[i + 2];

    sum += u_val * v_val;
    sum1 += u_val1 * v_val1;
    sum2 += u_val2 * v_val2;
  }

  for (; i < length; i++) {
    u_val = u->data[i];
    v_val = v->data[i];

    sum += u_val * u_val;
  }

  return sum + sum1 + sum2;
}
