CREATE TABLE "profiles" (
    "uid" UUID NOT NULL UNIQUE,
    "mail" TEXT NOT NULL UNIQUE,
    "first_name" TEXT NOT NULL,
    "second_name" TEXT NOT NULL,
    "third_name" TEXT NOT NULL,
    "position" TEXT NOT NULL,
    "department" TEXT NOT NULL,
    PRIMARY KEY("uid")
);

CREATE INDEX "profiles_uid_index" ON "profiles"("uid");