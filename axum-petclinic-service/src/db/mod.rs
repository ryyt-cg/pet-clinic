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