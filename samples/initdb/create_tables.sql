DROP TABLE IF EXISTS "cases";
DROP SEQUENCE IF EXISTS cases_id_seq;
CREATE SEQUENCE cases_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."cases" (
    "id" integer DEFAULT nextval('cases_id_seq') NOT NULL,
    "start_time_stamp" timestamp,
    "end_time_stamp" timestamp,
    "customer_id" integer NOT NULL,
    "store_id" integer NOT NULL,
    "created_on" timestamp NOT NULL,
    "modified_on" timestamp NOT NULL,
    CONSTRAINT "cases_id" PRIMARY KEY ("id")
) WITH (oids = false);

DROP TABLE IF EXISTS "customers";
DROP SEQUENCE IF EXISTS customers_id_seq;
CREATE SEQUENCE customers_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."customers" (
    "id" integer DEFAULT nextval('customers_id_seq') NOT NULL,
    "first_name" text,
    "last_name" text,
    "age" integer,
    "email" text,
    "created_on" timestamp NOT NULL,
    "modified_on" timestamp NOT NULL,
    CONSTRAINT "costumers_id" PRIMARY KEY ("id")
) WITH (oids = false);


DROP TABLE IF EXISTS "stores";
DROP SEQUENCE IF EXISTS stores_id_seq;
CREATE SEQUENCE stores_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."stores" (
    "id" integer DEFAULT nextval('stores_id_seq') NOT NULL,
    "name" text,
    "address" text,
    "created_on" timestamp NOT NULL,
    "modified_on" timestamp NOT NULL,
    CONSTRAINT "stores_id" PRIMARY KEY ("id")
) WITH (oids = false);


ALTER TABLE ONLY "public"."cases" ADD CONSTRAINT "cases_customers_id_fkey" FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY "public"."cases" ADD CONSTRAINT "cases_stores_id_fkey" FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE NOT DEFERRABLE;