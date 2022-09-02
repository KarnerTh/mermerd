# Mermerd
[![tests](https://github.com/KarnerTh/mermerd/actions/workflows/test.yml/badge.svg)](https://github.com/KarnerTh/mermerd/actions/workflows/test.yml)

Create [Mermaid-Js](https://mermaid-js.github.io/mermaid/#/entityRelationshipDiagram) ERD diagrams from existing tables.

Want to see what has changed? Take a look at
the [Changelog](https://github.com/KarnerTh/mermerd/blob/master/changelog.md)

# Contents
<ul>
  <li><a href="#installation">Installation</a></li>
  <li><a href="#features">Features</a></li>
  <li><a href="#why-would-i-need-it--why-should-i-care">Why would I need it / Why should I care?</a></li>
  <li><a href="#how-does-it-work">How does it work</a></li>
  <li><a href="#parametersflags">Parameters/Flags</a></li>
  <li><a href="#global-configuration-file">Global configuration file</a></li>
  <li><a href="#use-a-predefined-run-configuration-eg-for-cicd">Use a predefined run configuration (e.g. for CI/CD)</a></li>
  <li><a href="#example-usages">Example usages</a></li>
  <li><a href="#connection-strings">Connection strings</a></li>
  <li><a href="#how-can-i-writeupdate-mermaid-js-diagrams">How can I write/update Mermaid-JS diagrams?</a></li>
  <li><a href="#how-does-mermerd-determine-the-constraints">How does mermerd determine the constraints?</a></li>
  <li><a href="#tests">Tests</a></li>
  <li><a href="#roadmap">Roadmap</a></li>
</ul>

## Installation
```sh
go install github.com/KarnerTh/mermerd@latest
```

or

just head over to the [Releases](https://github.com/KarnerTh/mermerd/releases) page and download the right executable
for your operating system. To be able to use it globally on your system, add the executable to your path.

## Features

* Supports PostgreSQL MySQL, and MSSQL
* Select from available schemas
* Select only the tables you are interested in
* Show only the constraints that you are interested in
* Interactive cli (multiselect, search for tables and schemas, etc.)
* Use it in CI/CD pipeline via a run configuration
* Either generate plain mermaid syntax or enclose it with mermaid backticks to use directly in e.g. GitHub markdown

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

## Parameters/Flags

Some configurations can be set via command line parameters/flags. The available options can also be viewed
via `mermerd -h`

```
  -c, --connectionString string       connection string that should be used
      --debug                         show debug logs        
  -e, --encloseWithMermaidBackticks   enclose output with mermaid backticks (needed for e.g. in markdown viewer)
  -h, --help                          help for mermerd
      --omitConstraintLabels          omit the constraint labels
  -o, --outputFileName string         output file name (default "result.mmd")
      --runConfig string              run configuration (replaces global configuration)
  -s, --schema string                 schema that should be used
      --selectedTables strings        tables to include
      --showAllConstraints            show all constraints, even though the table of the resulting constraint was not selected
      --useAllTables                  use all available tables
```

If the flag `--showAllConstraints` is provided, mermerd will print out all constraints of the selected tables, even when
the resulting constraint is not in the list of selected tables. These tables do not have any column info and are only
present via their table name.

## Global configuration file

Mermerd uses a yaml configuration file in your home directory called `.mermerd` (needs to be created, an example is
shown below). You can set options that you want by default (e.g. enabling the `showAllConstraints` for all runs) or
provide connection string suggestions for the cli.

```yaml
showAllConstraints: true
encloseWithMermaidBackticks: true
outputFileName: "my-db.mmd"
debug: false
omitConstraintLabels: false

# These connection strings are available as suggestions in the cli (use tab to access)
connectionStringSuggestions:
  - postgresql://user:password@localhost:5432/yourDb
  - mysql://root:password@tcp(127.0.0.1:3306)/yourDb
  - sqlserver://user:password@localhost:1433?database=yourDb
```

## Use a predefined run configuration (e.g. for CI/CD)

You can specify all parameters needed for generating the ERD via a run configuration. Just create a yaml file (example
shown below) and start mermerd via `mermerd --runConfig yourRunConfig.yaml`

```yaml
# Connection properties
connectionString: "postgresql://user:password@localhost:5432/yourDb"
schema: "public"

# Define what tables should be used
useAllTables: true
# or
selectedTables:
  - city
  - customer

# Additional flags
showAllConstraints: true
encloseWithMermaidBackticks: true
outputFileName: "my-db.mmd"
debug: true
omitConstraintLabels: true
```

## Example usages

```bash
# all parameters are provided via the interactive cli
mermerd

# same as previous one, but show all constraints even though the table of the resulting constraint was not selected
mermerd --showAllConstraints

# ERD is created via the provided run config
mermerd --runConfig yourRunConfig.yaml

# specify all connection properties so that only the table selection is done via the interactive cli
mermerd -c "postgresql://user:password@localhost:5432/yourDb" -s public

# same as previous one, but use all available tables without interaction
mermerd -c "postgresql://user:password@localhost:5432/yourDb" -s public --useAllTables

# same as previous one, but use a list of tables without interaction
mermerd -c "postgresql://user:password@localhost:5432/yourDb" -s public --selectedTables article,article_label
```

## Connection strings

Examples of valid connection strings:

* `postgresql://user:password@localhost:5432/yourDb`
* `mysql://root:password@tcp(127.0.0.1:3306)/yourDb`
* `sqlserver://user:password@localhost:1433?database=yourDb` 

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

## Tests

You can either use the Makefile targets to run the tests and have a pretty
output ([tparse](https://github.com/mfridman/tparse) is used for formatting) or start them manually via `go test`.

Mocks for unit tests are generated via [mockery](https://github.com/vektra/mockery) (can be created
via `make create-mocks` or `mockery --all`)

Local setup for integration tests:

1. `cd test`
2. `docker-compose up -d`
3. done - the required tables are created automatically at startup (`test/db-table-setup.sql` contains all test data)

Integration and unit tests are separated via the `--short` flag:

* `go test --short -v ./...` runs all unit tests
* `go test -v ./...` runs all unit and integration tests

or via the Makefile targets

* `make test-unit` runs all unit tests
* `make test-all` runs all unit and integration tests

## Roadmap

* [ ] Support `}o--o|` relation (currently displayed as `}o--||`)
* [ ] Take unique constraints into account
* [ ] Support ERD Attributes for FK and PK
