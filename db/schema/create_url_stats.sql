CREATE TABLE IF NOT EXISTS url_stats (
    time timestamp NOT NULL,
    url text NOT NULL,
    region text NOT NULL,
    time_namelookup double precision NOT NULL,
    time_connect double precision NOT NULL,
    time_appconnect double precision NOT NULL,
    time_pretransfer double precision NOT NULL,
    time_starttransfer double precision NOT NULL,
    time_total double precision NOT NULL
);

SELECT create_hypertable('url_stats', 'time'); 