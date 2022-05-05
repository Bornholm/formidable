test = 1

test1 = 2 + 1

foo = {
  description = "Ça fait des trucs"
  type = "object"
  properties = {
    type = "string"
    minLength = 5
  }
  test = [
    "foo", 
    {
      test = "foo"
    },
    5 + 5.2
  ]
}