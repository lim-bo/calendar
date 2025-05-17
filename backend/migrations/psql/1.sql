CREATE TABLE "attachments" (
    "event_id" char(24) NOT NULL,
    "name" text NOT NULL,
);

CREATE INDEX "attachments_eventID_index" on "attachments"("eventID");