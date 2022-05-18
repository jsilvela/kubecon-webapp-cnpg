--liquibase formatted sql

--changeset jaime.silvela:1 labels:kubecon-demo
--comment: let's start off with 50 stocks 
create table stocks as
    select 'stock_' || generate_series as stock
    from generate_series(1, 50);
--rollback DROP TABLE stocks;

--changeset jaime.silvela:2 labels:kubecon-demo
--comment: lets add a bunch of random stock values
create table stock_values as
    with dates as (
        SELECT generate_series as date
        FROM generate_series('2020-01-01 00:00'::timestamp,
                '2022-05-15 00:00', '24 hours')
    )
    select stock, date, random() as stock_value
    from stocks cross join dates;
--rollback DROP TABLE stock_values;
