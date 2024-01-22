-- Put your ddl queries here.
-- NO NEED TO PUT CREATE DATABASE statement here (assuming we already create the database when starting postgres docker container)

CREATE TABLE provinces
(
    id         BIGINT PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE cities
(
    id          BIGINT PRIMARY KEY        NOT NULL,
    name        VARCHAR                   NOT NULL,
    province_id BIGINT REFERENCES provinces (id),
    created_at  TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at  TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE user_roles
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

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

CREATE TABLE doctor_specializations
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    image      VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE doctor_profiles
(
    user_id                  BIGSERIAL PRIMARY KEY REFERENCES users (id),
    name                     VARCHAR                   NOT NULL,
    profile_photo            VARCHAR                   NOT NULL,
    starting_year            INTEGER                   NOT NULL,
    doctor_certificate       VARCHAR                   NOT NULL,
    doctor_specialization_id BIGINT                    NOT NULL REFERENCES doctor_specializations (id),
    consultation_fee         NUMERIC                   NOT NULL,
    is_online                BOOL                      NOT NULL,
    created_at               TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at               TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at               TIMESTAMPTZ DEFAULT NULL
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
    city         BIGINT                    NOT NULL REFERENCES cities (id),
    province     BIGINT                    NOT NULL REFERENCES provinces (id),
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

CREATE TABLE pharmacies
(
    id                    BIGSERIAL PRIMARY KEY,
    name                  VARCHAR                   NOT NULL,
    address               TEXT                      NOT NULL,
    sub_district          VARCHAR                   NOT NULL,
    district              VARCHAR                   NOT NULL,
    city                  BIGINT                    NOT NULL REFERENCES cities (id),
    province              BIGINT                    NOT NULL REFERENCES provinces (id),
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
    token      VARCHAR                   NOT NULL,
    is_valid   BOOLEAN     DEFAULT TRUE  NOT NULL,
    expired_at TIMESTAMPTZ               NOT NULL,
    email      VARCHAR                   NOT NULL,
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
    image      VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE product_categories
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR UNIQUE            NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE drug_classifications
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

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
    UNIQUE (pharmacy_id, product_id),
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

CREATE TABLE product_stock_mutation_request_statuses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE product_stock_mutations
(
    id                             BIGSERIAL PRIMARY KEY,
    pharmacy_product_id            BIGINT                    NOT NULL REFERENCES pharmacy_products (id),
    product_stock_mutation_type_id BIGINT                    NOT NULL REFERENCES product_stock_mutation_types (id),
    stock                          INT                       NOT NULL,
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

CREATE TABLE payment_methods
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE transaction_statuses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE transactions
(
    id         BIGSERIAL PRIMARY KEY,
    date TIMESTAMPTZ NOT NULL ,
    payment_proof VARCHAR NOT NULL ,
    transaction_status_id  BIGINT                    NOT NULL REFERENCES transaction_statuses (id),
    payment_method_id  BIGINT                    NOT NULL REFERENCES payment_methods (id),
    address       VARCHAR                   NOT NULL ,
    user_id            BIGINT                    NOT NULL REFERENCES users (id),
    total_payment NUMERIC NOT NULL,
    created_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at         TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE orders
(
    id                 BIGSERIAL PRIMARY KEY,
    date               TIMESTAMPTZ               NOT NULL,
    pharmacy_id        BIGINT                    NOT NULL REFERENCES pharmacies (id),
    no_of_items        INTEGER                   NOT NULL,
    pharmacy_address   VARCHAR                   NOT NULL,
    shipping_method_id BIGINT                    NOT NULL REFERENCES shipping_methods (id),
    shipping_cost      NUMERIC                   NOT NULL,
    total_payment      NUMERIC                   NOT NULL,
    transaction_id  BIGINT                    NOT NULL REFERENCES transactions (id),
    created_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT now() NOT NULL,
    deleted_at         TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE order_status_logs
(
    id              BIGSERIAL PRIMARY KEY,
    order_id        BIGINT                    NOT NULL REFERENCES orders (id),
    order_status_id BIGINT                    NOT NULL REFERENCES order_statuses (id),
    is_latest BOOL NOT NULL ,
    description TEXT NOT NULL ,
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


-- INSERT DATA
INSERT INTO user_roles (name)
values ('Admin'),
       ('Pharmacy Admin'),
       ('Doctor'),
       ('User');

INSERT INTO users (email, password, user_role_id, is_verified)
VALUES ('byebyesick@gmail.com', '$2a$04$MYf2/GkfNPUUZUj8zInF.ej7KqSVO3KlJrbNEwkCtCerFXzqbOsDe', 1, true),
       ('yafi.tamfan08@gmail.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 2, true),
       ('tifan@email.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 2, true),
       ('random@email.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 2, true),
       ('wasikamin4@gmail.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 3, true),
       ('lumbanraja.boy@gmail.com', '$2a$04$s0eeWG0MEJ6b.ffuGsibcuhqHZJyIxHbb5Cc/EckWD2GY9ZnvUj9S', 4, true);

INSERT INTO doctor_specializations (name, image)
values ('General Practitioners',
        'https://byebyesick-bucket.irfancen.com/doctor_specializations/doctor-specs.jpg'),
       ('Pediatric Specialist',
        'https://byebyesick-bucket.irfancen.com/doctor_specializations/doctor-specs.jpg');

INSERT INTO doctor_profiles (user_id, name, profile_photo, starting_year, doctor_certificate, doctor_specialization_id, consultation_fee, is_online)
VALUES (5, 'dokter wasik', '', 2021, '', 1, 10, false);

INSERT INTO user_profiles(user_id, name, profile_photo, date_of_birth)
VALUES (6, 'lumban boy', '', '2000-11-25');

INSERT INTO manufacturers (name, image)
values ('Soho Industri Pharmasi', 'https://byebyesick-bucket.irfancen.com/doctor_specializations/soho.png'),
       ('Amarox Pharma Global', 'https://byebyesick-bucket.irfancen.com/doctor_specializations/amarox.jpeg');

INSERT INTO shipping_methods (name)
values ('Official Instant'),
       ('Official Same Day'),
       ('Non Official');

INSERT INTO product_categories (name)
values ('Obat'),
       ('Non Obat');

INSERT INTO drug_classifications (name)
values ('Obat Bebas'),
       ('Obat Keras'),
       ('Obat Bebas Terbatas'),
       ('Non Obat');

INSERT INTO product_stock_mutation_types (name)
values ('Addition'),
       ('Deduction');

INSERT INTO product_stock_mutation_request_statuses (name)
values ('Pending'),
       ('Accepted'),
       ('Rejected');

INSERT INTO consultation_session_statuses (name)
values ('Open'),
       ('Close');

INSERT INTO payment_methods (name)
values ('Bank Transfer');

INSERT INTO order_statuses (name)
values ('Waiting for Pharmacy'),
       ('Processed'),
       ('Sent'),
       ('Order Confirmed'),
        ('Canceled by Pharmacy'),
        ('Canceled by User');


INSERT INTO transaction_statuses (name)
values ('Unpaid'),
       ('Waiting for Confirmation'),
       ('Payment Rejected'),
       ('Paid'),
        ('Canceled');

INSERT INTO provinces (id, name)
VALUES (1, 'Bali');
INSERT INTO provinces (id, name)
VALUES (2, 'Bangka Belitung');
INSERT INTO provinces (id, name)
VALUES (3, 'Banten');
INSERT INTO provinces (id, name)
VALUES (4, 'Bengkulu');
INSERT INTO provinces (id, name)
VALUES (5, 'DI Yogyakarta');
INSERT INTO provinces (id, name)
VALUES (6, 'DKI Jakarta');
INSERT INTO provinces (id, name)
VALUES (7, 'Gorontalo');
INSERT INTO provinces (id, name)
VALUES (8, 'Jambi');
INSERT INTO provinces (id, name)
VALUES (9, 'Jawa Barat');
INSERT INTO provinces (id, name)
VALUES (10, 'Jawa Tengah');
INSERT INTO provinces (id, name)
VALUES (11, 'Jawa Timur');
INSERT INTO provinces (id, name)
VALUES (12, 'Kalimantan Barat');
INSERT INTO provinces (id, name)
VALUES (13, 'Kalimantan Selatan');
INSERT INTO provinces (id, name)
VALUES (14, 'Kalimantan Tengah');
INSERT INTO provinces (id, name)
VALUES (15, 'Kalimantan Timur');
INSERT INTO provinces (id, name)
VALUES (16, 'Kalimantan Utara');
INSERT INTO provinces (id, name)
VALUES (17, 'Kepulauan Riau');
INSERT INTO provinces (id, name)
VALUES (18, 'Lampung');
INSERT INTO provinces (id, name)
VALUES (19, 'Maluku');
INSERT INTO provinces (id, name)
VALUES (20, 'Maluku Utara');
INSERT INTO provinces (id, name)
VALUES (21, 'Nanggroe Aceh Darussalam (NAD)');
INSERT INTO provinces (id, name)
VALUES (22, 'Nusa Tenggara Barat (NTB)');
INSERT INTO provinces (id, name)
VALUES (23, 'Nusa Tenggara Timur (NTT)');
INSERT INTO provinces (id, name)
VALUES (24, 'Papua');
INSERT INTO provinces (id, name)
VALUES (25, 'Papua Barat');
INSERT INTO provinces (id, name)
VALUES (26, 'Riau');
INSERT INTO provinces (id, name)
VALUES (27, 'Sulawesi Barat');
INSERT INTO provinces (id, name)
VALUES (28, 'Sulawesi Selatan');
INSERT INTO provinces (id, name)
VALUES (29, 'Sulawesi Tengah');
INSERT INTO provinces (id, name)
VALUES (30, 'Sulawesi Tenggara');
INSERT INTO provinces (id, name)
VALUES (31, 'Sulawesi Utara');
INSERT INTO provinces (id, name)
VALUES (32, 'Sumatera Barat');
INSERT INTO provinces (id, name)
VALUES (33, 'Sumatera Selatan');
INSERT INTO provinces (id, name)
VALUES (34, 'Sumatera Utara');

INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Badung', 17, 1);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bangli', 32, 1);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Buleleng', 94, 1);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Denpasar', 114, 1);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Gianyar', 128, 1);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Jembrana', 161, 1);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Karangasem', 170, 1);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Klungkung', 197, 1);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tabanan', 447, 1);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bangka', 27, 2);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bangka Barat', 28, 2);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bangka Selatan', 29, 2);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bangka Tengah', 30, 2);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Belitung', 56, 2);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Belitung Timur', 57, 2);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Pangkal Pinang', 334, 2);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Cilegon', 106, 3);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lebak', 232, 3);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pandeglang', 331, 3);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Serang', 402, 3);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Serang', 403, 3);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tangerang', 455, 3);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tangerang', 456, 3);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tangerang Selatan', 457, 3);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bengkulu', 62, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bengkulu Selatan', 63, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bengkulu Tengah', 64, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bengkulu Utara', 65, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kaur', 175, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepahiang', 183, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lebong', 233, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Muko Muko', 294, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Rejang Lebong', 379, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Seluma', 397, 4);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bantul', 39, 5);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Gunung Kidul', 135, 5);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kulon Progo', 210, 5);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sleman', 419, 5);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Yogyakarta', 501, 5);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Jakarta Barat', 151, 6);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Jakarta Pusat', 152, 6);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Jakarta Selatan', 153, 6);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Jakarta Timur', 154, 6);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Jakarta Utara', 155, 6);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Seribu', 189, 6);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Boalemo', 77, 7);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bone Bolango', 88, 7);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Gorontalo', 129, 7);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Gorontalo', 130, 7);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Gorontalo Utara', 131, 7);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pohuwato', 361, 7);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Batang Hari', 50, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bungo', 97, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Jambi', 156, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kerinci', 194, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Merangin', 280, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Muaro Jambi', 293, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sarolangun', 393, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Sungaipenuh', 442, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tanjung Jabung Barat', 460, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tanjung Jabung Timur', 461, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tebo', 471, 8);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bandung', 22, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bandung', 23, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bandung Barat', 24, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Banjar', 34, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bekasi', 54, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bekasi', 55, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bogor', 78, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bogor', 79, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ciamis', 103, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Cianjur', 104, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Cimahi', 107, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Cirebon', 108, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Cirebon', 109, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Depok', 115, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Garut', 126, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Indramayu', 149, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Karawang', 171, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kuningan', 211, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Majalengka', 252, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pangandaran', 332, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Purwakarta', 376, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Subang', 428, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sukabumi', 430, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Sukabumi', 431, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sumedang', 440, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tasikmalaya', 468, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tasikmalaya', 469, 9);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Banjarnegara', 37, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Banyumas', 41, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Batang', 49, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Blora', 76, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Boyolali', 91, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Brebes', 92, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Cilacap', 105, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Demak', 113, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Grobogan', 134, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Jepara', 163, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Karanganyar', 169, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kebumen', 177, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kendal', 181, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Klaten', 196, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kudus', 209, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Magelang', 249, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Magelang', 250, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pati', 344, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pekalongan', 348, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Pekalongan', 349, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pemalang', 352, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Purbalingga', 375, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Purworejo', 377, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Rembang', 380, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Salatiga', 386, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Semarang', 398, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Semarang', 399, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sragen', 427, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sukoharjo', 433, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Surakarta (Solo)', 445, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tegal', 472, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tegal', 473, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Temanggung', 476, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Wonogiri', 497, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Wonosobo', 498, 10);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bangkalan', 31, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Banyuwangi', 42, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Batu', 51, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Blitar', 74, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Blitar', 75, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bojonegoro', 80, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bondowoso', 86, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Gresik', 133, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Jember', 160, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Jombang', 164, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kediri', 178, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Kediri', 179, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lamongan', 222, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lumajang', 243, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Madiun', 247, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Madiun', 248, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Magetan', 251, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Malang', 256, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Malang', 255, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mojokerto', 289, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Mojokerto', 290, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nganjuk', 305, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ngawi', 306, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pacitan', 317, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pamekasan', 330, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pasuruan', 342, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Pasuruan', 343, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ponorogo', 363, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Probolinggo', 369, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Probolinggo', 370, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sampang', 390, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sidoarjo', 409, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Situbondo', 418, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sumenep', 441, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Surabaya', 444, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Trenggalek', 487, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tuban', 489, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tulungagung', 492, 11);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bengkayang', 61, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kapuas Hulu', 168, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kayong Utara', 176, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ketapang', 195, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kubu Raya', 208, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Landak', 228, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Melawi', 279, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pontianak', 364, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Pontianak', 365, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sambas', 388, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sanggau', 391, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sekadau', 395, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Singkawang', 415, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sintang', 417, 12);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Balangan', 18, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Banjar', 33, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Banjarbaru', 35, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Banjarmasin', 36, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Barito Kuala', 43, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Hulu Sungai Selatan', 143, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Hulu Sungai Tengah', 144, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Hulu Sungai Utara', 145, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kotabaru', 203, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tabalong', 446, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tanah Bumbu', 452, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tanah Laut', 454, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tapin', 466, 13);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Barito Selatan', 44, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Barito Timur', 45, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Barito Utara', 46, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Gunung Mas', 136, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kapuas', 167, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Katingan', 174, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kotawaringin Barat', 205, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kotawaringin Timur', 206, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lamandau', 221, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Murung Raya', 296, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Palangka Raya', 326, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pulang Pisau', 371, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Seruyan', 405, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sukamara', 432, 14);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Balikpapan', 19, 15);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Berau', 66, 15);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bontang', 89, 15);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kutai Barat', 214, 15);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kutai Kartanegara', 215, 15);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kutai Timur', 216, 15);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Paser', 341, 15);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Penajam Paser Utara', 354, 15);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Samarinda', 387, 15);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bulungan (Bulongan)', 96, 16);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Malinau', 257, 16);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nunukan', 311, 16);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tana Tidung', 450, 16);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tarakan', 467, 16);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Batam', 48, 17);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bintan', 71, 17);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Karimun', 172, 17);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Anambas', 184, 17);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lingga', 237, 17);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Natuna', 302, 17);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tanjung Pinang', 462, 17);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bandar Lampung', 21, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lampung Barat', 223, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lampung Selatan', 224, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lampung Tengah', 225, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lampung Timur', 226, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lampung Utara', 227, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mesuji', 282, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Metro', 283, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pesawaran', 355, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pesisir Barat', 356, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pringsewu', 368, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tanggamus', 458, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tulang Bawang', 490, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tulang Bawang Barat', 491, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Way Kanan', 496, 18);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Ambon', 14, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Buru', 99, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Buru Selatan', 100, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Aru', 185, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Maluku Barat Daya', 258, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Maluku Tengah', 259, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Maluku Tenggara', 260, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Maluku Tenggara Barat', 261, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Seram Bagian Barat', 400, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Seram Bagian Timur', 401, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tual', 488, 19);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Halmahera Barat', 138, 20);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Halmahera Selatan', 139, 20);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Halmahera Tengah', 140, 20);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Halmahera Timur', 141, 20);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Halmahera Utara', 142, 20);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Sula', 191, 20);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pulau Morotai', 372, 20);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Ternate', 477, 20);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tidore Kepulauan', 478, 20);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Barat', 1, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Barat Daya', 2, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Besar', 3, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Jaya', 4, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Selatan', 5, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Singkil', 6, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Tamiang', 7, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Tengah', 8, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Tenggara', 9, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Timur', 10, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Aceh Utara', 11, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Banda Aceh', 20, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bener Meriah', 59, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bireuen', 72, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Gayo Lues', 127, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Langsa', 230, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Lhokseumawe', 235, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nagan Raya', 300, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pidie', 358, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pidie Jaya', 359, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Sabang', 384, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Simeulue', 414, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Subulussalam', 429, 21);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bima', 68, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bima', 69, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Dompu', 118, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lombok Barat', 238, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lombok Tengah', 239, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lombok Timur', 240, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lombok Utara', 241, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Mataram', 276, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sumbawa', 438, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sumbawa Barat', 439, 22);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Alor', 13, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Belu', 58, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ende', 122, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Flores Timur', 125, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kupang', 212, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Kupang', 213, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lembata', 234, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Manggarai', 269, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Manggarai Barat', 270, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Manggarai Timur', 271, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nagekeo', 301, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ngada', 304, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Rote Ndao', 383, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sabu Raijua', 385, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sikka', 412, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sumba Barat', 434, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sumba Barat Daya', 435, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sumba Tengah', 436, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sumba Timur', 437, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Timor Tengah Selatan', 479, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Timor Tengah Utara', 480, 23);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Asmat', 16, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Biak Numfor', 67, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Boven Digoel', 90, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Deiyai (Deliyai)', 111, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Dogiyai', 117, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Intan Jaya', 150, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Jayapura', 157, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Jayapura', 158, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Jayawijaya', 159, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Keerom', 180, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Yapen (Yapen Waropen)', 193, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lanny Jaya', 231, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mamberamo Raya', 263, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mamberamo Tengah', 264, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mappi', 274, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Merauke', 281, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mimika', 284, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nabire', 299, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nduga', 303, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Paniai', 335, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pegunungan Bintang', 347, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Puncak', 373, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Puncak Jaya', 374, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sarmi', 392, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Supiori', 443, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tolikara', 484, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Waropen', 495, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Yahukimo', 499, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Yalimo', 500, 24);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Fakfak', 124, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kaimana', 165, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Manokwari', 272, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Manokwari Selatan', 273, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Maybrat', 277, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pegunungan Arfak', 346, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Raja Ampat', 378, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sorong', 424, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Sorong', 425, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sorong Selatan', 426, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tambrauw', 449, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Teluk Bintuni', 474, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Teluk Wondama', 475, 25);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bengkalis', 60, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Dumai', 120, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Indragiri Hilir', 147, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Indragiri Hulu', 148, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kampar', 166, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Meranti', 187, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kuantan Singingi', 207, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Pekanbaru', 350, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pelalawan', 351, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Rokan Hilir', 381, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Rokan Hulu', 382, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Siak', 406, 26);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Majene', 253, 27);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mamasa', 262, 27);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mamuju', 265, 27);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mamuju Utara', 266, 27);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Polewali Mandar', 362, 27);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bantaeng', 38, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Barru', 47, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bone', 87, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bulukumba', 95, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Enrekang', 123, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Gowa', 132, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Jeneponto', 162, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Luwu', 244, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Luwu Timur', 245, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Luwu Utara', 246, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Makassar', 254, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Maros', 275, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Palopo', 328, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pangkajene Kepulauan', 333, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Parepare', 336, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pinrang', 360, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Selayar (Kepulauan Selayar)', 396, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sidenreng Rappang/Rapang', 408, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sinjai', 416, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Soppeng', 423, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Takalar', 448, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tana Toraja', 451, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Toraja Utara', 486, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Wajo', 493, 28);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Banggai', 25, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Banggai Kepulauan', 26, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Buol', 98, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Donggala', 119, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Morowali', 291, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Palu', 329, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Parigi Moutong', 338, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Poso', 366, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sigi', 410, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tojo Una-Una', 482, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Toli-Toli', 483, 29);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bau-Bau', 53, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bombana', 85, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Buton', 101, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Buton Utara', 102, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Kendari', 182, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kolaka', 198, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kolaka Utara', 199, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Konawe', 200, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Konawe Selatan', 201, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Konawe Utara', 202, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Muna', 295, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Wakatobi', 494, 30);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bitung', 73, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bolaang Mongondow (Bolmong)', 81, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bolaang Mongondow Selatan', 82, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bolaang Mongondow Timur', 83, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Bolaang Mongondow Utara', 84, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Sangihe', 188, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Siau Tagulandang Biaro (Sitaro)', 190, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Talaud', 192, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Kotamobagu', 204, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Manado', 267, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Minahasa', 285, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Minahasa Selatan', 286, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Minahasa Tenggara', 287, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Minahasa Utara', 288, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tomohon', 485, 31);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Agam', 12, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Bukittinggi', 93, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Dharmasraya', 116, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Kepulauan Mentawai', 186, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lima Puluh Koto/Kota', 236, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Padang', 318, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Padang Panjang', 321, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Padang Pariaman', 322, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Pariaman', 337, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pasaman', 339, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pasaman Barat', 340, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Payakumbuh', 345, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pesisir Selatan', 357, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Sawah Lunto', 394, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Sijunjung (Sawah Lunto Sijunjung)', 411, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Solok', 420, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Solok', 421, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Solok Selatan', 422, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tanah Datar', 453, 32);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Banyuasin', 40, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Empat Lawang', 121, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Lahat', 220, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Lubuk Linggau', 242, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Muara Enim', 292, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Musi Banyuasin', 297, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Musi Rawas', 298, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ogan Ilir', 312, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ogan Komering Ilir', 313, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ogan Komering Ulu', 314, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ogan Komering Ulu Selatan', 315, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Ogan Komering Ulu Timur', 316, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Pagar Alam', 324, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Palembang', 327, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Prabumulih', 367, 33);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Asahan', 15, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Batu Bara', 52, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Binjai', 70, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Dairi', 110, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Deli Serdang', 112, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Gunungsitoli', 137, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Humbang Hasundutan', 146, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Karo', 173, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Labuhan Batu', 217, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Labuhan Batu Selatan', 218, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Labuhan Batu Utara', 219, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Langkat', 229, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Mandailing Natal', 268, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Medan', 278, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nias', 307, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nias Barat', 308, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nias Selatan', 309, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Nias Utara', 310, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Padang Lawas', 319, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Padang Lawas Utara', 320, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Padang Sidempuan', 323, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Pakpak Bharat', 325, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Pematang Siantar', 353, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Samosir', 389, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Serdang Bedagai', 404, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Sibolga', 407, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Simalungun', 413, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tanjung Balai', 459, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tapanuli Selatan', 463, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tapanuli Tengah', 464, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Tapanuli Utara', 465, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kota Tebing Tinggi', 470, 34);
INSERT INTO cities (name, id, province_id)
VALUES ('Kab. Toba Samosir', 481, 34);

