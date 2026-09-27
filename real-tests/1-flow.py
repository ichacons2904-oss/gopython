print("--- IFs ---")

if True:
    print("OK")
else:
    print("ERROR")

if False:
    print("ERROR")
else:
    print("OK")

print("--- WHILEs ---")
i = 0
match = 0
while i <= 3:
    match = i
    i = i + 1

if match == 3:
    print("OK")
else:
    print("ERROR")

print("--- FOR ---")
match = 0
for i in range(0, 6):
    match = i

if match == 5:
    print("OK")
else:
    print("ERROR")
