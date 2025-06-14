# API Json Type System

- `String` => String
- `Number` => Number
- `NULL` => Null
- `?` => Optional (value or undefined)
- `|` => Or
- `<>` => Specific type (just for strings)
  - `URL` = Url
  - `Hex` = Hexadecimal
  - `RGB` = RGB
  - `Number` = Number
  - `Bin` = Binary
  - `Path` = File/folder path
  - `UUID` = Uuid
- `[]` => Array

## Example

```json
{
  "id": "String<UUID>",
  "name": "String",
  "age": "Number",
  "description": "String?",
  "friends": "String?[]",
  "website": "String<URL>?",
  "avatar": "String<URL> | NULL"
}
```
