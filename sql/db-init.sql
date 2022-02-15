DROP TABLE IF EXISTS spells;
CREATE TABLE spells (
	name TEXT CONSTRAINT spells_pk PRIMARY KEY
	,level TEXT
	,school TEXT
	,casting_time TEXT
	,range TEXT
	,components TEXT
	,duration TEXT
	,classes TEXT
	,description TEXT
	,source TEXT
);

DROP TABLE IF EXISTS monsters;
CREATE TABLE monsters (
	name TEXT CONSTRAINT monsters_pk PRIMARY KEY
	,hit_points TEXT
	,speed TEXT
	,strength TEXT
	,dexterity TEXT
	,constitution TEXT
	,intelligence TEXT
	,wisdom TEXT
	,charisma TEXT
	,skills TEXT
	,damage_resistances TEXT
	,condition_immunities TEXT
	,senses TEXT
	,languages TEXT
	,challenge TEXT
	,traits TEXT
	,actions TEXT
	,source TEXT
);
