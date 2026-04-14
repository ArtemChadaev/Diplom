-- Up migration for reference data tables.

CREATE TABLE "countries" (
    "code"     VARCHAR(3) PRIMARY KEY,
    "name_ru"  VARCHAR(255) NOT NULL
);

CREATE TABLE "atc_codes" (
    "code"     VARCHAR(10) PRIMARY KEY,
    "name_ru"  VARCHAR(255) NOT NULL
);

-- Seed some basic reference data
INSERT INTO "countries" ("code", "name_ru") VALUES
('RU', 'Россия'),
('BY', 'Беларусь'),
('KZ', 'Казахстан');

INSERT INTO "atc_codes" ("code", "name_ru") VALUES
('A01AB', 'Противомикробные препараты'),
('B01AC', 'Антиагреганты'),
('C01AA', 'Гликозиды наперстянки');
