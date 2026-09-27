total = 0
for i in range(1500):
    for j in range(1500):
        total = total + (i * j) % 7
print(total)
