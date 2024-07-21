# soidh-media-service

This is Go web server for managing media for soidh.

## Setup

1. Copy `.env.example` to `.env` and fill in your configuration.
2. Run `go mod tidy` to install dependencies.
3. Run the database migrations: `psql -U your_username -d your_database -a -f scripts/migration.sql`
4. Start the server: `go run cmd/server/main.go`

## Project Structure

- `cmd/`: Contains the main applications of the project.
- `internal/`: Houses the private application and library code.
- `pkg/`: Stores the public library code.
- `scripts/`: Holds database migrations and other scripts.

<!-- ## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests. -->

## License

This project is licensed under the MIT License - see the LICENSE file for details.