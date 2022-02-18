# Policy based Dataset synchronization
The Sync Service allows for a user defined policy to be applied to object store bucket notifications allowing the user to apply a custom filter to what files are synchronized between datasets.

## Policy Specification
The following section defines the Sync Policy JSON object that is used to filter which files are synced between datasets.

> Note: all fields are required unless noted otherwise

### PolicyObject
```json
{
  "Version": "1",
  "Effect": EffectField,
  "Statements": [StatementObject, ...]
}
```

A PolicyObject consists of the following fields:
* Version: Currently only Version 1 is defined and supported
* Effect: An EffectField setting how the statements should be combined together
  - Optional, defaults to `"OR"`
* Statements: A list of StatementObjects whose logic will be combined together for a final policy decision

### StatementObject
```json
{
  "Id": String,
  "Effect": EffectField,
  "Conditions": [ConditionObject, ...]
}
```

A StatementObject consists of the following fields:
* Id: An identifier for the statement
  - No meaning is given to this field by the Hoss
* Effect: An EffectField setting how the conditions should be combined together
  - Optional, defaults to `"AND"`
* Conditions: A list of ConditionObjects whose logic will be combined together for a final statement decision

### ConditionObject
```json
{
  "Left": LeftOperandField,
  "Right": RightOperandField,
  "Operator": OperatorField,
}
```

A ConditionObject consists of the following fields:
* Left: The left side of the condition
* Right: The right side of the condition
* Operator: The operation to apply to the two operands

### EffectField
```json
"Effect": String
```

An EffectField can have one of the following values
* `"AND"`: Logically and all of the sub-object logic together
* `"OR"`: Logically inclusivly or all of the sub-object logic together

### LeftOperandField
```json
"Left": String
```

An LeftOperandField can have one of the following values
* `"event:operation"`: The type of notification event generated for the file
  - The value is either `"PUT"` (for create or update) or `"DELETE"` (for delete)
* `"object:key"`: The key of the object that is the focus of the notification message
* `"object:size"`: The size of the object (in bytes) that is the focus of the notification message
* `"object:metadata"`: The object's metadata dictionary
  - Used with the `"has"` operator to check if a metadata key exists
* `"object:metadata:<key>"`: The key of the object's metadata to use in the conditional
  - `"<key>"`: Is a string containing the name of the metadata key
  - All metadata values are strings

### RightOperandField
```json
"Right": String|Number
```

An RightOperandField can have one of the following values
* `""`: An empty string can be used to verify that a metadata value doesn't exist
* `"<glob>"`: A glob expression to match against the `Left` operand
  - Only supports `"=="` and `"!="` operators
* `<number>: An integer or float number

### OperatorField
```json
"Operator": String
```

An OperatorField can have one of the following values
* `"=="`: Returns true if the two operands are equal
* `"!="`: Returns true if the two operands are not equal
* `">"`: Returns true if the left operand is greater than the right operand
* `"<"`: Returns true if the left operand is less than the right operand
* `">="`: Returns true if the left operand is greater than or equal to the right operand
* `"<="`: Returns true if the left operand is less than or equal to the right operand
* `"has"`: Returns true if the left operand dictionary contains the right operand key

## Example Policies

### Existing Behavior
This policy will reproduce the current (pre-policy) sync service behavior replicating all file create, update & delete operations, along with updating the search index for the server where the sync service is running.

```json
{
  "Version": "1",
  "Statements":[]
}
```

### Example Conditional Statements
This policy prevents the sync of raw files (any file ending with `.raw`) and only allows syncing of files with the specific metadata value.

> Note: If a raw file contains the specific metadata value it will be synced, because the different StatementObjects policy decisions are `or`ed together for the final policy decision

```json
{
  "Version": "1",
  "Effect": "OR",
  "Statements":[
    {
      "Id": "IgnoreRawData",
      "Conditions":[
        {
          "Left": "object:key",
          "Right": "*.raw",
          "Operator": "!="
        }
      ]
    },  
    {
      "Id": "RequireMetadataValue",
      "Conditions":[
        {
          "Left": "object:metadata:my-key-1",
          "Right": "expected-value-1",
          "Operator": "=="
        }
      ]
    }                    
  ]
}
```
