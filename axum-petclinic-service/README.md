# Pet Clinic REST API using [Axum](https://www.youtube.com/watch?v=Wnb_n5YktO8&t=716s) & [Tokio](https://tokio.rs/)

## Introduction
* Tokio is an asynchronous runtime for the Rust programming language. It provides the building blocks needed for writing network applications. It gives you the ability to write fast, reliable, and highly concurrent applications.
* Axum is a web framework built on top of Tokio and Tower. It is designed to be both ergonomic and performant. It is built on top of hyper, a fast and low-level HTTP implementation in Rust.

[Send + Sync](https://www.youtube.com/watch?v=yOezcP-XaIw)

## Technologies Stacks
* [Rust](https://www.rust-lang.org/) - programming language
* [Axum](https://docs.rs/axum/latest/axum/) - web framework
* [Clippy](https://github.com/rust-lang/rust-clippy) - linter
* [Derive Builder](https://github.com/colin-kiegel/rust-derive-builder) - builder pattern
* [Mockall](https://github.com/asomers/mockall) - mocking
* [sqlx](https://github.com/launchbadge/sqlx) - async SQL client
* [Serde](https://serde.rs/) - serialization
* [Tower HTTP](https://github.com/tower-rs/tower-http) - middleware
* [Tokio](https://tokio.rs/) - async runtime
* [tarpaulin](https://github.com/xd009642/tarpaulin) - code coverage
* [Validation](https://github.com/Keats/validator) - validation
* [zld](https://github.com/michaeleisel/zld) - faster linker

## Create a new project
Create a new project using the following cargo commands.

```bash
cargo new axum-petclinic-service --bin
cd axum-petclinic-service
```
Add tokio, serde, tower-http and axum dependencies to Cargo.toml

```toml
axum = { version = "0.7.5", features = ["macros"] }
tokio = { version = "1.37.0", features = ["full"] }
serde = { version = "1.0.199", features = ["derive"] }
tower-http = { version = "0.5.2" }

[dev-dependencies]
anyhow = { version = "1.0.82" }
httpc-test = { version = "0.1.9" }
```

## Simple health & info REST API
Create a module named api and add health and info modules to it.

health module
```rust
use axum::{Json, Router};
use axum::routing::get;
use serde::Serialize;

#[derive(Debug, Serialize)]
struct HealthCheck {
    status: String,
}

pub fn routes() -> Router {
    // let app_state = AppState { mc };
    Router::new().route("/health", get(check))
}

async fn check() -> Json<HealthCheck> {
    let health_check = HealthCheck {
        status: "UP".to_string(),
    };

    Json(health_check)
}

```

info module
```rust
use axum::{Json, Router};
use axum::routing::get;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize)]
struct InfoResponse {
    pub version: String,
    pub name: String,
    pub description: String,
    pub author: String,
    pub license: String,
}

pub fn routes() -> Router {
    // let app_state = AppState { mc };
    Router::new().route("/info", get(info))
}

async fn info() -> Json<InfoResponse> {
    let info = InfoResponse {
        version: "0.1.0".to_string(),
        name: "Rust API".to_string(),
        description: "Rust API using Axum".to_string(),
        author: "Ethan J. Tran".to_string(),
        license: "MIT".to_string(),
    };

    Json(info)
}
```

main.rs
```rust
mod api;
mod db;
pub mod error;
mod middleware;
mod repository;

use std::time::Duration;
use axum::{Router};
use sqlx::postgres::PgPoolOptions;
use tokio::net::TcpListener;
use tokio::signal;
use tower_http::timeout::TimeoutLayer;
use tower_http::trace::TraceLayer;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};
use api::{health, info, owner};


#[tokio::main]
async fn main() {
    // Setup tracing
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "example_tls_graceful_shutdown=debug".into()),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    // Create a connection pool to the database.
    let conn_string = "postgresql://postgres:mysecretpassword@localhost:5432/petclinic?sslmode=disable";
    let pool = PgPoolOptions::new()
        .min_connections(1)
        .max_connections(5)
        .connect(conn_string)
        .await
        .unwrap();

    // Create a new router
    let app_routes =
        Router::new()
            .merge(health::routes())
            .merge(info::routes())
            .merge(owner::routes(pool.clone()))
            .layer((
                       TraceLayer::new_for_http(),
                       // Graceful shutdown will wait for outstanding requests to complete.
                       // Add a timeout so requests don't respond forever.
                       TimeoutLayer::new(Duration::from_secs(10))),);

    let listener = TcpListener::bind("127.0.0.1:3003")
        .await
        .unwrap();
    println!("Listening on {}", listener.local_addr().unwrap());

    // Start the server
    axum::serve(listener, app_routes)
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}

/// Wait for a shutdown signal to gracefully shutdown the server.
async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };

    #[cfg(unix)]
        let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
        let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }
}
```

### Run the application
Run the application using the following commands.
* 1st terminal:
```bash
make watch-src-change
```

### Test the application
Test the application using the following commands.
* 2nd terminal:
```bash
make watch-test-change
```

## Step 2 - Refactor Connection Pool
Create db module and refactor the connection pool code.

db/mod.rs
```rust
use sqlx::{Pool, Postgres};
use sqlx::postgres::PgPoolOptions;

pub async fn create_connection() -> Pool<Postgres> {
    let db_url = "postgresql://postgres:mysecretpassword@localhost:5432/petclinic?sslmode=disable";
    PgPoolOptions::new()
        .min_connections(1)
        .max_connections(5)
        .connect(db_url)
        .await.unwrap()
}
```
In main.rs replace the connection pool code with the following code.
```rust
    // Create a connection pool to the database.
    let pool = db::create_connection().await;
```

## Step 3 - Abstract the Database Queries into Repository
* Create repository module
* Create owner module
* Move the database queries into the Owner repository module

```rust
use sqlx::{Error, Pool, Postgres, Row};
use sqlx::postgres::PgRow;
use crate::repository::schema::{LogDate, Owner};

pub trait Repositorier {
    fn new(db: Pool<Postgres>) -> Self;
    async fn find_by_id(&self, id: i32) -> Result<Owner, Error>;
}

#[derive(Debug, Clone)]
pub struct Repository {
    pub db: Pool<Postgres>,
}

impl Repositorier for Repository {
   fn new(db: Pool<Postgres>) -> Self {
        Repository {
            db
        }
    }

    async fn find_by_id(&self, id: i32) -> Result<Owner, Error> {
        let owner_query = sqlx::query(
            "SELECT id, first_name, last_name, address, city, telephone, created_at, updated_at, deleted_at \
                    FROM owners \
                    WHERE id = $1")
            .bind(id);

        let owner: Owner = owner_query.map(|row: PgRow| Owner {
            id: row.get("id"),
            first_name: row.get("first_name"),
            last_name: row.get("last_name"),
            address: row.get("address"),
            city: row.get("city"),
            telephone: row.get("telephone"),
            log_date: LogDate {
                created_at: row.get("created_at"),
                updated_at: row.get("updated_at"),
                deleted_at: row.get("deleted_at"),
            }
        }).fetch_one(&self.db).await.unwrap();

        Ok(owner)
    }
}
```
In main.rs, replace the owner module with the following code.
```rust
    // Create a connection pool to the database.
    let pool = db::create_connection().await;
    let repo = Repository::new(pool);
    
    ...

    .merge(owner::routes(repo.clone()))
    
  ```

## Step 4 - Create a Service Module for Business Logic
* Create a service module for Owner
* Inject the repository into the service module
* Inject the service into routes

```rust
```


## List of example repositories

| Description                           | Link                                                                               |
|---------------------------------------|------------------------------------------------------------------------------------|
| Error handling in Rust - a Deep Dive  | https://www.lpalmieri.com/posts/error-handling-rust/                               |
| zero-to-production in Rust            | https://github.com/jeastham1993/zero-to-production-rust                            |
|                                       | https://github.com/LukeMathWalker/zero-to-production                               |
| http client - reqwest                 | https://www.youtube.com/watch?v=j9MsMYz9hBw&t=656s                                 |
|                                       |                                                                                    |
| Shelter Project - Base microservice   | https://github.com/sapati/shelter-project/tree/ep-12                               |
| Clean Architecture                    | https://github.com/MSC29/clean-architecture-rust                                   |
| Fullstack clean architecture          | https://github.com/flosse/clean-architecture-with-rust                             |
| Hexagonal Architecture                | https://jameseastham.co.uk/post/software-development/hexagaonal-architecture-rust/ |
| Clean and scalable architecture       | https://kerkour.com/rust-web-application-clean-architecture                        |
| Black Hat                             | https://github.com/skerkour/black-hat-rust                                         |
|                                       |                                                                                    |
| Polymorphism with Traits              | https://www.youtube.com/watch?v=CHRNj5oubwc                                        |
| Implement Rust Traits                 | https://youtube.com/watch?v=Lrayq0UW7nA                                            |
| Dynamic Dispatch                      | https://www.youtube.com/watch?v=3biW5NkNnrk&t=455s                                 |
| Store Data on the Heap with Box       | https://www.youtube.com/watch?v=br6nGvqT48w                                        |
| Axum Diesel real world                | https://github.com/Quentin-Piot/axum-diesel-real-world/tree/master                 |
| Stock Metrics                         | https://github.com/yuk1ty/stock-metrics                                            |
| Mock multiple impls of generic struct | https://github.com/asomers/mockall/issues/271                                      |
| Jeremy Chone Github                   | https://github.com/jeremychone?tab=repositories                                    |
|                                       |                                                                                    |
| CRUD API                              | https://www.youtube.com/watch?v=NJsTgmayHZY&t=556s                                 |
| Error Handler Patterns                | https://www.youtube.com/watch?v=f82wn-1DPas                                        |
|                                       | https://www.youtube.com/watch?v=kHxjiTv8r18                                        |


https://github.com/krojew/springtime/tree/master/springtime-web-axum

https://www.reddit.com/r/rust/comments/15c68sk/seeking_guidance_on_writing_a_rust_api_with_ddd/

https://stackoverflow.com/questions/76497219/handling-dependency-injection-and-errors-in-rust-using-axum

https://users.rust-lang.org/t/realworld-axum-sqlx-refactoring-and-unit-testing-for-an-even-more-real-world/97607/3

https://github.com/brannan/realworld-axum-sqlx/tree/main

https://github.com/tokio-rs/axum/discussions/358

https://audunhalland.github.io/blog/

https://audunhalland.github.io/blog/testability-reimagining-oop-design-patterns-in-rust/

