CREATE TABLE IF NOT EXISTS "public"."ingestion_data" (
  "uuid" VARCHAR(36) PRIMARY KEY,
  "source_type" TEXT,
  "source_id" TEXT,
  "data" TEXT,
  "data_id" VARCHAR(255),
  "metadata" JSONB,
  "metadata_type" VARCHAR(255),
  "lang" VARCHAR(10),
  "tenant_id" VARCHAR(255),
  "created_at" TIMESTAMP WITH TIME ZONE,
  "updated_at" TIMESTAMP WITH TIME ZONE
);

-- IDEALLY, we would like to have indexes on the following columns:
-- (source_type, source_id, data_id)
-- (source_type)
-- (source_id)
-- (tenant_id)
