--liquibase formatted sql

--changeset jaime.silvela:1 labels:kubecon-demo
--comment: let's start off with 50 bonds 
create table bonds as
    select 'bn_' || generate_series as bond
    from generate_series(1, 50);
--rollback DROP TABLE bonds;

--changeset jaime.silvela:2 labels:kubecon-demo
--comment: lets add a bunch of factors
create table factors as
    with dates as (
        SELECT generate_series as date
        FROM generate_series('2020-01-01 00:00'::timestamp,
                '2022-05-15 00:00', '24 hours')
    )
    select bond, date, random() as factor
    from bonds cross join dates;
--rollback DROP TABLE factors;
