
package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Connect(connStr string) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pgx.Connect(ctx, connStr)
}

func RunMigrations(conn *pgx.Conn) error {
	files, err := migrationsFS.ReadDir("migrations")
	if err != nil { return err }
	for _, f := range files {
		if f.IsDir() { continue }
		content, err := migrationsFS.ReadFile("migrations/" + f.Name())
		if err != nil { return err }
		stmt := string(content)
		// Split on ; but ignore inside simple cases
		parts := strings.Split(stmt, ";")
		for _, p := range parts {
			q := strings.TrimSpace(p)
			if q == "" { continue }
			if _, err := conn.Exec(context.Background(), q); err != nil {
				log.Printf("migração %s falhou: %v", f.Name(), err)
				return err
			}
		}
		log.Printf("migração %s aplicada", f.Name())
	}
	return nil
}

func ConnString(host string, port string, user string, pass string, dbname string, sslmode string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbname, sslmode)
}
