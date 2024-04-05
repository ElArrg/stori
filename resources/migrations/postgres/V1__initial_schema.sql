create table public.account
(
    id        varchar(36) not null
        constraint account_pk
            primary key,
    firstname varchar(50) not null,
    lastname  varchar(50) not null,
    email     varchar(50) not null
        constraint account_email_unique
            unique
);

-- Inserting into account table IF YOU EDIT THE ID, THEN YOU SHOULD UPDATE DOWN FOR notifications_settings
INSERT INTO public.account (id, firstname, lastname, email)
VALUES
    ('acc1', 'James', 'Smith', 'james.smith@example.com'), --edit email
    ('acc2', 'Maria', 'Garcia', 'maria.garcia@example.com'); -- edit email


create table public.transactions
(
    id         varchar(36) not null
        constraint transactions_pk
            primary key,
    account_id varchar(36) not null,
    amount     bigint not null,
    type       varchar(25),
    date       timestamp   not null,
    year       integer not null,
    month      integer not null
);

create index transactions_account_id_idx
    on public.transactions (account_id);

create table if not exists public.templates
(
    id          varchar(36)           not null
        constraint templates_pk
            primary key,
    name        varchar(50)           not null,
    channel     varchar(50)           not null,
    source      varchar(100)          not null,
    source_type varchar(50)           not null,
    active      boolean default false not null,
    constraint templates_name_channel_unique
        unique (name, channel)
);

INSERT INTO public.templates (id, operation, channel, source, source_type, active) VALUES ('tmp1', 'account-summary', 'email', 'd-4ad8d2f9840c444caadb7d53dfabdac7', 'sendgrid', true);


create table public.notifications_settings
(
    id         varchar(36)           not null
        constraint notifications_settings_pk
            primary key,
    account_id varchar(36)           not null,
    channel    varchar(25)           not null,
    enabled    boolean default false not null,
    constraint notifications_settings_account_id_channel_unique
        unique (account_id, channel)
);

INSERT INTO public.notifications_settings (id, account_id, channel, enabled) VALUES ('ns1', 'acc1', 'email', true);
INSERT INTO public.notifications_settings (id, account_id, channel, enabled) VALUES ('ns2', 'acc2', 'email', true);

