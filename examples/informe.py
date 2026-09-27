# Informe numérico: para cada número del 1 al 20 muestra su numeral romano,
# si es primo y cuántos pasos necesita para llegar a 1 según Collatz.


def repeat(text, times):
    if times == 0:
        return ""
    return text + repeat(text, times - 1)


def roman_digit(digit, one, five, ten):
    if digit == 9:
        return one + ten
    elif digit >= 5:
        return five + repeat(one, digit - 5)
    elif digit == 4:
        return one + five
    return repeat(one, digit)


def roman(n):
    thousands = repeat("M", n // 1000)
    hundreds = roman_digit(n // 100 % 10, "C", "D", "M")
    tens = roman_digit(n // 10 % 10, "X", "L", "C")
    units = roman_digit(n % 10, "I", "V", "X")
    return thousands + hundreds + tens + units


def is_prime(n):
    if n < 2:
        return False
    divisor = 2
    while divisor * divisor <= n:
        if n % divisor == 0:
            return False
        divisor = divisor + 1
    return True


def collatz_steps(n):
    steps = 0
    while n != 1:
        if n % 2 == 0:
            n = n // 2
        else:
            n = 3 * n + 1
        steps = steps + 1
    return steps


def make_accumulator():
    total = 0

    def add(amount):
        nonlocal total
        total = total + amount
        return total

    return add


primes = make_accumulator()
steps = make_accumulator()
longest = ""
limit = 20

for n in range(1, limit + 1):
    numeral = roman(n)
    kind = "no primo"
    if is_prime(n):
        kind = "primo"
        primes(1)
    if len(numeral) > len(longest):
        longest = numeral
    print(n, numeral, kind, collatz_steps(n))
    steps(collatz_steps(n))

print("primos:", primes(0))
print("promedio de pasos de Collatz:", steps(0) / limit)
print("numeral más largo:", longest)
print("año actual:", roman(2026))
