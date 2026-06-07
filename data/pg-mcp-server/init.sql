-- ══════════════════════════════════════════════════════════════════════════════
--  pg-mcp-server — Schema d'initialisation
--  Version: 2.0.0
--  Exécuté automatiquement par `lfa setup` après installation PostgreSQL.
-- ══════════════════════════════════════════════════════════════════════════════

-- ─── Extension UUID (optionnel mais recommandé) ─────────────────────────────
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ─── Projets ────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS projects (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    root_path   TEXT,
    tech_stack  JSONB DEFAULT '[]'::jsonb,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- ─── Entités (graphe de connaissance) ───────────────────────────────────────
CREATE TABLE IF NOT EXISTS entities (
    id           SERIAL PRIMARY KEY,
    project_id   INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    entity_type  VARCHAR(50) NOT NULL,
    observations TEXT[] DEFAULT '{}',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(project_id, name)
);

-- ─── Relations entre entités ────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS relations (
    id                SERIAL PRIMARY KEY,
    project_id        INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    source_entity_id  INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    target_entity_id  INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    relation_type     VARCHAR(50) NOT NULL DEFAULT 'related_to',
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(project_id, source_entity_id, target_entity_id, relation_type)
);

-- ─── Tâches ─────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS tasks (
    id          SERIAL PRIMARY KEY,
    project_id  INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title       VARCHAR(500) NOT NULL,
    description TEXT,
    priority    VARCHAR(10) NOT NULL DEFAULT 'medium'
                CHECK (priority IN ('critical','high','medium','low')),
    status      VARCHAR(20) NOT NULL DEFAULT 'pending'
                CHECK (status IN ('pending','in_progress','completed','cancelled','blocked')),
    agent_type  VARCHAR(100),
    tags        TEXT[] DEFAULT '{}',
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- ─── ADR (Architecture Decision Records) ────────────────────────────────────
CREATE TABLE IF NOT EXISTS adrs (
    id            SERIAL PRIMARY KEY,
    project_id    INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title         VARCHAR(500) NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'proposed'
                  CHECK (status IN ('proposed','accepted','deprecated','superseded')),
    context       TEXT NOT NULL,
    decision      TEXT NOT NULL,
    consequences  TEXT NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

-- ─── Design Tokens ──────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS design_tokens (
    id          SERIAL PRIMARY KEY,
    project_id  INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    category    VARCHAR(30) NOT NULL
                CHECK (category IN ('color','typography','spacing','shadow','radius','opacity','font','breakpoint')),
    name        VARCHAR(255) NOT NULL,
    value       TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(project_id, category, name)
);

-- ─── Conversations journalisées ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS conversations (
    id          SERIAL PRIMARY KEY,
    project_id  INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title       VARCHAR(500) NOT NULL,
    summary     TEXT NOT NULL,
    agent_type  VARCHAR(100),
    tags        TEXT[] DEFAULT '{}',
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- ─── Code Findings ──────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS code_findings (
    id            SERIAL PRIMARY KEY,
    project_id    INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    finding_type  VARCHAR(20) NOT NULL
                  CHECK (finding_type IN ('bug','smell','security','performance','style','duplication','coverage')),
    severity      VARCHAR(10) NOT NULL
                  CHECK (severity IN ('critical','high','medium','low','info')),
    title         VARCHAR(500) NOT NULL,
    description   TEXT,
    file_path     TEXT,
    line_start    INTEGER,
    suggestion    TEXT,
    status        VARCHAR(10) NOT NULL DEFAULT 'open'
                  CHECK (status IN ('open','closed','wontfix')),
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

-- ─── Index ──────────────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_entities_project   ON entities(project_id);
CREATE INDEX IF NOT EXISTS idx_entities_name      ON entities(name);
CREATE INDEX IF NOT EXISTS idx_entities_type      ON entities(entity_type);
CREATE INDEX IF NOT EXISTS idx_relations_project  ON relations(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_project      ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status       ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_adrs_project       ON adrs(project_id);
CREATE INDEX IF NOT EXISTS idx_tokens_project     ON design_tokens(project_id);
CREATE INDEX IF NOT EXISTS idx_conversations_project ON conversations(project_id);
CREATE INDEX IF NOT EXISTS idx_findings_project   ON code_findings(project_id);
CREATE INDEX IF NOT EXISTS idx_findings_status    ON code_findings(status);

-- ─── Fonctions utilitaires ──────────────────────────────────────────────────

-- search_entities: recherche plein texte dans les entités et leurs observations
CREATE OR REPLACE FUNCTION search_entities(
    p_project_id UUID,
    p_query      TEXT,
    p_limit      INTEGER DEFAULT 10
) RETURNS TABLE(
    id            INTEGER,
    name          VARCHAR,
    entity_type   VARCHAR,
    observations  TEXT[],
    relevance     REAL
) LANGUAGE SQL STABLE AS $$
    SELECT id, name, entity_type, observations,
           ts_rank(
               to_tsvector('french', name || ' ' || array_to_string(observations, ' ')),
               plainto_tsquery('french', p_query)
           ) AS relevance
    FROM entities
    WHERE project_id = p_project_id
      AND to_tsvector('french', name || ' ' || array_to_string(observations, ' '))
          @@ plainto_tsquery('french', p_query)
    ORDER BY relevance DESC
    LIMIT p_limit;
$$;

-- project_dashboard: tableau de bord complet d'un projet
CREATE OR REPLACE FUNCTION project_dashboard(p_project_id UUID)
RETURNS JSONB LANGUAGE SQL STABLE AS $$
    SELECT jsonb_build_object(
        'project',       row_to_json(p.*),
        'entity_count',  (SELECT COUNT(*) FROM entities WHERE project_id = p.id),
        'task_count',    (SELECT COUNT(*) FROM tasks WHERE project_id = p.id),
        'open_tasks',    (SELECT COUNT(*) FROM tasks WHERE project_id = p.id AND status NOT IN ('completed','cancelled')),
        'adr_count',     (SELECT COUNT(*) FROM adrs WHERE project_id = p.id),
        'token_count',   (SELECT COUNT(*) FROM design_tokens WHERE project_id = p.id),
        'finding_count', (SELECT COUNT(*) FROM code_findings WHERE project_id = p.id),
        'open_findings', (SELECT COUNT(*) FROM code_findings WHERE project_id = p.id AND status = 'open'),
        'conversations', (SELECT COUNT(*) FROM conversations WHERE project_id = p.id)
    )
    FROM projects p WHERE p.id = p_project_id;
$$;
