drop table league_offers;
drop table offer_tips;
drop table offers;
drop table leagues;
drop table users;

create table leagues(
	league_id int primary key auto_increment,
	title varchar(30)
);
create table offers(
	offer_id int primary key,
    game varchar(40),
    time_played datetime,
    tv_channel varchar(10),
    has_statistics bool
);

create table league_offers(
	league_id int not null,
    offer_id int not null,
    constraint fk_league_offers_leagues foreign key (league_id) references leagues(league_id),
    constraint fk_league_offers_offers foreign key (offer_id) references offers(offer_id),
    constraint pk_league_offers primary key(league_id, offer_id)
);

create table offer_tips(
	offer_id int not null,
    tip varchar(25) not null,
    coefficient decimal(4,2),
    constraint fk_offer_tips_offers foreign key (offer_id) references offers(offer_id),
    constraint pk_offer_tip primary key(offer_id, tip)
);

