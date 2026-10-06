package admin

import (
	template2 "html/template"
	"strconv"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// NewUUIDTable initializes a standard GoAdmin table configured for PostgreSQL with a UUID primary key ("id").
// This ensures that GoAdmin parses the primary key as a UUID string rather than defaulting to an integer.
func NewUUIDTable(ctx *context.Context) table.Table {
	return table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql").SetPrimaryKey("id", db.UUID))
}

// AddIDColumn adds the standard sortable UUID primary key column to an InfoPanel.
// Note: FieldCopyable is intentionally avoided because GoAdmin's table action buttons
// interpolate {{(index $info $PrimaryKey).Content}} into URLs, which would inject raw HTML markup.
func AddIDColumn(info *types.InfoPanel) {
	info.AddField("ID", "id", db.UUID).FieldSortable()
}

// AddTimestampColumns adds standard sortable created_at and updated_at columns to an InfoPanel.
func AddTimestampColumns(info *types.InfoPanel) {
	info.AddField("Created At", "created_at", db.Timestamp).FieldSortable()
	info.AddField("Updated At", "updated_at", db.Timestamp).FieldSortable()
}

// AddIDFormField adds a read-only ID field to a FormPanel that is disabled during creation and updates.
func AddIDFormField(formList *types.FormPanel) {
	formList.AddField("ID", "id", db.UUID, form.Text).
		FieldDisableWhenCreate().
		FieldDisableWhenUpdate()
}

// AddTimestampFormFields registers created_at and updated_at on the FormPanel as hidden/disabled fields.
func AddTimestampFormFields(formList *types.FormPanel) {
	formList.AddField("Created At", "created_at", db.Timestamp, form.Datetime).
		FieldDisableWhenCreate().
		FieldDisableWhenUpdate()
	formList.AddField("Updated At", "updated_at", db.Timestamp, form.Datetime).
		FieldDisableWhenCreate().
		FieldDisableWhenUpdate()
}

// WithAutoTimestamps attaches a PreProcessFn that automatically stamps created_at and updated_at.
// Optional custom PreProcess functions can be passed and will be called in sequence.
func WithAutoTimestamps(formList *types.FormPanel, customFns ...types.FormPreProcessFn) {
	formList.SetPreProcessFn(func(values form2.Values) form2.Values {
		now := time.Now().Format("2006-01-02 15:04:05.000000-07:00")
		if values.IsInsertPost() {
			values.Add("created_at", now)
			values.Add("updated_at", now)
		} else if values.IsUpdatePost() || values.IsSingleUpdatePost() {
			values.Add("updated_at", now)
		}
		for _, fn := range customFns {
			values = fn(values)
		}
		return values
	})
}

// NullableIntDisplay returns a FieldFilterFn for InfoPanel that renders nullDisplay when the integer is NULL/0/"".
func NullableIntDisplay(column string, nullDisplay string) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		if value.Row[column] == nil || value.Value == "0" || value.Value == "" {
			return nullDisplay
		}
		return value.Value
	}
}

// NullableIntFormDisplay returns a FieldFilterFn for FormPanel that clears "0" or nil back to "" so placeholders can display.
func NullableIntFormDisplay(column string) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		if value.Row[column] == nil || value.Value == "0" || value.Value == "" {
			return ""
		}
		return value.Value
	}
}

// NullableIntPostFilter converts empty strings from form submissions to nil (SQL NULL), and parses non-empty values as integers.
func NullableIntPostFilter() types.PostFieldFilterFn {
	return func(value types.PostFieldModel) interface{} {
		val := strings.TrimSpace(value.Value.Value())
		if val == "" {
			return nil
		}
		num, err := strconv.Atoi(val)
		if err != nil {
			return val
		}
		return num
	}
}

// AddNullableIntColumn adds a sortable integer column to an InfoPanel that renders fallback text when NULL.
func AddNullableIntColumn(info *types.InfoPanel, head, field, nullDisplay string) *types.InfoPanel {
	return info.AddField(head, field, db.Int).
		FieldSortable().
		FieldDisplay(NullableIntDisplay(field, nullDisplay))
}

// AddNullableIntFormField registers a text-based integer input on a FormPanel that supports SQL NULL when left empty.
func AddNullableIntFormField(formList *types.FormPanel, head, field, placeholder, helpMsg string) *types.FormPanel {
	return formList.AddField(head, field, db.Int, form.Text).
		FieldPlaceholder(placeholder).
		FieldHelpMsg(template2.HTML(helpMsg)).
		FieldDisplay(NullableIntFormDisplay(field)).
		FieldPostFilterFn(NullableIntPostFilter())
}

// NullableUUIDPostFilter converts empty strings from form submissions to nil (SQL NULL).
func NullableUUIDPostFilter() types.PostFieldFilterFn {
	return func(value types.PostFieldModel) interface{} {
		val := strings.TrimSpace(value.Value.Value())
		if val == "" {
			return nil
		}
		return val
	}
}

