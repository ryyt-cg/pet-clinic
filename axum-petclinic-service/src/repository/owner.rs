use sqlx::{Error, Pool, Postgres, Row};
use sqlx::postgres::{PgRow};
use crate::repository::BaseQuery;
use crate::repository::schema::{LogDate, Owner};

// pub trait Repositorier: BaseQuery<Owner> {}

#[derive(Clone)]
pub struct Repository {
    pub db: Pool<Postgres>,
}

impl Repository {
    pub fn new(db: Pool<Postgres>) -> Self {
        Repository {
            db,
        }
    }
}

impl BaseQuery<Owner> for Repository {
    async fn find_all(&self) -> Result<Vec<Owner>, Error> {
        let owners_query =
            sqlx::query("SELECT id, first_name, last_name, address, city, telephone, created_at, updated_at, deleted_at
                        FROM owners");

        let owners: Vec<Owner> = owners_query.map(|row: PgRow| Owner {
            id: row.get("id"),
            first_name: row.get("first_name"),
            last_name: row.get("last_name"),
            address: row.get("address"),
            city: row.get("city"),
            telephone: row.get("telephone"),
            log_date: Some(LogDate {
                created_at: row.get("created_at"),
                updated_at: row.get("updated_at"),
                deleted_at: row.get("deleted_at"),
            }),
        }).fetch_all(&self.db).await.unwrap();

        Ok(owners)
    }

    async fn find_by_id(&self, id: i32) -> Result<Owner, Error> {
        let owner_query =
            sqlx::query("SELECT id, first_name, last_name, address, city, telephone, created_at, updated_at, deleted_at
                        FROM owners
                        WHERE id = $1")
            .bind(id);

        let owner: Owner = owner_query.map(|row: PgRow| Owner {
            id: row.get("id"),
            first_name: row.get("first_name"),
            last_name: row.get("last_name"),
            address: row.get("address"),
            city: row.get("city"),
            telephone: row.get("telephone"),
            log_date: Some(LogDate {
                created_at: row.get("created_at"),
                updated_at: row.get("updated_at"),
                deleted_at: row.get("deleted_at"),
            }),
        }).fetch_one(&self.db).await.unwrap();

        Ok(owner)
    }
}