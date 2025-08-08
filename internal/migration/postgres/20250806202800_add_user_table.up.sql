CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
    guid                UUID PRIMARY KEY        NOT NULL,
    name                VARCHAR(255)            NOT NULL,
    occupation          TEXT                    NOT NULL,
    is_deleted          BOOLEAN                 NOT NULL DEFAULT false,
    created_at          TIMESTAMPTZ             NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ             NOT NULL DEFAULT now()
);

COMMENT ON COLUMN users.guid          IS 'GUID пользователя';
COMMENT ON COLUMN users.name          IS 'Имя пользователя';
COMMENT ON COLUMN users.occupation    IS 'Место работы';
COMMENT ON COLUMN users.is_deleted    IS 'Признак удален ли пользователь';
COMMENT ON COLUMN users.created_at    IS 'Дата создания';
COMMENT ON COLUMN users.updated_at    IS 'Дата обновления';

INSERT INTO users (
        guid,
        name,
        occupation,
        is_deleted,
        created_at,
        updated_at
    )
VALUES (
        gen_random_uuid(),
        'Иванова Ариадна Евгеньевна',
        'МУП ДЭС',
        false,
        now(),
        now()
    );