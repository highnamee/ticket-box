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
	ticketTable = NewUUIDTable(ctx)

	info := ticketTable.GetInfo().SetFilterFormLayout(form.LayoutTwoCol)
	AddIDColumn(info)
	info.AddField("Name", "name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Price ($)", "price", db.Decimal).
		FieldSortable()
	info.AddField("Total Qty", "total_quantity", db.Int).
		FieldSortable()
	info.AddField("Available Stock", "available_stock", db.Int).
		FieldSortable()
	AddNullableIntColumn(info, "Max/User", "max_booking_per_user", "Unlimited")
	info.AddField("Status", "status", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "ACTIVE", Text: "ACTIVE"},
			{Value: "INACTIVE", Text: "INACTIVE"},
		})
	info.AddField("Category ID", "category_id", db.UUID).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Venue", "venue", db.Varchar)
	info.AddField("Date", "date", db.Varchar)
	info.AddField("Featured", "featured", db.Boolean).
		FieldBool()
	AddTimestampColumns(info)

	info.SetTable("tickets").
		SetTitle("Tickets").
		SetDescription("Manage Event Tickets and Inventory")

	formList := ticketTable.GetForm()
	AddIDFormField(formList)
	formList.AddField("Name", "name", db.Varchar, form.Text).
		FieldMust()
	formList.AddField("Category", "category_id", db.UUID, form.Select).
		FieldOptionsFromTable("categories", "name", "id").
		FieldPostFilterFn(NullableUUIDPostFilter())
	formList.AddField("Venue", "venue", db.Varchar, form.Text)
	formList.AddField("Date", "date", db.Varchar, form.Text)
	formList.AddField("Time", "time", db.Varchar, form.Text)
	formList.AddField("Image URL", "image_url", db.Text, form.Text)
	formList.AddField("Featured", "featured", db.Boolean, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "true", Text: "Yes"},
			{Value: "false", Text: "No"},
		}).
		FieldDefault("false")
	formList.AddField("Description", "description", db.Text, form.TextArea)
	formList.AddField("Price ($)", "price", db.Decimal, form.Currency).
		FieldMust()
	formList.AddField("Total Quantity", "total_quantity", db.Int, form.Number).
		FieldMust()
	formList.AddField("Available Stock", "available_stock", db.Int, form.Number).
		FieldMust()
	AddNullableIntFormField(formList, "Max Per User", "max_booking_per_user", "Unlimited (leave blank)", "Optional maximum tickets a single user can book. Leave blank for unlimited (NULL).")
	formList.AddField("Status", "status", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "ACTIVE", Text: "ACTIVE"},
			{Value: "INACTIVE", Text: "INACTIVE"},
		}).
		FieldDefault("INACTIVE")
	AddTimestampFormFields(formList)
	WithAutoTimestamps(formList)

	formList.SetTable("tickets").
		SetTitle("Ticket Details").
		SetDescription("Create or update ticket information")

	return
}

// GetUserTable returns the table definition for User model management
func GetUserTable(ctx *context.Context) (userTable table.Table) {
	userTable = NewUUIDTable(ctx)

	info := userTable.GetInfo().SetFilterFormLayout(form.LayoutTwoCol)
	AddIDColumn(info)
	info.AddField("Email", "email", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Full Name", "full_name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Role", "role", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "USER", Text: "USER"},
			{Value: "ADMIN", Text: "ADMIN"},
		})
	info.AddField("Status", "status", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "ACTIVE", Text: "ACTIVE"},
			{Value: "INACTIVE", Text: "INACTIVE"},
			{Value: "BANNED", Text: "BANNED"},
		})
	AddTimestampColumns(info)

	info.SetTable("users").
		SetTitle("Users").
		SetDescription("Manage Platform Users and Roles")

	formList := userTable.GetForm()
	AddIDFormField(formList)
	formList.AddField("Email", "email", db.Varchar, form.Email).
		FieldMust()
	formList.AddField("Full Name", "full_name", db.Varchar, form.Text).
		FieldMust()
	formList.AddField("Role", "role", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "USER", Text: "USER"},
			{Value: "ADMIN", Text: "ADMIN"},
		}).
		FieldDefault("USER")
	formList.AddField("Status", "status", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "ACTIVE", Text: "ACTIVE"},
			{Value: "INACTIVE", Text: "INACTIVE"},
			{Value: "BANNED", Text: "BANNED"},
		}).
		FieldDefault("ACTIVE")
	AddTimestampFormFields(formList)
	WithAutoTimestamps(formList)

	formList.SetTable("users").
		SetTitle("User Details").
		SetDescription("View and update user status and role")

	return
}

// GetCategoryTable returns the table definition for Category model management
func GetCategoryTable(ctx *context.Context) (categoryTable table.Table) {
	categoryTable = NewUUIDTable(ctx)

	info := categoryTable.GetInfo().SetFilterFormLayout(form.LayoutTwoCol)
	AddIDColumn(info)
	info.AddField("Name", "name", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Slug", "slug", db.Varchar).
		FieldSortable().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Description", "description", db.Text)
	AddTimestampColumns(info)

	info.SetTable("categories").
		SetTitle("Categories").
		SetDescription("Manage Event Categories")

	formList := categoryTable.GetForm()
	AddIDFormField(formList)
	formList.AddField("Name", "name", db.Varchar, form.Text).
		FieldMust()
	formList.AddField("Slug", "slug", db.Varchar, form.Text).
		FieldMust()
	formList.AddField("Description", "description", db.Text, form.TextArea)
	AddTimestampFormFields(formList)
	WithAutoTimestamps(formList)

	formList.SetTable("categories").
		SetTitle("Category Details").
		SetDescription("Create or update event category")

	return
}

