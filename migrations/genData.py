

template = '''INSERT INTO user (userID, fullname) VALUES ('{}', '{}');
INSERT INTO account (balance, accountID) VALUES (0, '{}');'''

uids = ['935f871a-660f-4f19-801e-916c04bb0324', 'a89b7b78-b9c1-4129-8cff-380bf53f3a49', 'a98cd0f5-d6b2-4899-a1fb-ddf308d6f5c8', 'a679ac51-08e8-45c7-80d7-019bf9dad64b', '55b36756-6089-4756-bbd2-b0f66e50ee07', '5a1e760e-76ea-4709-98ba-e1a701a4d340', '201bef83-cc46-4acb-9c25-2eef60a59a9a', '1c3e7209-fb42-4643-bfa6-c6a3fb42bf92', '084e135f-78c7-406e-a347-94e38fa55b60', '8a180d2b-0965-4095-ba17-a880d196f04d']
names = ['Tim', 'Alex', 'Arthur', 'Ray', 'HD', 'peko', 'miko', 'rushia', 'gura', 'Ame']

def generateData():
    for u, n in zip(uids, names):
        #print("\"{}\": \"{}\",".format(n, u))
        print(template.format(u, n, u))

if __name__ == '__main__':
    generateData()
