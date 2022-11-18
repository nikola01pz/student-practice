create table league(
	IDleague smallint primary key,  -- starts with 0xxx
	title varchar(30),
    offer smallint
);

create table offer(
	IDoffer smallint primary key, -- starts with 1xxx
    game varchar(40),
    sport varchar(15),
    tip smallint,
    time_played datetime,
    tv_channel varchar(10),
	constraint fkTip foreign key(tip) references tip(IDtip)
);

create table tip(
	IDtip smallint primary key, -- starts with 2xxx
    tip varchar(25),
	coefficient smallint
);

create table league_offers( -- relation n:m between league and offer
	IDleague smallint not null,
    IDoffer smallint not null,
    constraint fkLeague foreign key (IDleague) references league(IDleague),
    constraint fkOffer foreign key (IDoffer) references offer(IDoffer),
    constraint pkOffers primary key(IDleague, IDoffer)
);