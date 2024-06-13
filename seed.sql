CREATE TABLE Colour
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE SIZE
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE Topic
(
    id   SERIAL,
    name TEXT NOT NULL UNIQUE,
    PRIMARY KEY (id)
);
CREATE TABLE Photos
(
    id                  int8 PRIMARY KEY,
    ImageNo             int8 NOT NULL,
    Width               int8 NOT NULL,
    Height              int8 NOT NULL,
    DominantColor       TEXT    NOT NULL,
    DominantColorOpaque TEXT    NOT NULL,
    URL                 TEXT    NOT NULL,
    IsMain              BOOLEAN NOT NULL,
    HighResolution      TEXT ,
    IsSuspicious        BOOLEAN NOT NULL,
    FullSizeURL         TEXT    NOT NULL,
    IsHidden            BOOLEAN NOT NULL
);


CREATE TABLE Thumbnails
(
    id       SERIAL PRIMARY KEY,
    Type     TEXT    NOT NULL,
    URL      TEXT    NOT NULL,
    Width    int8 NOT NULL,
    Height   int8 NOT NULL,
    photo_id int8 NOT NULL,
    FOREIGN KEY (photo_id) REFERENCES Photos (id)
);

CREATE TABLE Item
(
    id                       int8 PRIMARY KEY,
    title                    TEXT    NOT NULL,
    price                    NUMERIC NOT NULL,
    is_visible               int8 NOT NULL,
    discount                 NUMERIC,
    currency                 TEXT    NOT NULL,
    brand_title              TEXT    NOT NULL,
    user_id                  int8 NOT NULL,
    url                      TEXT    NOT NULL,
    promoted                 BOOLEAN NOT NULL,
    photo_id                 int8 NOT NULL,
    favourite_count          int8 NOT NULL,
    is_favourite             BOOLEAN NOT NULL,
    badge                    TEXT,
    conversion               TEXT,
    service_fee              TEXT    NOT NULL,
    total_item_price         NUMERIC NOT NULL,
    total_item_price_rounded NUMERIC,
    view_count               int8 NOT NULL,
    size_title               TEXT    NOT NULL,
    content_source           TEXT    NOT NULL,
    status                   TEXT,
    icon_badges              TEXT,
    search_tracking_params   TEXT,
    topic_id                 int8,
    FOREIGN KEY (topic_id) REFERENCES Topic (id),
    Foreign Key (photo_id) REFERENCES Photos (id)
);

CREATE TABLE Item_Topic
(
    topic_id int8 NOT NULL,
    item_id  int8 NOT NULL,
    FOREIGN KEY (topic_id) REFERENCES Topic (id),
    FOREIGN KEY (item_id) REFERENCES Item (id)
);

CREATE TABLE Item_Colour
(
    colour_id int8 NOT NULL,
    item_id   int8 NOT NULL,
    FOREIGN KEY (colour_id) REFERENCES Colour (id),
    FOREIGN KEY (item_id) REFERENCES Item (id)
);

CREATE TABLE Item_Size
(
    size_id int8 NOT NULL,
    item_id int8 NOT NULL,
    FOREIGN KEY (size_id) REFERENCES SIZE (id),
    FOREIGN KEY (item_id) REFERENCES Item (id)
);




