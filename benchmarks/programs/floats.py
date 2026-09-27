total = 0.0
sign = 1.0
for k in range(3000000):
    total = total + sign / (2 * k + 1)
    sign = -sign
print(4 * total)
