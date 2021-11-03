-- migrate:up


CREATE TABLE klausimų_kategorijos
(
	id serial,
	name char (22) NOT NULL,
	PRIMARY KEY(id)
);
INSERT INTO klausimų_kategorijos(id, name)VALUES (1, 'mokėjimai');
INSERT INTO klausimų_kategorijos(id, name)VALUES (2, 'vairuotojo_pažymėjimas');
INSERT INTO klausimų_kategorijos(id, name)VALUES (3, 'rezervacijos');
INSERT INTO klausimų_kategorijos(id, name)VALUES (4, 'kita');

CREATE TABLE kuro_tipai
(
	id serial,
	name char (9) NOT NULL,
	PRIMARY KEY(id)
);
INSERT INTO kuro_tipai(id, name)VALUES (1, 'benzinas');
INSERT INTO kuro_tipai(id, name)VALUES (2, 'dyzelinas');
INSERT INTO kuro_tipai(id, name)VALUES (3, 'elektra');

CREATE TABLE mokėjimo_būsenos
(
	id serial,
	name char (11) NOT NULL,
	PRIMARY KEY(id)
);
INSERT INTO mokėjimo_būsenos(id, name)VALUES (1, 'sėkmingas');
INSERT INTO mokėjimo_būsenos(id, name)VALUES (2, 'atmestas');
INSERT INTO mokėjimo_būsenos(id, name)VALUES (3, 'neužbaigtas');

CREATE TABLE pavarų_dėžės
(
	id serial,
	name char (10) NOT NULL,
	PRIMARY KEY(id)
);
INSERT INTO pavarų_dėžės(id, name) VALUES(1, 'automatinė');
INSERT INTO pavarų_dėžės(id, name) VALUES(2, 'mechaninė');

CREATE TABLE rolės
(
	id serial,
	name char (32) NOT NULL,
	PRIMARY KEY(id)
);
INSERT INTO rolės(id, name) VALUES(1, 'klientas');
INSERT INTO rolės(id, name) VALUES(2, 'klientų_aptarnavimo_specialistas');
INSERT INTO rolės(id, name) VALUES(3, 'administratorius');

CREATE TABLE vairuotojo_pažymėjimo_būsenos
(
	id serial,
	name char (12) NOT NULL,
	PRIMARY KEY(id)
);
INSERT INTO vairuotojo_pažymėjimo_būsenos(id, name)VALUES (1, 'pateiktas');
INSERT INTO vairuotojo_pažymėjimo_būsenos(id, name)VALUES (2, 'patvirtintas');
INSERT INTO vairuotojo_pažymėjimo_būsenos(id, name)VALUES (3, 'atmestas');

CREATE TABLE automobiliai
(
	valstybiniai_numeriai varchar (255),
	markė varchar (255),
	modelis varchar (255),
	spalva varchar (255),
	pozicijos_platuma decimal,
	pozicijos_ilguma decimal,
	minutės_kaina decimal,
	valandos_kaina decimal,
	paros_kaina decimal,
	kilometro_kaina decimal,
	kondicionierius boolean,
	usb boolean,
	bluetooth boolean,
	navigacija boolean,
	vaikiška_kėdutė boolean,
	pašalintas boolean,
	pavarų_dėžė integer,
	kuro_tipas integer,
	id serial,
	PRIMARY KEY(id),
	FOREIGN KEY(pavarų_dėžė) REFERENCES pavarų_dėžės (id),
	FOREIGN KEY(kuro_tipas) REFERENCES kuro_tipai (id)
);

CREATE TABLE dažniausiai_užduodami_klausimai
(
	klausimas varchar (255),
	atsakymas varchar (255),
	kategorija integer,
	id serial,
	PRIMARY KEY(id),
	FOREIGN KEY(kategorija) REFERENCES klausimų_kategorijos (id)
);

CREATE TABLE vartotojai
(
	vardas varchar (255),
	pavardė varchar (255),
	el_paštas varchar (255),
	gimimo_data date,
	slaptažodis varchar (255),
	balansas integer,
	asmens_kodas varchar (255),
	rolė integer,
	id serial,
	PRIMARY KEY(id),
	FOREIGN KEY(rolė) REFERENCES rolės (id)
);

