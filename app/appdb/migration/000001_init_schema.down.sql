-- Put your ddl queries here.
-- NO NEED TO PUT CREATE DATABASE statement here (assuming we already create the database when

DROP TABLE IF EXISTS order_details;
DROP TABLE IF EXISTS order_status_logs;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS order_statuses;
DROP TABLE IF EXISTS sick_leave_forms;
DROP TABLE IF EXISTS prescription_products;
DROP TABLE IF EXISTS prescriptions;
DROP TABLE IF EXISTS consultation_messages;
DROP TABLE IF EXISTS consultation_sessions;
DROP TABLE IF EXISTS consultation_session_statuses;
DROP TABLE IF EXISTS product_stock_mutation_requests;
DROP TABLE IF EXISTS product_stock_mutations;
DROP TABLE IF EXISTS product_stock_mutation_request_statuses;
DROP TABLE IF EXISTS product_stock_mutation_types;
DROP TABLE IF EXISTS pharmacy_products;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS drug_classifications;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS manufacturers;
DROP TABLE IF EXISTS forgot_password_tokens;
DROP TABLE IF EXISTS verification_tokens;
DROP TABLE IF EXISTS pharmacy_shipping_methods;
DROP TABLE IF EXISTS pharmacies;
DROP TABLE IF EXISTS shipping_methods;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS doctor_profiles;
DROP TABLE IF EXISTS doctor_specializations;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS cities;
DROP TABLE IF EXISTS provinces;