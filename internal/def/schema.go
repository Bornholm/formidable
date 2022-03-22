package def

import "github.com/santhosh-tekuri/jsonschema/v5"

const rawSchema = `
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
    "title": "Formidable default schema",
    "type": ["null", "boolean", "object", "array", "number", "string"]
}`

var Schema = jsonschema.MustCompileString("", rawSchema)
