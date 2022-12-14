drop table league_offers;
drop table user_bet_slips;
drop table bet_slip;
drop table bet;
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
    constraint fk_league_offers_from_leagues foreign key (league_id) references leagues(league_id),
    constraint fk_league_offers_from_offers foreign key (offer_id) references offers(offer_id),
    constraint pk_league_offers primary key(league_id, offer_id)
);

create table offer_tips(
	offer_id int not null,
    tip varchar(25) not null,
    coefficient decimal(4,2),
    constraint fk_offer_tips_from_offers foreign key (offer_id) references offers(offer_id),
    constraint pk_offer_tip primary key(offer_id, tip)
);

create table users(
	user_id int primary key auto_increment,
    username varchar(25),
	email varchar(80),
	password_hash varchar(60),
    first_name varchar(25),
    last_name varchar(25),
	birth_date date
);

create table bet(
	bet_id int primary key auto_increment,
    offer_id int,
    offer_tip varchar(25),
    constraint fk_bet_from_offer_tips foreign key (offer_id, offer_tip) references offer_tips(offer_id, tip)
);

create table bet_slip(
	bet_slip_id int primary key auto_increment,
    bet_id int,
    constraint fk_bet_slip_from_bet foreign key(bet_id) references bet(bet_id)
);

create table user_bet_slips(
	user_id int,
    user_bet_slip_id int,
	stake decimal(6,2),
    payout decimal(7,2),
    constraint fk_user_bet_slips_from_users foreign key(user_id) references users(user_id),
    constraint fk_user_bet_slips_from_bet_slip foreign key(user_bet_slip_id) references bet_slip(bet_slip_id),
    constraint pk_user_bet_slips primary key(user_id, user_bet_slip_id)
);
