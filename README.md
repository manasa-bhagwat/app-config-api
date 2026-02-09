# App Config API

A simple REST API built with Go (Gin) and PostgreSQL to manage mobile app pages and widgets.

Acts as a configuration backend that lets apps render UI dynamically.

---

## Tech Stack

* **Go** (Gin) - lightweight, high-performance HTTP server
* **PostgreSQL** - reliable relational database
* **UUID** - globally unique identifiers for entities
* **JSONB** - flexible storage for dynamic widget configuration
* **Raw SQL** - clear, explicit queries with minimal abstraction

---

## Features

### Pages

- Create page
- List all pages
- Get single page (with widgets)
- Update page
- Delete page (home page protected)

### Widgets

- Add widget to a page
- Update widget
- Delete widget
- Reorder widgets on a page

### Validations

- Page name required
- Route must be unique
- Only one home page allowed
- Home page cannot be deleted
- Widget types allowed:
    - banner
    - product_grid
    - text
- Widget config stored as JSONB
- Consistent error response format

---

## Architecture

This service stores:

- Pages (app screens)
- Widgets (UI components)
- Widget layout order

Mobile apps fetch this configuration and render UI dynamically.


Below is a high-level view of how pages, widgets, and layout configuration flow through the system.

![Architecture](images/architecture.png)

### Full Design Overview

The original Excalidraw file is included for reference and can be edited if needed:

You can open it by:
1. Visiting https://excalidraw.com
2. Clicking "Open"
3. Uploading the `.excalidraw`

---

## Project Structure

```
app-config-api/
├── design/
├── images/
├── postman/
├── main.go
├── db.go
├── models.go
├── handlers.go
├── routes.go
├── middleware.go
├── schema.sql
└── README.md
```

---

## Run Locally

### 1) Clone repo

```
git clone https://github.com/manasa-bhagwat/app-config-api.git
cd app-config-api
```

### 2) Create database

Using CLI:

```
createdb configapp
psql -d configapp -f schema.sql
```

Or using pgAdmin:

- Create DB: `configapp`
- Open Query tool.
- Run `schema.sql`

### 3) Update DB connection

In `db.go`:

```
postgres://username:password@localhost:5432/configapp?sslmode=disable
```

### 4) Run server

```
go run .
```

Server runs at:

```
http://localhost:8080
```

---

## API Endpoints

### Pages

- **POST** `/pages`
- **GET** `/pages`
- **GET** `/pages/:id`
- **PUT** `/pages/:id`
- **DELETE** `/pages/:id`

### Widgets

- **POST** `/pages/:id/widgets`
- **PUT** `/widgets/:id`
- **DELETE** `/widgets/:id`
- **POST** `/pages/:id/widgets/reorder`

---

## Postman Test Flow

**Base URL**

```
http://localhost:8080
```

### 1) Create Home Page

**POST /pages**

```json
{
  "name": "Home",
  "route": "/home",
  "is_home": true
}
```

Save: `HOME_PAGE_ID`

---

### 2) Create Products Page

**POST /pages**

```json
{
  "name": "Products",
  "route": "/products",
  "is_home": false
}
```

Save: `PRODUCTS_PAGE_ID`

---

### 3) Add Widgets to Home Page

**POST /pages/HOME_PAGE_ID/widgets**

Banner:

```json
{
  "type": "banner",
  "position": 1,
  "config": {
    "image_url": "https://store.com/sale.jpg",
    "title": "Mega Sale",
    "cta": "Shop Now"
  }
}
```

Product Grid:

```json
{
  "type": "product_grid",
  "position": 2,
  "config": {
    "collection": "featured",
    "columns": 2
  }
}
```

Text:

```json
{
  "type": "text",
  "position": 3,
  "config": {
    "content": "Welcome to our store!"
  }
}
```

---

### 4) Fetch Full Page Config

**GET /pages/HOME_PAGE_ID**

Returns:

- Page details
- All widgets
- Ordered by position

---

### 5) Reorder Widgets

**POST /pages/HOME_PAGE_ID/widgets/reorder**

```json
{
  "widget_ids": [
    "TEXT_WIDGET_ID",
    "BANNER_WIDGET_ID",
    "GRID_WIDGET_ID"
  ]
}
```

---

### 6) Update Widget

**PUT /widgets/WIDGET_ID**

```json
{
  "type": "banner",
  "position": 1,
  "config": {
    "title": "BIGGEST SALE"
  }
}
```

---

### 7) Delete Widget

**DELETE /widgets/WIDGET_ID**

---

### Postman Collection

A ready-to-use Postman collection export is included in the repository for quick testing of all API endpoints.

File location:

```
postman/app-config-api.postman_collection.json
```

To use it:

1. Open Postman
2. Click **Import**
3. Select the JSON file from the `postman/` folder
4. Set the base URL if needed:

```
http://localhost:8080
```

This collection contains pre-configured requests for:

* Pages (Create, List, Get, Update, Delete)
* Widgets (Add, Update, Delete)
* Reorder widgets

---
## Request Logging

Each request logs:

- Method
- Path
- Status
- Response time


Example:

```
POST /pages | 201 | 2ms
GET /pages | 200 | 1ms
```
![Logging](logging.png)
---
