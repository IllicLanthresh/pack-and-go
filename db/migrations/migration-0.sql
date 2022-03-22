CREATE TABLE Trip
(
    id             INTEGER
        CONSTRAINT Trip_pk
            PRIMARY KEY AUTOINCREMENT,
    origin_id      INTEGER NOT NULL,
    destination_id INTEGER NOT NULL,
    dates_bitmask  INTEGER NOT NULL,
    price          INTEGER NOT NULL
);

CREATE INDEX Trip_origin_destination_index
    ON Trip (origin_id, destination_id);

CREATE INDEX Trip_dates_index
    ON Trip (dates_bitmask);
