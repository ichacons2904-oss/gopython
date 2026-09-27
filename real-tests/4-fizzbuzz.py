def check(value, expected):
    if value != expected:
        print("ERROR")
    else:
        print("OK")


def fizzbuzz(n):
    res = ""
    for i in range(1, n + 1):
        if i % 3 == 0 and i % 5 == 0:
            res = res + "FizzBuzz"
        elif i % 3 == 0:
            res = res + "Fizz"
        elif i % 5 == 0:
            res = res + "Buzz"
    return res


print("--- FIZZ BUZZ ---")
check(fizzbuzz(2), "")
check(fizzbuzz(3), "Fizz")
check(fizzbuzz(5), "FizzBuzz")
check(fizzbuzz(6), "FizzBuzzFizz")
check(fizzbuzz(15), "FizzBuzzFizzFizzBuzzFizzFizzBuzz")
