package database

import (
	"fmt"
	"log/slog"
	"regexp"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var identRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
var qualifiedIdentRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?$`)

func quoteIdent(ident string) string {
	return `"` + ident + `"`
}

func validateIdent(ident string) error {
	if !identRe.MatchString(ident) {
		return fmt.Errorf("invalid identifier: %q", ident)
	}
	return nil
}

func validateQualifiedIdent(ident string) error {
	if !qualifiedIdentRe.MatchString(ident) {
		return fmt.Errorf("invalid qualified identifier: %q", ident)
	}
	return nil
}

// EnsureGeneratedTSVectorColumn ensures a stored generated tsvector column exists.
//
// expr must be a safe SQL expression controlled by code (not user input), e.g.:
// setweight(to_tsvector('simple', coalesce(name,”)), 'A') || ...
func EnsureGeneratedTSVectorColumn(db *gorm.DB, table, column, expr string) {
	if err := validateIdent(table); err != nil {
		slog.Error("fts: invalid table identifier", "table", table, "error", err)
		return
	}
	if err := validateIdent(column); err != nil {
		slog.Error("fts: invalid column identifier", "column", column, "error", err)
		return
	}
	if expr == "" {
		slog.Error("fts: empty tsvector expression", "table", table, "column", column)
		return
	}

	// Use a DO block for broad Postgres compatibility (generated columns + IF NOT EXISTS).
	sql := fmt.Sprintf(`DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = '%s' AND column_name = '%s'
  ) THEN
    ALTER TABLE %s ADD COLUMN %s tsvector GENERATED ALWAYS AS (%s) STORED;
  END IF;
END $$;`, table, column, quoteIdent(table), quoteIdent(column), expr)

	if err := db.Exec(sql).Error; err != nil {
		slog.Error("fts: failed to ensure generated tsvector column", "table", table, "column", column, "error", err)
	}
}

// EnsureGinIndex ensures a GIN index exists on the given table column.
func EnsureGinIndex(db *gorm.DB, indexName, table, column string) {
	if err := validateIdent(indexName); err != nil {
		slog.Error("fts: invalid index identifier", "index", indexName, "error", err)
		return
	}
	if err := validateIdent(table); err != nil {
		slog.Error("fts: invalid table identifier", "table", table, "error", err)
		return
	}
	if err := validateIdent(column); err != nil {
		slog.Error("fts: invalid column identifier", "column", column, "error", err)
		return
	}

	sql := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON %s USING gin (%s);`, quoteIdent(indexName), quoteIdent(table), quoteIdent(column))
	if err := db.Exec(sql).Error; err != nil {
		slog.Error("fts: failed to ensure gin index", "index", indexName, "table", table, "column", column, "error", err)
	}
}

// EnsureTrigramGinIndex ensures a GIN trigram index exists (pg_trgm).
func EnsureTrigramGinIndex(db *gorm.DB, indexName, table, column string) {
	if err := validateIdent(indexName); err != nil {
		slog.Error("trgm: invalid index identifier", "index", indexName, "error", err)
		return
	}
	if err := validateIdent(table); err != nil {
		slog.Error("trgm: invalid table identifier", "table", table, "error", err)
		return
	}
	if err := validateIdent(column); err != nil {
		slog.Error("trgm: invalid column identifier", "column", column, "error", err)
		return
	}

	sql := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON %s USING gin (%s gin_trgm_ops);`, quoteIdent(indexName), quoteIdent(table), quoteIdent(column))
	if err := db.Exec(sql).Error; err != nil {
		slog.Error("trgm: failed to ensure trigram index", "index", indexName, "table", table, "column", column, "error", err)
	}
}

// WebSearchScope returns a WHERE scope matching any of the provided tsvector columns.
// Uses websearch_to_tsquery for user-friendly search syntax.
func WebSearchScope(term string, vectorColumns ...string) (func(*gorm.DB) *gorm.DB, error) {
	if term == "" || len(vectorColumns) == 0 {
		return func(db *gorm.DB) *gorm.DB { return db }, nil
	}
	for _, col := range vectorColumns {
		if err := validateQualifiedIdent(col); err != nil {
			return nil, err
		}
	}

	conds := make([]string, 0, len(vectorColumns))
	vars := make([]any, 0, len(vectorColumns))
	for _, col := range vectorColumns {
		conds = append(conds, fmt.Sprintf("%s @@ websearch_to_tsquery('simple', ?)", col))
		vars = append(vars, term)
	}

	sql := "(" + joinOr(conds) + ")"
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(sql, vars...)
	}, nil
}

// WebSearchRankOrder returns an ORDER BY expression ordering by combined rank.
func WebSearchRankOrder(term string, vectorColumns ...string) (clause.Expr, error) {
	for _, col := range vectorColumns {
		if err := validateQualifiedIdent(col); err != nil {
			return clause.Expr{}, err
		}
	}
	exprs := make([]string, 0, len(vectorColumns))
	vars := make([]any, 0, len(vectorColumns))
	for _, col := range vectorColumns {
		exprs = append(exprs, fmt.Sprintf("ts_rank_cd(%s, websearch_to_tsquery('simple', ?))", col))
		vars = append(vars, term)
	}

	sql := joinPlus(exprs) + " DESC"
	return clause.Expr{SQL: sql, Vars: vars}, nil
}

func joinOr(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += " OR " + parts[i]
	}
	return out
}

func joinPlus(parts []string) string {
	if len(parts) == 0 {
		return "0"
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += " + " + parts[i]
	}
	return out
}
