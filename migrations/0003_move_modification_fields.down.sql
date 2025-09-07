-- Rollback: Move modification-related fields back from project_usages to oss_versions
ALTER TABLE oss_versions ADD COLUMN modified BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE oss_versions ADD COLUMN modification_description TEXT;
ALTER TABLE oss_versions ADD COLUMN supplier_type TEXT;
ALTER TABLE oss_versions ADD COLUMN fork_origin_url TEXT;

-- Note: Data migration would be needed here to move data back
-- from project_usages to oss_versions during rollback

-- Remove modification fields from project_usages
ALTER TABLE project_usages DROP COLUMN modified;
ALTER TABLE project_usages DROP COLUMN modification_description;
ALTER TABLE project_usages DROP COLUMN supplier_type;
ALTER TABLE project_usages DROP COLUMN fork_origin_url;