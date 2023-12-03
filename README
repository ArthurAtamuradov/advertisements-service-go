# Advertisement Service

This is a Go-based REST API service for managing advertisements. The project follows the principles of Clean Architecture and uses the Gin framework for HTTP routing, MySQL for the relational database, and Redis for caching.

## Features

- **Add Advertisement:** Create a new advertisement with a title, description, price, and activation status.
- **Update Advertisement:** Modify an existing advertisement's details.
- **Delete Advertisement:** Remove an advertisement from the system.
- **View Advertisement:** Retrieve information about a specific advertisement.
- **List Advertisements:** Fetch a paginated list of all advertisements.
- **Pagination and Sorting:** Customize the number of advertisements per page and sort them by price or creation date.

## Setup

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/ArthurAtamuradov/advertisements-service-go.git
   cd advertisement-service-go
   ```

2. **Install Dependencies:**

   ```bash
   go mod download
   ```

3. **Database Setup:**

   - Create a MySQL database and configure the connection in `config.yaml`.
   - Migrations applied on start of application.

4. **Run the Service:**

   ```bash
   go run main.go
   ```

5. **Access the API:**

   The API will be available at `http://localhost:8080`.

## Configuration

The application configuration is stored in `config.yaml`. You can adjust settings such as server port, database connection, and Redis connection there.

## Docker

Run the service in a Docker container along with MySQL and Redis using `docker-compose`:

```bash
docker-compose build
docker-compose up
```
