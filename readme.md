# Mermerd

Create [Mermaid-Js](https://mermaid-js.github.io/mermaid/#/entityRelationshipDiagram) ERD diagrams from existing tables.

## Features

* Supports PostgreSQL and MySQL
* Select from available schemas
* Select only the tables you are interested in
* Show only the constraints that you are interested in
* interactive cli (multiselect, search for tables and schemas, etc.)

## Why would I need it / Why should I care?

Documenting stuff and keeping it updated is hard and tedious, but having the right documentation can help to make the
right decisions. Mermerd was designed to be able to export an existing database schema in a format that can be used to
prototype and plan new features based on the existing schema. The resulting output is an ERD diagram
in [Mermaid-Js](https://mermaid-js.github.io/mermaid/#/entityRelationshipDiagram) format that can easily be updated and
extended.

## How does it work

1. Specify the connection string (via parameter or interactive cli)
2. Specify the schema that should be used (via parameter or interactive cli)
3. Select the tables that you are interested in (multiselect, at least 1)
4. Enjoy your current database schema in Mermaid-JS format

https://user-images.githubusercontent.com/22556363/149669994-bd5cfd8d-670c-4f64-9fe9-4892866d6763.mp4

## Parameters

Some configurations can be set via command line parameters. The parameters can also be viewed via `mermerd -h`

| Parameter | Description                                                       | Example                                                            |
|-----------|-------------------------------------------------------------------|--------------------------------------------------------------------|
| c         | the connection string                                             | `mermerd -c="postgresql://user:password@localhost:5432/dvdrental"` |
| s         | the schema                                                        | `mermerd -s=public`                                                |
| ac        | toggle if all constraints should be included (default false) [^1] | `mermerd -ac`                                                      |

[^1]: If the flag `-ac` is provided, mermerd will print out all constraints of the selected tables, even when the
resulting constraint is not in the list of selected tables. These tables do not have any column infos and are only
present via their table name.

## Connection strings

Examples for connection strings:

* `postgresql://user:password@localhost:5432/yourDb`
* `mysql://root:password@tcp(127.0.0.1:3306)/yourDb`

## How can I write/update Mermaid-JS diagrams?

* All information can be found here: [Mermaid-JS](https://mermaid-js.github.io/mermaid/#/entityRelationshipDiagram)
* I also recommend using an IDE with an Mermaid-JS extension,
  e.g. [VS Code](https://marketplace.visualstudio.com/items?itemName=tomoyukim.vscode-mermaid-editor)

## How does mermerd determine the constraints?

The table constraints are analysed and interpreted as listed:

| Nr. | Constraint type                        | Criteria                                                                 |
|-----|----------------------------------------|--------------------------------------------------------------------------|
| 1   | <code>a &#124;o--&#124;&#124; b</code> | If table a has a FK to table b and that column is the only PK of table a |
| 2   | <code>a }o--&#124;&#124; b</code>      | Same as 1, but table a has multiple PK                                   |
|     |                                        | Same as 1, but the FK is not a PK                                        |
| 3   | <code>a }o--o&#124; b</code>           | Same as 2, but the FK is nullable                                        |

## Roadmap

* [ ] Unit tests
* [ ] Configurable suggestions for connection string input
* [ ] Support `}o--o|` relation (currently displayed as `}o--||`)
* [ ] Improve output file naming
* [ ] Take unique constraints into account
* [ ] Support ERD Attributes for FK and PK
