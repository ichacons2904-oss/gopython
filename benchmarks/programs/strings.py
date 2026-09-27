matches = 0
for i in range(300000):
    text = ""
    for j in range(10):
        text = text + "ab"
    if text == "abababababababababab":
        matches = matches + 1
print(matches)
