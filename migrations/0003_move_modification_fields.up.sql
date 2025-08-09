-- Move modification-related fields from oss_versions to project_usages
ALTER TABLE project_usages ADD COLUMN modified BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE project_usages ADD COLUMN modification_description TEXT;
ALTER TABLE project_usages ADD COLUMN supplier_type TEXT;
ALTER TABLE project_usages ADD COLUMN fork_origin_url TEXT;

-- Note: Data migration would be needed here to move existing data
-- from oss_versions to project_usages if there are existing records

-- Remove modification fields from oss_versions
ALTER TABLE oss_versions DROP COLUMN modified;
ALTER TABLE oss_versions DROP COLUMN modification_description;
ALTER TABLE oss_versions DROP COLUMN supplier_type;
ALTER TABLE oss_versions DROP COLUMN fork_origin_url;