CREATE TABLE mokėjimai
(
	suma decimal,
	būsena integer,
	id serial,
	fk_Vartotojas integer NOT NULL,
	PRIMARY KEY(id),
	FOREIGN KEY(būsena) REFERENCES mokėjimo_būsenos (id),
	CONSTRAINT atlieka FOREIGN KEY(fk_Vartotojas) REFERENCES vartotojai (id)
);

CREATE TABLE rezervacijos
(
	sukurta timestamp with time zone,
	atšaukta timestamp with time zone,
	pradzios_adresas varchar (255),
	pabaigos_adresas varchar (255),
	id serial,
	fk_Automobilis integer NOT NULL,
	fk_Vartotojas integer NOT NULL,
	PRIMARY KEY(id),
	CONSTRAINT priklauso FOREIGN KEY(fk_Automobilis) REFERENCES automobiliai (id),
	CONSTRAINT sukuria FOREIGN KEY(fk_Vartotojas) REFERENCES vartotojai (id)
);

CREATE TABLE sesijos
(
	galiojimo_pabaiga timestamp with time zone,
	id serial,
	fk_Vartotojas integer NOT NULL,
	PRIMARY KEY(id),
	CONSTRAINT sukuriama FOREIGN KEY(fk_Vartotojas) REFERENCES vartotojai (id)
);

CREATE TABLE užklausos
(
	sukurta timestamp with time zone,
	užbaigta timestamp with time zone,
	id serial,
	fk_klientas integer NOT NULL,
	fk_klientų_aptarnavimo_specialistas integer,
	PRIMARY KEY(id),
	CONSTRAINT sukuria FOREIGN KEY(fk_klientas) REFERENCES vartotojai (id),
	CONSTRAINT priima FOREIGN KEY(fk_klientų_aptarnavimo_specialistas) REFERENCES vartotojai (id)
);

CREATE TABLE vairuotojo_pažymėjimai
(
	nr varchar (255),
	galiojimo_pabaiga date,
	būsena integer,
	id serial,
	fk_Vartotojas integer NOT NULL,
	PRIMARY KEY(id),
	FOREIGN KEY(būsena) REFERENCES vairuotojo_pažymėjimo_būsenos (id),
	CONSTRAINT prideda FOREIGN KEY(fk_Vartotojas) REFERENCES vartotojai (id)
);

CREATE TABLE įvertinimai
(
	žvaigždutės integer,
	komentaras varchar (255),
	data timestamp with time zone,
	id serial,
	fk_Uzklausa integer NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(fk_Uzklausa),
	CONSTRAINT paliekamas FOREIGN KEY(fk_Uzklausa) REFERENCES užklausos (id)
);

CREATE TABLE kelionės
(
	pradžios_laikas timestamp with time zone,
	pabaigos_laikas timestamp with time zone,
	pabaigos_taško_platuma decimal,
	pabaigos_taško_ilguma decimal,
	trukmė time,
	id serial,
	fk_Rezervacija integer NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(fk_Rezervacija),
	CONSTRAINT priklauso FOREIGN KEY(fk_Rezervacija) REFERENCES rezervacijos (id)
);

CREATE TABLE vairuotojo_pažymėjimo_nuotraukos
(
	nuoroda varchar (255),
	id serial,
	fk_Vairuotojo_pazymejimas integer NOT NULL,
	PRIMARY KEY(id),
	CONSTRAINT turi FOREIGN KEY(fk_Vairuotojo_pazymejimas) REFERENCES vairuotojo_pažymėjimai (id)
);

CREATE TABLE žinutės
(
	tekstas varchar (255),
	išsiųsta timestamp with time zone,
	id serial,
	fk_Vartotojas integer NOT NULL,
	fk_Uzklausa integer NOT NULL,
	PRIMARY KEY(id),
	CONSTRAINT siunčia FOREIGN KEY(fk_Vartotojas) REFERENCES vartotojai (id),
	CONSTRAINT priklauso FOREIGN KEY(fk_Uzklausa) REFERENCES užklausos (id)
);


-- migrate:down

