CREATE TABLE "users"
(
    "id"            uuid         NOT NULL,
    "username"      varchar(255) NOT NULL UNIQUE,
    "password_hash" varchar(255) NOT NULL,
    CONSTRAINT "users_pk" PRIMARY KEY ("id")
) WITH (
      OIDS= FALSE
    );
