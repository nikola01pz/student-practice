create table leagues(
	league_id int primary key,  -- starts with 1xxxx
	title varchar(30)
);

create table offers(
	offer_id int primary key, -- starts with 2xxxx
    game varchar(40),
    sport varchar(15),
    time_played varchar(25),
    tv_channel varchar(10)
);

create table league_offers( -- relation n:m between league and offer
	league_id int not null,
    offer_id int not null,
    constraint fk_league_offers_leagues foreign key (league_id) references leagues(league_id),
    constraint fk_league_offers_offers foreign key (offer_id) references offers(offer_id),
    constraint pk_league_offers primary key(league_id, offer_id)
);

create table offer_tips( -- starts with 3xxxx
	offer_id int not null,
    tip varchar(25) not null,
    columne decimal(4,2),
    constraint fk_offer_tips_offers foreign key (offer_id) references offers(offer_id),
    constraint pk_offer_tip primary key(offer_id, tip)
);
