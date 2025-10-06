CREATE TYPE ORDER_STATUS AS ENUM (
    'draft',
    'completed',
    'deleted'
);

CREATE TABLE orders(
    guid                UUID PRIMARY KEY        NOT NULL,
    user_guid           UUID                    NOT NULL,
    number              VARCHAR(20)             NOT NULL,
    amount              DECIMAL(14, 2)          NOT NULL,
    status              ORDER_STATUS            NOT NULL,
    created_at          TIMESTAMPTZ             NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ             NOT NULL DEFAULT now()
);

COMMENT ON COLUMN orders.guid           IS 'GUID';
COMMENT ON COLUMN orders.user_guid      IS 'Guid пользователя';
COMMENT ON COLUMN orders.number         IS 'Номер заказа';
COMMENT ON COLUMN orders.amount         IS 'Сумма заказа';
COMMENT ON COLUMN orders.status         IS 'Статус заказа';
COMMENT ON COLUMN orders.created_at     IS 'Дата создания';
COMMENT ON COLUMN orders.updated_at     IS 'Дата обновления';