CREATE TABLE items (
    hash        TEXT NOT NULL,  -- base64_urlencode
    firstseen   INT,            -- unix ctime
    conttype    string,         -- MIME-content type as reported at upload
    raw         BLOB,           -- item's binary data
    PRIMARY KEY(hash)
);

