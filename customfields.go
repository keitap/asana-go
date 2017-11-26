package asana

import "fmt"

type EnumValue struct {
	WithID
	WithName
	Enabled bool   `json:"enabled"`
	Color   string `json:"color"`
}

type FieldType string

// FieldTypes for CustomField.Type field
const (
	Text   FieldType = "text"
	Enum   FieldType = "enum"
	Number FieldType = "number"
)

// Custom Fields store the metadata that is used in order to add user-
// specified information to tasks in Asana. Be sure to reference the Custom
// Fields developer documentation for more information about how custom fields
// relate to various resources in Asana.
type CustomField struct {
	Expandable

	WithName
	WithCreated

	// The type of the custom field. Must be one of the given values:
	// 'text', 'enum', 'number'
	Type FieldType `json:"type"`

	// Only relevant for custom fields of type ‘Enum’. This array specifies
	// the possible values which an enum custom field can adopt.
	EnumOptions []*EnumValue `json:"enum_options,omitempty"`

	// Only relevant for custom fields of type ‘Number’. This field dictates
	// the number of places after the decimal to round to, i.e. 0 is integer
	// values, 1 rounds to the nearest tenth, and so on.
	Precision int `json:"precision,omitempty"`
}

type CustomFieldSetting struct {
	WithID

	CustomField *CustomField `json:"custom_field"`

	Project *Project `json:"project,omitempty"`
}

// When a custom field is associated with a project, tasks in that project can
// carry additional custom field values which represent the value of the field
// on that particular task - for instance, the selected item from an enum type
// custom field. These custom fields will appear as an array in a
// custom_fields property of the task, along with some basic information which
// can be used to associate the custom field value with the custom field
// metadata.
type CustomFieldValue struct {
	CustomField

	// Custom fields of type text will return a text_value property containing
	// the string of text for the field.
	TextValue string `json:"text_value,omitempty"`

	// Custom fields of type number will return a number_value property
	// containing the number for the field.
	NumberValue string `json:"number_value,omitempty"`

	// Custom fields of type enum will return an enum_value property
	// containing an object that represents the selection of the enum value.
	EnumValue *EnumValue `json:"enum_value,omitempty"`
}

// CustomField retrieves a custom field record by ID
func (c *Client) CustomField(id int64) *CustomField {
	result := &CustomField{}
	result.init(id, c)
	return result
}

// Expand loads the full details for this CustomField
func (f *CustomField) Expand() error {
	f.trace("Loading details for custom field %q", f.ID)

	if f.expanded {
		return nil
	}

	_, err := f.client.get(fmt.Sprintf("/custom_fields/%d", f.ID), nil, f)
	return err
}

// CustomFields returns the compact records for all custom fields in the workspace
func (w *Workspace) CustomFields(options ...*Options) ([]*CustomField, *NextPage, error) {
	w.trace("Listing custom fields in workspace %d...\n", w.ID)
	var result []*CustomField

	// Make the request
	nextPage, err := w.client.get(fmt.Sprintf("/workspaces/%d/custom_fields", w.ID), nil, &result, options...)
	return result, nextPage, err
}

// AllCustomFields repeatedly pages through all available custom fields in a workspace
func (w *Workspace) AllCustomFields(options ...*Options) ([]*CustomField, error) {
	allCustomFields := []*CustomField{}
	nextPage := &NextPage{}

	var customFields []*CustomField
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  100,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		customFields, nextPage, err = w.CustomFields(allOptions...)
		if err != nil {
			return nil, err
		}

		allCustomFields = append(allCustomFields, customFields...)
	}
	return allCustomFields, nil
}
