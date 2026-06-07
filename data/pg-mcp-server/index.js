#!/usr/bin/env node
/**
 * pg-mcp-server — MCP serveur PostgreSQL pour OpenCode.
 *
 * Outils génériques : pg_health, pg_query, pg_list_tables
 * Outils métier    : memory_*, task_*, project_*, adr_*, design_token_*, conversation_*
 *
 * Version avec sous-processus psql (fiable, pas de dépendance module pg).
 */

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import { execSync } from "child_process";

// ─── Configuration ───────────────────────────────────────────────────────────

const DATABASE_URL = process.env.DATABASE_URL;
const url = DATABASE_URL ? new URL(DATABASE_URL) : null;
const PGUSER = url?.username || "postgres";
const PGPASSWORD = url?.password || "postgres";
const PGHOST = url?.hostname || "localhost";
const PGPORT = url?.port || "5432";
const PGDATABASE = url?.pathname?.replace(/^\//, "") || "opencode_db";

const ENV = { PGPASSWORD, PGUSER, PGHOST, PGPORT, PGDATABASE, PATH: process.env.PATH || "" };
const PSQL = `psql -h ${PGHOST} -p ${PGPORT} -U ${PGUSER} -d ${PGDATABASE} -t -A --no-psqlrc`;

// ─── Helper psql ─────────────────────────────────────────────────────────────

function query(sql, params = []) {
  let i = 0;
  const escaped = sql.replace(/\$(\d+)/g, (_, n) => {
    const val = params[parseInt(n) - 1];
    if (val === null || val === undefined) return "NULL";
    if (typeof val === "number") return String(val);
    if (typeof val === "boolean") return val ? "TRUE" : "FALSE";
    return `'${String(val).replace(/'/g, "''")}'`;
  });

  const out = execSync(`${PSQL} -c "${escaped.replace(/"/g, '\\"')}"`, {
    env: ENV, encoding: "utf-8", timeout: 15000, maxBuffer: 10 * 1024 * 1024,
  });

  const lines = out.trim().split("\n").filter(Boolean);
  if (!lines.length) return { rows: [], rowCount: 0 };

  // Try JSON output for structured queries
  try {
    const jsonOut = execSync(
      `${PSQL} -c "SELECT json_agg(row_to_json(t)) FROM (${escaped.replace(/;$/, '')}) t"`,
      { env: ENV, encoding: "utf-8", timeout: 15000, maxBuffer: 10 * 1024 * 1024 }
    );
    const raw = jsonOut.trim();
    if (raw && raw !== "NULL") {
      const rows = JSON.parse(raw);
      return { rows, rowCount: rows.length };
    }
  } catch {}

  return { rows: lines.map(l => ({ value: l })), rowCount: lines.length };
}

function queryOne(sql, params = []) {
  const r = query(sql, params);
  return r.rows?.[0] || null;
}

function jsonReply(data) {
  return { content: [{ type: "text", text: JSON.stringify(data, null, 2) }] };
}

function textReply(text) {
  return { content: [{ type: "text", text }] };
}

// ─── Server ──────────────────────────────────────────────────────────────────

const server = new McpServer({
  name: "pg-mcp-server",
  version: "2.0.0",
});

// ═══════════════════════════════════════════════════════════════════════════════
//  OUTILS GÉNÉRIQUES
// ═══════════════════════════════════════════════════════════════════════════════

server.registerTool("pg_health", {
  title: "Postgres Health",
  description: "Statut de la connexion PostgreSQL et métriques du projet.",
  inputSchema: z.object({}).strict(),
  annotations: { readOnlyHint: true, destructiveHint: false, idempotentHint: true, openWorldHint: false },
}, async () => {
  try {
    const version = queryOne("SELECT version()");
    const info = queryOne(`
      SELECT current_database() AS db,
             pg_size_pretty(pg_database_size(current_database())) AS size,
             (SELECT COUNT(*)::int FROM projects) AS projects,
             (SELECT COUNT(*)::int FROM entities) AS entities,
             (SELECT COUNT(*)::int FROM tasks) AS tasks,
             (SELECT COUNT(*)::int FROM adrs) AS adrs,
             (SELECT COUNT(*)::int FROM code_findings WHERE status='open') AS open_findings
    `);
    return jsonReply({ connected: true, version: version?.version || "?", ...info });
  } catch (err) {
    return textReply(`Erreur connexion: ${err.message}`);
  }
});

server.registerTool("pg_query", {
  title: "Postgres Query",
  description: "Exécute une requête SQL. Utilise $1, $2 pour les paramètres.",
  inputSchema: z.object({
    sql: z.string().min(1).max(500000).describe("Requête SQL"),
    params: z.array(z.any()).optional().describe("Paramètres positionnels ($1, $2...)"),
  }).strict(),
  annotations: { readOnlyHint: true, destructiveHint: false, idempotentHint: true, openWorldHint: false },
}, async ({ sql, params }) => {
  try {
    return jsonReply(query(sql, params || []));
  } catch (err) {
    return textReply(`Erreur SQL: ${err.message}`);
  }
});

server.registerTool("pg_list_tables", {
  title: "List Tables",
  description: "Liste les tables, colonnes et types d'un schéma.",
  inputSchema: z.object({ schema: z.string().default("public") }).strict(),
  annotations: { readOnlyHint: true, destructiveHint: false, idempotentHint: true, openWorldHint: false },
}, async ({ schema }) => {
  try {
    const tables = query("SELECT table_name, table_type FROM information_schema.tables WHERE table_schema = $1 ORDER BY table_name", [schema]);
    const cols = query(
      "SELECT table_name, column_name, data_type, is_nullable, ordinal_position FROM information_schema.columns WHERE table_schema = $1 ORDER BY table_name, ordinal_position",
      [schema]
    );
    const byTable = {};
    for (const c of cols.rows) {
      if (!byTable[c.table_name]) byTable[c.table_name] = [];
      byTable[c.table_name].push(`${c.column_name} ${c.data_type}${c.is_nullable === 'NO' ? ' PK/NN' : ''}`);
    }
    const lines = [`# Schéma: ${schema}\n`];
    for (const t of tables.rows) {
      lines.push(`## ${t.table_name} (${t.table_type === 'VIEW' ? 'vue' : 'table'})`);
      for (const cdef of (byTable[t.table_name] || [])) lines.push(`  - ${cdef}`);
      lines.push("");
    }
    return textReply(lines.join("\n"));
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

// ═══════════════════════════════════════════════════════════════════════════════
//  OUTILS MÉTIER — Knowledge Graph (mémoire persistante)
// ═══════════════════════════════════════════════════════════════════════════════

server.registerTool("memory_search", {
  title: "Mémoire — Recherche",
  description: `Recherche dans le graphe de connaissance (entités + observations).
  Retourne les entités les plus pertinentes avec leurs relations.

Args:
  - project_slug (string): Slug du projet
  - query (string): Texte à rechercher (recherche plein texte)
  - limit (number): Max résultats (défaut: 10)`,
  inputSchema: z.object({
    project_slug: z.string().min(1).describe("Slug du projet (ex: 'mon-projet')"),
    query: z.string().min(1).describe("Recherche plein texte"),
    limit: z.number().int().min(1).max(50).default(10),
  }).strict(),
  annotations: { readOnlyHint: true, destructiveHint: false, idempotentHint: true, openWorldHint: false },
}, async ({ project_slug, query: q, limit: lim }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    const entities = query("SELECT * FROM search_entities($1, $2, $3)", [pid.id, q, lim]);
    return jsonReply(entities);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

server.registerTool("memory_add", {
  title: "Mémoire — Ajouter entité",
  description: `Ajoute une entité au graphe de connaissance (ou met à jour ses observations).

Args:
  - project_slug (string): Slug du projet
  - name (string): Nom unique de l'entité
  - entity_type (string): Type (component, module, api, pattern, decision, concept, person, tool)
  - observations (string[]): Liste d'observations/facts
  - connect_to (string, optional): Nom d'une entité à relier (relation "related_to")`,
  inputSchema: z.object({
    project_slug: z.string().min(1),
    name: z.string().min(1).max(255).describe("Nom de l'entité"),
    entity_type: z.enum(["component", "module", "api", "pattern", "decision", "concept", "person", "tool", "skill"]),
    observations: z.array(z.string()).default([]).describe("Faits / observations"),
    connect_to: z.string().optional().describe("Nom d'une entité existante à relier"),
  }).strict(),
  annotations: { readOnlyHint: false, destructiveHint: false, idempotentHint: false, openWorldHint: false },
}, async ({ project_slug, name, entity_type, observations, connect_to }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable. Crée-le d'abord.`);

    // Upsert
    const existing = queryOne("SELECT id, observations FROM entities WHERE project_id = $1 AND name = $2", [pid.id, name]);
    if (existing) {
      const merged = [...new Set([...(existing.observations || []), ...observations])];
      query("UPDATE entities SET observations = $1, entity_type = $2, updated_at = NOW() WHERE id = $3",
        [merged, entity_type, existing.id]);
    } else {
      query("INSERT INTO entities (project_id, name, entity_type, observations) VALUES ($1, $2, $3, $4)",
        [pid.id, name, entity_type, observations]);
    }

    // Relation optionnelle
    if (connect_to) {
      const target = queryOne("SELECT id FROM entities WHERE project_id = $1 AND name = $2", [pid.id, connect_to]);
      if (target) {
        const source = queryOne("SELECT id FROM entities WHERE project_id = $1 AND name = $2", [pid.id, name]);
        query("INSERT INTO relations (project_id, source_entity_id, target_entity_id, relation_type) VALUES ($1, $2, $3, 'related_to') ON CONFLICT DO NOTHING",
          [pid.id, source.id, target.id]);
      }
    }

    return textReply(`✓ Entité "${name}" (${entity_type}) enregistrée dans "${project_slug}".`);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

server.registerTool("memory_graph", {
  title: "Mémoire — Graphe complet",
  description: "Retourne le graphe de connaissance complet d'un projet.",
  inputSchema: z.object({
    project_slug: z.string().min(1),
  }).strict(),
  annotations: { readOnlyHint: true, destructiveHint: false, idempotentHint: true, openWorldHint: false },
}, async ({ project_slug }) => {
  try {
    const pid = queryOne("SELECT id, name FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    const entities = query(
      "SELECT name, entity_type, observations, (SELECT COUNT(*) FROM relations WHERE source_entity_id = e.id) AS relation_count FROM entities e WHERE project_id = $1 ORDER BY entity_type, name",
      [pid.id]
    );
    const relations = query(
      `SELECT e1.name AS source, r.relation_type, e2.name AS target
       FROM relations r
       JOIN entities e1 ON e1.id = r.source_entity_id
       JOIN entities e2 ON e2.id = r.target_entity_id
       WHERE r.project_id = $1 ORDER BY e1.name`,
      [pid.id]
    );

    return jsonReply({ project: pid, entities: entities.rows, relations: relations.rows });
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

// ═══════════════════════════════════════════════════════════════════════════════
//  OUTILS MÉTIER — Tâches
// ═══════════════════════════════════════════════════════════════════════════════

server.registerTool("task_list", {
  title: "Tâches — Lister",
  description: "Liste les tâches d'un projet, filtrées par statut.",
  inputSchema: z.object({
    project_slug: z.string().min(1),
    status: z.string().optional().describe("Filtre: pending, in_progress, completed, cancelled, blocked"),
    limit: z.number().int().min(1).max(100).default(20),
  }).strict(),
  annotations: { readOnlyHint: true, destructiveHint: false, idempotentHint: true, openWorldHint: false },
}, async ({ project_slug, status, limit: lim }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    let sql = "SELECT title, status, priority, agent_type, tags, created_at, updated_at FROM tasks WHERE project_id = $1";
    const params = [pid.id];
    if (status) { params.push(status); sql += " AND status = $2"; }
    sql += " ORDER BY priority DESC, created_at DESC LIMIT $" + (params.length + 1);
    params.push(lim);

    const tasks = query(sql, params);
    return jsonReply(tasks);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

server.registerTool("task_create", {
  title: "Tâches — Créer",
  description: "Crée une nouvelle tâche.",
  inputSchema: z.object({
    project_slug: z.string().min(1),
    title: z.string().min(1).max(500),
    description: z.string().optional(),
    priority: z.enum(["critical", "high", "medium", "low"]).default("medium"),
    status: z.enum(["pending", "in_progress", "completed", "cancelled", "blocked"]).default("pending"),
    agent_type: z.string().optional(),
    tags: z.array(z.string()).optional(),
  }).strict(),
  annotations: { readOnlyHint: false, destructiveHint: false, idempotentHint: false, openWorldHint: false },
}, async (params) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [params.project_slug]);
    if (!pid) return textReply(`Projet "${params.project_slug}" introuvable.`);

    query(
      "INSERT INTO tasks (project_id, title, description, priority, status, agent_type, tags) VALUES ($1,$2,$3,$4,$5,$6,$7)",
      [pid.id, params.title, params.description || null, params.priority, params.status, params.agent_type || null, params.tags || null]
    );
    return textReply(`✓ Tâche "${params.title}" créée.`);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

// ═══════════════════════════════════════════════════════════════════════════════
//  OUTILS MÉTIER — Projets
// ═══════════════════════════════════════════════════════════════════════════════

server.registerTool("project_create", {
  title: "Projet — Créer",
  description: "Crée un nouveau projet.",
  inputSchema: z.object({
    name: z.string().min(1).max(255),
    slug: z.string().min(1).max(255).regex(/^[a-z0-9-]+$/, "slug: seulement a-z, 0-9, tirets"),
    description: z.string().optional(),
    root_path: z.string().optional(),
    tech_stack: z.array(z.string()).optional(),
  }).strict(),
  annotations: { readOnlyHint: false, destructiveHint: false, idempotentHint: false, openWorldHint: false },
}, async ({ name, slug, description, root_path, tech_stack }) => {
  try {
    query("INSERT INTO projects (name, slug, description, root_path, tech_stack) VALUES ($1,$2,$3,$4,$5)",
      [name, slug, description || null, root_path || null, JSON.stringify(tech_stack || [])]);
    return textReply(`✓ Projet "${name}" (${slug}) créé.`);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

server.registerTool("project_dashboard", {
  title: "Projet — Dashboard",
  description: "Tableau de bord complet d'un projet (stats, mémoire, tâches, décisions).",
  inputSchema: z.object({
    project_slug: z.string().min(1),
  }).strict(),
  annotations: { readOnlyHint: true, destructiveHint: false, idempotentHint: true, openWorldHint: false },
}, async ({ project_slug }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    const dash = queryOne("SELECT project_dashboard($1)", [pid.id]);
    return jsonReply(dash?.project_dashboard || dash);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

// ═══════════════════════════════════════════════════════════════════════════════
//  OUTILS MÉTIER — ADR (Architecture Decision Records)
// ═══════════════════════════════════════════════════════════════════════════════

server.registerTool("adr_list", {
  title: "ADR — Lister",
  description: "Liste les décisions d'architecture d'un projet.",
  inputSchema: z.object({
    project_slug: z.string().min(1),
    status: z.string().optional().describe("Filtre: proposed, accepted, deprecated, superseded"),
  }).strict(),
  annotations: { readOnlyHint: true, destructiveHint: false, idempotentHint: true, openWorldHint: false },
}, async ({ project_slug, status }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    let sql = "SELECT title, status, LEFT(context, 200) AS context_short, created_at FROM adrs WHERE project_id = $1";
    const params = [pid.id];
    if (status) { params.push(status); sql += " AND status = $2"; }
    sql += " ORDER BY created_at DESC";

    return jsonReply(query(sql, params));
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

server.registerTool("adr_create", {
  title: "ADR — Créer",
  description: "Enregistre une décision d'architecture.",
  inputSchema: z.object({
    project_slug: z.string().min(1),
    title: z.string().min(1).max(500),
    status: z.enum(["proposed", "accepted", "deprecated", "superseded"]).default("proposed"),
    context: z.string().describe("Contexte du problème"),
    decision: z.string().describe("Décision prise"),
    consequences: z.string().describe("Conséquences (positives et négatives)"),
  }).strict(),
  annotations: { readOnlyHint: false, destructiveHint: false, idempotentHint: false, openWorldHint: false },
}, async ({ project_slug, title, status, context, decision, consequences }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    query("INSERT INTO adrs (project_id, title, status, context, decision, consequences) VALUES ($1,$2,$3,$4,$5,$6)",
      [pid.id, title, status, context, decision, consequences]);
    return textReply(`✓ ADR "${title}" [${status}] enregistré.`);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

// ═══════════════════════════════════════════════════════════════════════════════
//  OUTILS MÉTIER — Design Tokens
// ═══════════════════════════════════════════════════════════════════════════════

server.registerTool("design_token_list", {
  title: "Design Tokens — Lister",
  description: "Liste les design tokens d'un projet (couleurs, typographie, espacements, ombres...).",
  inputSchema: z.object({
    project_slug: z.string().min(1),
    category: z.string().optional().describe("Filtre: color, typography, spacing, shadow, radius, opacity, font, breakpoint"),
  }).strict(),
  annotations: { readOnlyHint: true, destructiveHint: false, idempotentHint: true, openWorldHint: false },
}, async ({ project_slug, category }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    let sql = "SELECT category, name, value, description FROM design_tokens WHERE project_id = $1";
    const params = [pid.id];
    if (category) { params.push(category); sql += " AND category = $2"; }
    sql += " ORDER BY category, name";

    return jsonReply(query(sql, params));
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

server.registerTool("design_token_add", {
  title: "Design Tokens — Ajouter",
  description: "Ajoute un design token à un projet.",
  inputSchema: z.object({
    project_slug: z.string().min(1),
    category: z.enum(["color", "typography", "spacing", "shadow", "radius", "opacity", "font", "breakpoint"]),
    name: z.string().min(1).max(255),
    value: z.string().min(1),
    description: z.string().optional(),
  }).strict(),
  annotations: { readOnlyHint: false, destructiveHint: false, idempotentHint: false, openWorldHint: false },
}, async ({ project_slug, category, name, value, description }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    query(
      "INSERT INTO design_tokens (project_id, category, name, value, description) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (project_id, category, name) DO UPDATE SET value = $4, description = COALESCE($5, design_tokens.description)",
      [pid.id, category, name, value, description || null]
    );
    return textReply(`✓ Token ${category}.${name} = ${value}`);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

// ═══════════════════════════════════════════════════════════════════════════════
//  OUTILS MÉTIER — Conversations / Code Findings
// ═══════════════════════════════════════════════════════════════════════════════

server.registerTool("conversation_log", {
  title: "Conversation — Journaliser",
  description: "Enregistre un résumé de conversation pour mémoire.",
  inputSchema: z.object({
    project_slug: z.string().min(1),
    title: z.string().max(500),
    summary: z.string(),
    agent_type: z.string().optional(),
    tags: z.array(z.string()).optional(),
  }).strict(),
  annotations: { readOnlyHint: false, destructiveHint: false, idempotentHint: false, openWorldHint: false },
}, async ({ project_slug, title, summary, agent_type, tags }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    query("INSERT INTO conversations (project_id, title, summary, agent_type, tags) VALUES ($1,$2,$3,$4,$5)",
      [pid.id, title, summary, agent_type || null, tags || null]);
    return textReply(`✓ Conversation "${title}" journalisée.`);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

server.registerTool("finding_add", {
  title: "Code Finding — Ajouter",
  description: "Enregistre un résultat d'analyse de code (bug, sécurité, perf...).",
  inputSchema: z.object({
    project_slug: z.string().min(1),
    finding_type: z.enum(["bug", "smell", "security", "performance", "style", "duplication", "coverage"]),
    severity: z.enum(["critical", "high", "medium", "low", "info"]),
    title: z.string().min(1).max(500),
    description: z.string().optional(),
    file_path: z.string().optional(),
    line_start: z.number().int().optional(),
    suggestion: z.string().optional(),
  }).strict(),
  annotations: { readOnlyHint: false, destructiveHint: false, idempotentHint: false, openWorldHint: false },
}, async ({ project_slug, finding_type, severity, title, description, file_path, line_start, suggestion }) => {
  try {
    const pid = queryOne("SELECT id FROM projects WHERE slug = $1", [project_slug]);
    if (!pid) return textReply(`Projet "${project_slug}" introuvable.`);

    query(
      "INSERT INTO code_findings (project_id, finding_type, severity, title, description, file_path, line_start, suggestion) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
      [pid.id, finding_type, severity, title, description || null, file_path || null, line_start || null, suggestion || null]
    );
    return textReply(`✓ Finding [${severity}] ${title}`);
  } catch (err) {
    return textReply(`Erreur: ${err.message}`);
  }
});

// ═══════════════════════════════════════════════════════════════════════════════
//  STARTUP
// ═══════════════════════════════════════════════════════════════════════════════

async function main() {
  try {
    const v = queryOne("SELECT version()");
    console.error(`pg-mcp-server v2 connecté: ${String(v?.version || v?.value || '?').slice(0, 60)}`);
  } catch (err) {
    console.error(`pg-mcp-server ERREUR: ${err.message}`);
    process.exit(1);
  }
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("pg-mcp-server v2 prêt (stdio) — 16 outils");
}

main().catch(err => { console.error("pg-mcp-server fatal:", err); process.exit(1); });
