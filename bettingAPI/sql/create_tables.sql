drop table league_offers;
drop table bet;
drop table user_bet_slips;
drop table offer_tips;
drop table offers;
drop table leagues;
drop table users;


create table leagues(
	id int primary key auto_increment,
	title varchar(30)
);
create table offers(
	id int primary key,
    game varchar(40),
    time_played datetime,
    tv_channel varchar(10),
    has_statistics bool
);

create table league_offers(
	league_id int not null,
    offer_id int not null,
    constraint fk_league_offers_from_leagues foreign key (league_id) references leagues(id),
    constraint fk_league_offers_from_offers foreign key (offer_id) references offers(id),
    constraint pk_league_offers primary key(league_id, offer_id)
);

create table offer_tips(
	offer_id int not null,
    tip varchar(25) not null,
    coefficient decimal(4,2),
    constraint fk_offer_tips_from_offers foreign key (offer_id) references offers(id),
    constraint pk_offer_tip primary key(offer_id, tip)
);

create table users(
	id int primary key auto_increment,
    username varchar(25),
	email varchar(80),
	password_hash varchar(60),
    first_name varchar(25),
    last_name varchar(25),
	birth_date date,
    balance decimal(6,2)
);

create table user_bet_slips(
    id int primary key auto_increment,
	user_id int,
	stake decimal(6,2),
    coefficient decimal(6,2),
    payout decimal(7,2),
    constraint fk_user_bet_slips_from_users foreign key(user_id) references users(id)
);

create table bet(
	id int primary key auto_increment,
	user_bet_slip_id int,
    offer_id int,
    tip varchar(25),
    coefficient decimal(4,2),
	constraint fk_user_bet_slip_id_from_user_bet_slips foreign key(user_bet_slip_id) references user_bet_slips(id)
);