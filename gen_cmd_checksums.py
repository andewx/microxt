cmds = ['RADC', 'RFFT', 'PDAT', 'TDAT', 'DDAT', 'DONE', 'RPST', 'SRPS', 'GRPS', 'RESP', 'GNFD']
i = 0
sums = []
for cmd in cmds:
    sum = 0
    for char in cmd:
        sum += ord(char)
    sums.append({cmds[i],sum})
    i += 1

print(sums)
