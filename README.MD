# REST API based on Clean Architecture

Application will be designed around layers:

- **Transport Layer**, handles HTTP requests & JWT authentication and authorization
- **Service Layer**, handles the business logic
- **Repository Layer**, handles the DB interaction

## Tech Stack

- **Golang aka Go**
- **jwt-go** - for working with JSON Web Tokens
- **golang-migrate** - for database migrations
- **sqlx** - for easier database interaction
- **go.uuid** - for working with UUID's
- **zerolog** - for nice logs
- **testify** - for easier tests
- **Package validator**   - implements value validations for structs and individual fields based on tags 
- **go.cmp** - for comparing structs for tests