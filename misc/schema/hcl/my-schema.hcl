title = "My schema"
description = "Test"
type = "object"
required = [ "foo" ]
properties = {
  foo = {
    description = "Ça fait des trucs"
    type = "object"
    properties = {
      bar = {
        type = "string"
        minLength = 5
      }
      enabled = {
        type = "boolean"
      }
      myItems = {
        type = "array"
        items = {
          type = "object"
          properties = {
            stringProp = {
              type = "string"
              minLength = 10
            }
            numericProp = {
              type = "integer"
            }
          }
        }
      }
    }
  }
}