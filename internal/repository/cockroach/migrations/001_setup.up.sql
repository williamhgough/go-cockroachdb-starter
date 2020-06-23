BEGIN;

CREATE TABLE test (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name STRING not null
);

COMMIT;