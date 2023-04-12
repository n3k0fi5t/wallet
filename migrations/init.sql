Drop Table If Exists account;
Drop Table If Exists user;
Drop Table If Exists TransactionLog;

CREATE TABLE IF NOT EXISTS user (
   id INT UNSIGNED NOT NULL AUTO_INCREMENT,
   userID varchar(50) NOT NULL,
   fullname varchar(50) NOT NULL,
   PRIMARY KEY (id),
   UNIQUE KEY userID (userID)
);

CREATE TABLE IF NOT EXISTS TransactionLog (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	accountID varchar(50) NOT NULL,
	action int(10) NOT NULL DEFAULT 0,
	amount BIGINT NOT NULL DEFAULT 0,
	timestampMS BIGINT NOT NULL,
	tradeID varchar(50) NOT NULL,
	PRIMARY KEY (id),
	KEY accountID (accountID),
	KEY timestampMS (timestampMS),
	KEY tradeID (tradeID)
);

CREATE TABLE IF NOT EXISTS account (
   id INT UNSIGNED NOT NULL AUTO_INCREMENT,
   accountID varchar(50) NOT NULL UNIQUE,
   balance BIGINT NOT NULL DEFAULT 0,
   PRIMARY KEY (id)
);

-- Pseudo user and account as system account
INSERT INTO user (userID, fullname) VALUES ('c1e395d9-8c00-4124-819a-85b0402900cf', 'PseudoUser');
INSERT INTO account (balance, accountID) VALUES (9223372036854775807, 'c1e395d9-8c00-4124-819a-85b0402900cf');

-- Seed users and accounts
INSERT INTO user (userID, fullname) VALUES ('935f871a-660f-4f19-801e-916c04bb0324', 'Tim');
INSERT INTO account (balance, accountID) VALUES (0, '935f871a-660f-4f19-801e-916c04bb0324');

INSERT INTO user (userID, fullname) VALUES ('a89b7b78-b9c1-4129-8cff-380bf53f3a49', 'Alex');
INSERT INTO account (balance, accountID) VALUES (0, 'a89b7b78-b9c1-4129-8cff-380bf53f3a49');

INSERT INTO user (userID, fullname) VALUES ('a98cd0f5-d6b2-4899-a1fb-ddf308d6f5c8', 'Arthur');
INSERT INTO account (balance, accountID) VALUES (0, 'a98cd0f5-d6b2-4899-a1fb-ddf308d6f5c8');

INSERT INTO user (userID, fullname) VALUES ('a679ac51-08e8-45c7-80d7-019bf9dad64b', 'Ray');
INSERT INTO account (balance, accountID) VALUES (0, 'a679ac51-08e8-45c7-80d7-019bf9dad64b');

INSERT INTO user (userID, fullname) VALUES ('55b36756-6089-4756-bbd2-b0f66e50ee07', 'HD');
INSERT INTO account (balance, accountID) VALUES (0, '55b36756-6089-4756-bbd2-b0f66e50ee07');

INSERT INTO user (userID, fullname) VALUES ('5a1e760e-76ea-4709-98ba-e1a701a4d340', 'peko');
INSERT INTO account (balance, accountID) VALUES (0, '5a1e760e-76ea-4709-98ba-e1a701a4d340');

INSERT INTO user (userID, fullname) VALUES ('201bef83-cc46-4acb-9c25-2eef60a59a9a', 'miko');
INSERT INTO account (balance, accountID) VALUES (0, '201bef83-cc46-4acb-9c25-2eef60a59a9a');

INSERT INTO user (userID, fullname) VALUES ('1c3e7209-fb42-4643-bfa6-c6a3fb42bf92', 'rushia');
INSERT INTO account (balance, accountID) VALUES (0, '1c3e7209-fb42-4643-bfa6-c6a3fb42bf92');

INSERT INTO user (userID, fullname) VALUES ('084e135f-78c7-406e-a347-94e38fa55b60', 'gura');
INSERT INTO account (balance, accountID) VALUES (0, '084e135f-78c7-406e-a347-94e38fa55b60');

INSERT INTO user (userID, fullname) VALUES ('8a180d2b-0965-4095-ba17-a880d196f04d', 'Ame');
INSERT INTO account (balance, accountID) VALUES (0, '8a180d2b-0965-4095-ba17-a880d196f04d');
