package admin

import (
	"fmt"
	"net/http"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	adminCfg "github.com/GoAdminGroup/go-admin/modules/config"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	"github.com/GoAdminGroup/go-admin/modules/language"
	adm "github.com/GoAdminGroup/go-admin/plugins/admin"
	admTable "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	_ "github.com/GoAdminGroup/themes/adminlte"

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/gin-gonic/gin"
	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/database"
)

// Generators maps table names to their GoAdmin table generators
var Generators = admTable.GeneratorList{
	"tickets": GetTicketTable,
}

// Mount initializes and attaches the GoAdmin engine to the Gin router
func Mount(r *gin.Engine, appCfg *config.Config) error {
	eng := engine.Default()

	// 1. Ensure GoAdmin schema and initial admin account exist in PostgreSQL using shared database package
	db, err := database.NewDatabase(appCfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database for GoAdmin: %w", err)
	}
	if err := InitSchema(db); err != nil {
		return fmt.Errorf("failed to initialize GoAdmin database schema: %w", err)
	}

	dsn := database.BuildDSN(appCfg)

	cfg := adminCfg.Config{
		Databases: adminCfg.DatabaseList{
			"default": {
				Dsn:    dsn,
				Driver: adminCfg.DriverPostgresql,
			},
		},
		UrlPrefix: "admin",
		IndexUrl:  "/info/tickets",
		Debug:     appCfg.AppEnv != "production",
		Language:  language.EN,
		Theme:     "adminlte",
		Title:     "Ticket Box Admin",
		Logo:      "<b>Ticket</b>Box",
		MiniLogo:  "<b>T</b>B",
	}

	adminPlugin := adm.NewAdmin(Generators)

	if err := eng.AddConfig(&cfg).
		AddPlugins(adminPlugin).
		Use(r); err != nil {
		return fmt.Errorf("failed to initialize GoAdmin engine: %w", err)
	}

	// Redirect root /admin to default dashboard (/admin/info/tickets)
	r.GET("/admin", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/info/tickets")
	})

	return nil
}
