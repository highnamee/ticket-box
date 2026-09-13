package admin

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// GetTicketTable returns the table definition for Ticket model management
func GetTicketTable(ctx *context.Context) (ticketTable table.Table) {
	ticketTable = table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := ticketTable.GetInfo().SetFilterFormLayout(form.LayoutTwoCol)
	info.AddField("ID", "id", db.UUID).
		FieldSortable().
		FieldCopyable()
	info.AddField("Name", "name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Price ($)", "price", db.Decimal).
		FieldSortable()
	info.AddField("Total Qty", "total_quantity", db.Int).
		FieldSortable()
	info.AddField("Available Stock", "available_stock", db.Int).
		FieldSortable()
	info.AddField("Status", "status", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "ACTIVE", Text: "ACTIVE"},
			{Value: "SOLD_OUT", Text: "SOLD_OUT"},
			{Value: "INACTIVE", Text: "INACTIVE"},
		})
	info.AddField("Created At", "created_at", db.Timestamp).
		FieldSortable()
	info.AddField("Updated At", "updated_at", db.Timestamp).
		FieldSortable()

	info.SetTable("tickets").
		SetTitle("Tickets").
		SetDescription("Manage Event Tickets and Inventory")

	formList := ticketTable.GetForm()
	formList.AddField("ID", "id", db.UUID, form.Text).
		FieldDisableWhenCreate().
		FieldDisableWhenUpdate()
	formList.AddField("Name", "name", db.Varchar, form.Text).
		FieldMust()
	formList.AddField("Description", "description", db.Text, form.TextArea)
	formList.AddField("Price ($)", "price", db.Decimal, form.Currency).
		FieldMust()
	formList.AddField("Total Quantity", "total_quantity", db.Int, form.Number).
		FieldMust()
	formList.AddField("Available Stock", "available_stock", db.Int, form.Number).
		FieldMust()
	formList.AddField("Status", "status", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "ACTIVE", Text: "ACTIVE"},
			{Value: "SOLD_OUT", Text: "SOLD_OUT"},
			{Value: "INACTIVE", Text: "INACTIVE"},
		}).
		FieldDefault("INACTIVE")

	formList.SetTable("tickets").
		SetTitle("Ticket Details").
		SetDescription("Create or update ticket information")

	return
}
