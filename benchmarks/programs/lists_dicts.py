counts = {}
values = []
for i in range(600000):
    key = i % 1000
    values.append(key)
    counts[key] = counts.get(key, 0) + 1

total = 0
for value in values:
    total = total + counts[value]

pairs = []
for i in range(1000):
    pairs.append((i, counts[i]))

print(len(values), len(counts), total, pairs[999])
