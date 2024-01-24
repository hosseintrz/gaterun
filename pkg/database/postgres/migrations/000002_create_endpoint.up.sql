CREATE TABLE IF NOT EXISTS endpoints(
    id BIGSERIAL INTEGER PRIMARY KEY,
    endpoint VARCHAR(50) NOT NULL,
    method VARCHAR(10) NOT NULL,
    timeout INTERVAL,
    cache_ttl INTERVAL,
    query_strings VARCHAR[],
    target_headers VARCHAR[],
    output_encoding VARCHAR
);