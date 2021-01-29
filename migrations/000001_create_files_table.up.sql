CREATE TABLE IF NOT EXISTS files (
    id varchar primary key,
    file_name varchar not null,
    doc_class varchar not null,
    doc_type varchar not null,
    doc_num varchar not null,
    metadata jsonb not null
);