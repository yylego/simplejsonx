# Function Naming

The extraction API uses verbs with seven letters: `Extract`, `Solicit`, `Inspect`, `Resolve`, `Inquire`, `Attempt`, `Explore`, and `Harvest`. Each verb describes a distinct extraction contract. Precise meaning takes precedence; fixed length is a preference.

`Harvest` means collecting a usable value from a field with mixed representations. Its options define which absent / null inputs are accepted and how strings are converted. It still validates the input according to those options.

Option names follow the intended meaning: `AbsentOk` accepts absent fields, and `Nullable` refers to JSON null. These names need not match the verb length. See the README and `HarvestOption` comments to check exact semantics.
