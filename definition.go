package stagparser

// Definition is a struct tag value element.
type Definition interface {
	// Name is a name of the definition
	Name() string
	// Attributes are attributes of the definition
	Attributes() map[string]interface{}
	// Attribute returns an attribute value and true if an attribute exists
	Attribute(name string) (interface{}, bool)
}

type definition struct {
	name       string
	attributes map[string]interface{}
}

func newDefinition(name string, attributes map[string]interface{}) Definition {
	return &definition{
		name:       name,
		attributes: attributes,
	}
}

func (d *definition) Name() string {
	return d.name
}

func (d *definition) Attributes() map[string]interface{} {
	return d.attributes
}

func (d *definition) Attribute(name string) (interface{}, bool) {
	v, ok := d.attributes[name]
	return v, ok
}
