CREATE TABLE "attachments" (
    "event_id" char(24) NOT NULL,
    "name" text UNIQUE NOT NULL,
);

CREATE INDEX "attachments_eventID_index" on "attachments"("eventID");