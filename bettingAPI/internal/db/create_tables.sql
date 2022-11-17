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
	constraint fkTip foreign key(tip) references tips(IDtip)
);

create table league_offers( -- relation n:m between league and offer
	IDleague smallint,
    IDoffer smallint,
    constraint pkOffers primary key(IDleague, IDoffer),
    constraint fkLeague foreign key (IDleague) references league(IDleague),
    constraint fkOffer foreign key (IDoffer) references offer(IDoffer)
);

create table tips(
	IDtip smallint primary key, -- starts with 2xxx
    tip varchar(25),
	coefficient smallint,
    constraint fkCoefficient foreign key(coefficient) references coefficients(IDcoefficient)
);

create table coefficients(
	IDcoefficient smallint primary key,
    worth decimal(4,2)
);


