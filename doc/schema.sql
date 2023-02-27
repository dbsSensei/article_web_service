-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-02-25T11:52:44.444Z

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "full_name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) NOT NULL,
  "hashed_password" VARCHAR(255) NOT NULL
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" int NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "articles" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "author_id" INT NOT NULL,
  "content" TEXT NOT NULL,
  "published_at" DATE NOT NULL
);

CREATE TABLE "categories" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL
);

CREATE TABLE "tags" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL
);

CREATE TABLE "article_categories" (
  "article_id" INT NOT NULL,
  "category_id" INT NOT NULL,
  PRIMARY KEY ("article_id", "category_id")
);

CREATE TABLE "article_tags" (
  "article_id" INT NOT NULL,
  "tag_id" INT NOT NULL,
  PRIMARY KEY ("article_id", "tag_id")
);

CREATE TABLE "comments" (
  "id" SERIAL PRIMARY KEY,
  "article_id" INT NOT NULL,
  "user_id" INT NOT NULL,
  "comment_date" DATE NOT NULL,
  "content" TEXT NOT NULL
);

CREATE TABLE "likes" (
  "id" SERIAL PRIMARY KEY,
  "article_id" INT NOT NULL,
  "user_id" INT NOT NULL
);

CREATE TABLE "views" (
  "id" SERIAL PRIMARY KEY,
  "article_id" INT NOT NULL,
  "view_date" DATE NOT NULL
);

CREATE TABLE "images" (
  "id" SERIAL PRIMARY KEY,
  "article_id" INT NOT NULL,
  "url" VARCHAR(255) NOT NULL
);

ALTER TABLE "articles" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "article_categories" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;

ALTER TABLE "article_categories" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE CASCADE;

ALTER TABLE "article_tags" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;

ALTER TABLE "article_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "likes" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;

ALTER TABLE "likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "views" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;

ALTER TABLE "images" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;
