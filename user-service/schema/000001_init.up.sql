CREATE TABLE "users"
(
    "id"            uuid         NOT NULL,
    "username"      varchar(255) NOT NULL UNIQUE,
    "email"         varchar(255) NOT NULL,
    "password_hash" varchar(255) NOT NULL,
    CONSTRAINT "users_pk" PRIMARY KEY ("id")
) WITH (
      OIDS= FALSE
    );
