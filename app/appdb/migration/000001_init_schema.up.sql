-- Put your ddl queries here.
-- NO NEED TO PUT CREATE DATABASE statement here (assuming we already create the database when

CREATE TABLE user_roles
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO user_roles (name)
values ('Admin'),
       ('Pharmacy Admin'),
       ('Doctor'),
       ('User');

CREATE TABLE users
(
    id           BIGSERIAL PRIMARY KEY,
    email        VARCHAR                   NOT NULL UNIQUE,
    password     VARCHAR                   NOT NULL,
    user_role_id BIGINT                    NOT NULL REFERENCES user_roles (id),
    is_verified  BOOLEAN                   NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at   TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE doctor_profiles
(
    user_id            BIGSERIAL PRIMARY KEY REFERENCES users (id),
    name               VARCHAR                   NOT NULL,
    profile_photo      VARCHAR                   NOT NULL,
    starting_year      INTEGER                   NOT NULL,
    doctor_certificate VARCHAR                   NOT NULL,
    specialization     VARCHAR                   NOT NULL,
    consultation_fee   NUMERIC                   NOT NULL,
    is_online          BOOL                      NOT NULL,
    created_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at         TIMESTAMPTZ DEFAULT NULL
);


CREATE TABLE user_profiles
(
    user_id       BIGSERIAL PRIMARY KEY REFERENCES users (id),
    name          VARCHAR                   NOT NULL,
    profile_photo VARCHAR                   NOT NULL,
    date_of_birth timestamptz               NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at    TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE addresses
(
    id           BIGSERIAL PRIMARY KEY,
    name         VARCHAR                   NOT NULL,
    address      TEXT                      NOT NULL,
    sub_district VARCHAR                   NOT NULL,
    district     VARCHAR                   NOT NULL,
    city         VARCHAR                   NOT NULL,
    province     VARCHAR                   NOT NULL,
    postal_code  VARCHAR                   NOT NULL,
    latitude     VARCHAR                   NOT NULL,
    longitude    VARCHAR                   NOT NULL,
    status       INTEGER                   NOT NULL,
    profile_id   BIGINT                    NOT NULL REFERENCES user_profiles (user_id),
    created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at   TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE shipping_methods
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO shipping_methods (name)
values ('Official Instant'),
       ('Official Same Day'),
       ('Non Official');

CREATE TABLE pharmacies
(
    id                    BIGSERIAL PRIMARY KEY,
    name                  VARCHAR                   NOT NULL,
    address               TEXT                      NOT NULL,
    sub_district          VARCHAR                   NOT NULL,
    district              VARCHAR                   NOT NULL,
    city                  VARCHAR                   NOT NULL,
    province              VARCHAR                   NOT NULL,
    postal_code           VARCHAR                   NOT NULL,
    latitude              VARCHAR                   NOT NULL,
    longitude             VARCHAR                   NOT NULL,
    pharmacist_name       VARCHAR                   NOT NULL,
    pharmacist_license_no VARCHAR                   NOT NULL,
    pharmacist_phone_no   VARCHAR                   NOT NULL,
    operational_hours     VARCHAR                   NOT NULL,
    operational_days      VARCHAR                   NOT NULL,
    pharmacy_admin_id     BIGINT                    NOT NULL REFERENCES users (id),
    created_at            TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at            TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at            TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE pharmacy_shipping_methods
(
    id                 BIGSERIAL PRIMARY KEY,
    pharmacy_id        BIGINT                    NOT NULL REFERENCES pharmacies (id),
    shipping_method_id BIGINT                    NOT NULL REFERENCES shipping_methods (id),
    created_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at         TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE verification_tokens
(
    id         BIGSERIAL PRIMARY KEY,
    token      VARCHAR              NOT NULL,
    is_valid   BOOLEAN DEFAULT TRUE NOT NULL,
    expired_at TIMESTAMPTZ          NOT NULL,
    email      VARCHAR              NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE forgot_password_tokens
(
    id         BIGSERIAL PRIMARY KEY,
    token      VARCHAR                   NOT NULL,
    is_valid   BOOLEAN     DEFAULT TRUE  NOT NULL,
    expired_at TIMESTAMPTZ               NOT NULL,
    user_id    BIGINT                    NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE manufacturers
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO manufacturers (name)
values ('Soho Industri Pharmasi'),
       ('Amarox Pharma Global');

CREATE TABLE product_categories
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO product_categories (name)
values ('Obat'),
       ('Non Obat');

CREATE TABLE drug_classifications
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO drug_classifications (name)
values ('Obat Bebas'),
       ('Obat Keras'),
       ('Obat Bebas Terbatas'),
       ('Non Obat');

CREATE TABLE products
(
    id                     BIGSERIAL PRIMARY KEY,
    name                   VARCHAR                   NOT NULL,
    generic_name           VARCHAR                   NOT NULL,
    content                VARCHAR                   NOT NULL,
    manufacturer_id        BIGINT                    NOT NULL REFERENCES manufacturers (id),
    description            TEXT                      NOT NULL,
    drug_classification_id BIGINT                    NOT NULL REFERENCES drug_classifications (id),
    product_category_id    BIGINT                    NOT NULL REFERENCES product_categories (id),
    drug_form              VARCHAR                   NOT NULL,
    unit_in_pack           VARCHAR                   NOT NULL,
    selling_unit           VARCHAR                   NOT NULL,
    weight                 FLOAT                     NOT NULL, -- in gram
    length                 FLOAT                     NOT NULL, -- in cm
    width                  FLOAT                     NOT NULL, -- in cm
    height                 FLOAT                     NOT NULL, -- in cm
    image                  VARCHAR                   NOT NULL,
    UNIQUE (name, generic_name, content, manufacturer_id),
    created_at             TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at             TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at             TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE cart_items
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT                    NOT NULL REFERENCES users (id),
    product_id BIGINT                    NOT NULL REFERENCES products (id),
    quantity   INT                       NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE pharmacy_products
(
    id          BIGSERIAL PRIMARY KEY,
    pharmacy_id BIGINT                    NOT NULL REFERENCES pharmacies (id),
    product_id  BIGINT                    NOT NULL REFERENCES products (id),
    is_active   BOOL                      NOT NULL,
    price       NUMERIC                   NOT NULL,
    stock       INT                       NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at  TIMESTAMPTZ DEFAULT NULL
);


CREATE TABLE product_stock_mutation_types
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO product_stock_mutation_types (name)
values ('Addition'),
       ('Deduction');


CREATE TABLE product_stock_mutation_request_statuses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO product_stock_mutation_request_statuses (name)
values ('Approved'),
       ('Canceled');


CREATE TABLE product_stock_mutations
(
    id                             BIGSERIAL PRIMARY KEY,
    pharmacy_product_id            BIGINT                    NOT NULL REFERENCES pharmacy_products (id),
    product_stock_mutation_type_id BIGINT                    NOT NULL REFERENCES product_stock_mutation_types (id),
    created_at                     TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at                     TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at                     TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE product_stock_mutation_requests
(
    id                                       BIGSERIAL PRIMARY KEY,
    pharmacy_product_origin_id               BIGINT                    NOT NULL REFERENCES pharmacy_products (id),
    pharmacy_product_dest_id                 BIGINT                    NOT NULL REFERENCES pharmacy_products (id),
    stock                                    INT                       NOT NULL,
    product_stock_mutation_request_status_id BIGINT                    NOT NULL REFERENCES product_stock_mutation_request_statuses (id),
    created_at                               TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at                               TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at                               TIMESTAMPTZ DEFAULT NULL
);


CREATE TABLE consultation_session_statuses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO consultation_session_statuses (name)
values ('Open'),
       ('Close');


CREATE TABLE consultation_sessions
(
    id                             BIGSERIAL PRIMARY KEY,
    user_id                        BIGINT                    NOT NULL REFERENCES users (id),
    doctor_id                      BIGINT                    NOT NULL REFERENCES users (id),
    consultation_session_status_id BIGINT                    NOT NULL REFERENCES consultation_session_statuses (id),
    created_at                     TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at                     TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at                     TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE consultation_messages
(
    id         BIGSERIAL PRIMARY KEY,
    session_id BIGINT                    NOT NULL REFERENCES consultation_sessions (id),
    sender_id  BIGINT                    NOT NULL REFERENCES users (id),
    message    VARCHAR                   NOT NULL,
    attachment VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE prescriptions
(
    id         BIGSERIAL PRIMARY KEY,
    session_id BIGINT                    NOT NULL REFERENCES consultation_sessions (id),
    symptoms   VARCHAR                   NOT NULL,
    diagnosis  VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE prescription_products
(
    id              BIGSERIAL PRIMARY KEY,
    prescription_id BIGINT                    NOT NULL REFERENCES prescriptions (id),
    product_id      BIGINT                    NOT NULL REFERENCES products (id),
    note            VARCHAR                   NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at      TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE sick_leave_forms
(
    id            BIGSERIAL PRIMARY KEY,
    session_id    BIGINT                    NOT NULL REFERENCES prescriptions (id),
    starting_date TIMESTAMPTZ               NOT NULL,
    ending_date   TIMESTAMPTZ               NOT NULL,
    description   TEXT                      NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at    TIMESTAMPTZ DEFAULT NULL
);


CREATE TABLE order_statuses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO order_statuses (name)
values ('Waiting for Payment'),
       ('Paid '),
       ('Waiting for Pharmacy'),
       ('Pharmacy Confirmed'),
       ('Order Packed'),
       ('Ready to Pickup'),
       ('Order Picked'),
       ('Order Sent'),
       ('Received'),
       ('Canceled');


CREATE TABLE payment_methods
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE orders
(
    id                 BIGSERIAL PRIMARY KEY,
    date               TIMESTAMPTZ               NOT NULL,
    user_id            BIGINT                    NOT NULL REFERENCES users (id),
    pharmacy_id        BIGINT                    NOT NULL REFERENCES pharmacies (id),
    payment_method_id  BIGINT                    NOT NULL REFERENCES payment_methods (id),
    no_of_items        INTEGER                   NOT NULL,
    total_payment      NUMERIC                   NOT NULL,
    payment_proof      VARCHAR                   NOT NULL,
    user_address       VARCHAR                   NOT NULL,
    pharmacy_address   VARCHAR                   NOT NULL,
    shipping_method_id BIGINT                    NOT NULL REFERENCES shipping_methods (id),
    shipping_cost      NUMERIC                   NOT NULL,
    created_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at         TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE order_status_logs
(
    id              BIGSERIAL PRIMARY KEY,
    order_id        BIGINT                    NOT NULL REFERENCES orders (id),
    order_status_id BIGINT                    NOT NULL REFERENCES order_statuses (id),
    created_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at      TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE order_details
(
    id           BIGSERIAL PRIMARY KEY,
    order_id     BIGINT                    NOT NULL REFERENCES orders (id),
    product_id   BIGINT                    NOT NULL REFERENCES products (id),
    quantity     INTEGER                   NOT NULL,
    name         VARCHAR                   NOT NULL,
    generic_name VARCHAR                   NOT NULL,
    content      VARCHAR                   NOT NULL,
    description  VARCHAR                   NOT NULL,
    image        VARCHAR                   NOT NULL,
    price        NUMERIC                   NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at   TIMESTAMPTZ DEFAULT NULL
);
