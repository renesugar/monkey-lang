#!./monkey-lang

prime := fn(n) {
  # Corner cases 
  if (n <= 1) {
    return false
  }
  if (n <= 3) {
    return true
  }
  
  # This is checked so that we can skip  
  # middle five numbers in below loop 
  if (n % 2 == 0 || n % 3 == 0) {
    return false
  }
  
  i := 5
  while (i * i <= n) {
    if (n % i == 0 || n % (i + 2) == 0) {
      return false
    }
    i = i + 6
  }

  return true
} 

if (len(args()) == 1) {
  N := int(args()[0])
} else {
  N := 100
}

n := 1
while (n < N) {
  if (prime(n)) {
    print(n)
  }
  n = n + 1
}