INSERT INTO products(name, generic_name, content, manufacturer_id, description, drug_classification_id,
                     product_category_id, drug_form, unit_in_pack, selling_unit, weight, length, width, height, image)
VALUES ('Panadol 500 mg 10 Kaplet', 'Panadol', 'Paracetamol', 1, 'Obat sakit kepala', 1,
        1, 'Blister', 10, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/panadol.jpg'),
       ('Saridon 4 Tablet', 'Saridon', 'Paracetamol 250 mg, propyphenazone 150 mg, caffeine 50 mg', 1,
        'Obat sakit kepala', 1,
        1, 'Tablet', 4, 'Strip', 0.05, 50, 50, 50,
        'https://byebyesick-bucket.irfancen.com/products/saridon.jpg');

INSERT INTO pharmacies(name, address, sub_district, district, city, province, postal_code, latitude, longitude,
                       pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days,
                       pharmacy_admin_id)
VALUES ('Kimia Farma Kuningan', 'Jalan Gatau', 'Kuningan', 'Setia Budi', 153, 6, '12950', '-6.230060', '106.827363',
        'M. Irfan Junaidi',
        '69696969', '08123456789', '0-20', 'mon,tue,wed,thu,fri', 2),
       ('Kimia Farma Pasar Minggu', 'Jalan Jalan', 'Ragunan', 'Pasar Minggu', 153, 6, '12560', '-6.290963',
        '106.817317',
        'M. Yafi Al Hakim',
        '42042042', '08998239082', '0-22', 'mon,tue,wed,thu,fri,sat,sun', 2),
       ('Apotek Sinar Jaya', 'Jalan Jaya', 'Merdeka', 'Medan Baru', 278, 34, '20222', '3.576816', '98.659355',
        'Victor Castor',
        '10090900', '0892308932', '0-18', 'mon,tue,wed,thu,fri', 4);

INSERT INTO pharmacy_products(pharmacy_id, product_id, is_active, price, stock)
VALUES (1, 1, true, '12000', 100),
       (1, 2, true, '5000', 95),
       (3, 1, true, '15000', 87);

INSERT INTO pharmacy_shipping_methods(pharmacy_id, shipping_method_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (2, 2);

INSERT INTO addresses(name, address, sub_district, district, city, province, postal_code, latitude, longitude, status, profile_id)
values ('rumah', 'jl kripat', 'ciangasna', 'gunung putri', 78,9, '16968', '-6.354846', '106.952082', 1, 6);

INSERT INTO transactions(date, payment_proof, transaction_status_id, payment_method_id, address, user_id, total_payment)
values (now(), '',1,1,'jl kripat',6,20000);

INSERT INTO orders(date, pharmacy_id, no_of_items, pharmacy_address, shipping_method_id, shipping_cost, total_payment, transaction_id)
values (now(), 1, 2, 'Jalan Gatau', 1, 5000, 10000, 1),
       (now(), 2, 1, 'Jalan Jalan', 1, 1000, 10000, 1),
       (now(), 3, 1, 'Jalan Jaya', 1, 0, 0, 1);

INSERT INTO order_details(order_id, product_id, quantity, name, generic_name, content, description, image, price)
values (1, 1, 1, 'Panadol 500 mg 10 Kaplet', 'Panadol', 'Paracetamol', 'Obat sakit kepala', '', 2500),
       (1, 2, 1, 'Saridon 4 Tablet', 'Saridon', 'Paracetamol 250 mg, propyphenazone 150 mg, caffeine 50 mg', 'Obat sakit kepala', '', 2500),
       (2, 2, 3, 'Saridon 4 Tablet', 'Saridon', 'Paracetamol 250 mg, propyphenazone 150 mg, caffeine 50 mg', 'Obat sakit kepala', '', 3000),
        (3, 2, 3, 'Saridon 4 Tablet', 'Saridon', 'Paracetamol 250 mg, propyphenazone 150 mg, caffeine 50 mg', 'Obat sakit kepala', '', 0);


INSERT INTO order_status_logs(order_id, order_status_id, is_latest, description)
values (1,1,false,''), (1,2,true,''), (2,1,true,''),(3,1,true,'');

-- CREATE FUNCTIONS --

-- SQL code to create a function that calculates the distance in kilometers
-- using the haversine formula

-- Define the radius of the Earth in kilometers
CREATE
OR REPLACE FUNCTION earth_radius()
    RETURNS DECIMAL AS
$$
BEGIN
RETURN 6371;
END;
$$
LANGUAGE plpgsql;

-- Create a function that takes two pairs of coordinates as input
-- and returns the distance between them as output
CREATE
OR REPLACE FUNCTION distance(lat1 VARCHAR, lon1 VARCHAR, lat2 VARCHAR, lon2 VARCHAR) RETURNS DECIMAL AS
$$
DECLARE
    -- Convert degrees to radians
radLat1 DECIMAL := RADIANS(lat1::DECIMAL);
    radLon1
DECIMAL := RADIANS(lon1::DECIMAL);
    radLat2
DECIMAL := RADIANS(lat2::DECIMAL);
    radLon2
DECIMAL := RADIANS(lon2::DECIMAL);
    -- Calculate the difference between the coordinates
    dLat
DECIMAL := radLat2 - radLat1;
    dLon
DECIMAL := radLon2 - radLon1;
    -- Apply the haversine formula
    a
DECIMAL := SIN(dLat / 2) ^ 2 + COS(radLat1) * COS(radLat2) * SIN(dLon / 2) ^ 2;
    c
DECIMAL := 2 * ATAN2(SQRT(a), SQRT(1 - a));
BEGIN
    -- Calculate the distance in kilometers
RETURN earth_radius() * c;
END
$$
LANGUAGE plpgsql;
