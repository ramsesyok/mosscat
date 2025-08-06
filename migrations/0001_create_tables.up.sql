CREATE TABLE tags (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE oss_components (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    normalized_name TEXT NOT NULL UNIQUE,
    homepage_url TEXT,
    repository_url TEXT,
    description TEXT,
    primary_language TEXT,
    default_usage_role TEXT,
    deprecated BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE oss_component_tags (
    oss_id UUID NOT NULL REFERENCES oss_components(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id),
    PRIMARY KEY (oss_id, tag_id)
);

CREATE TABLE oss_component_layers (
    oss_id UUID NOT NULL REFERENCES oss_components(id) ON DELETE CASCADE,
    layer TEXT NOT NULL,
    PRIMARY KEY (oss_id, layer)
);

CREATE TABLE oss_versions (
    id UUID PRIMARY KEY,
    oss_id UUID NOT NULL REFERENCES oss_components(id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    release_date DATE,
    license_expression_raw TEXT,
    license_concluded TEXT,
    purl TEXT,
    cpe_list TEXT[],
    hash_sha256 TEXT,
    modified BOOLEAN NOT NULL DEFAULT FALSE,
    modification_description TEXT,
    review_status TEXT NOT NULL,
    last_reviewed_at TIMESTAMPTZ,
    scope_status TEXT NOT NULL,
    supplier_type TEXT,
    fork_origin_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_oss_versions_oss_id ON oss_versions (oss_id);

CREATE TABLE projects (
    id UUID PRIMARY KEY,
    project_code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    department TEXT,
    manager TEXT,
    delivery_date DATE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE project_usages (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    oss_id UUID NOT NULL REFERENCES oss_components(id) ON DELETE CASCADE,
    oss_version_id UUID NOT NULL REFERENCES oss_versions(id) ON DELETE CASCADE,
    usage_role TEXT NOT NULL,
    scope_status TEXT NOT NULL,
    inclusion_note TEXT,
    direct_dependency BOOLEAN NOT NULL DEFAULT TRUE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    evaluated_at TIMESTAMPTZ,
    evaluated_by TEXT
);

CREATE INDEX idx_project_usages_project_id ON project_usages (project_id);
CREATE INDEX idx_project_usages_oss_version_id ON project_usages (oss_version_id);

CREATE TABLE scope_policies (
    id UUID PRIMARY KEY,
    runtime_required_default_in_scope BOOLEAN NOT NULL,
    server_env_included BOOLEAN NOT NULL,
    auto_mark_forks_in_scope BOOLEAN NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by TEXT NOT NULL
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    action TEXT NOT NULL,
    user_name TEXT NOT NULL,
    summary TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_type, entity_id);
