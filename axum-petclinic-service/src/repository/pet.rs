use sqlx::{Error, Pool, Postgres, Row};
use crate::repository::{BaseQuery, schema};
use crate::repository::schema::{Kind, LogDate, Owner, Pet};

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

impl BaseQuery<Pet> for Repository {
    async fn find_all(&self) -> Result<Vec<Pet>, Error> {
        let pets_query =
            sqlx::query("SELECT id, name, k.id, k.name created_at, updated_at, deleted_at
                        FROM pets");

        let pets: Vec<Pet> = pets_query.map(|row: sqlx::postgres::PgRow| Pet {
            id: row.get("id"),
            name: row.get("name"),
            birth_date: Default::default(),
            owner: Owner {
                id: Default::default(),
                first_name: Default::default(),
                last_name: Default::default(),
                address: Default::default(),
                city: Default::default(),
                telephone: Default::default(),
                log_date: Default::default(),
            },
            kind: Kind {
                id: row.get("k.id"),
                name: row.get("k.name"),
                log_date: Default::default(),
            },
            log_date: Some(LogDate {
                created_at: row.get("created_at"),
                updated_at: row.get("updated_at"),
                deleted_at: row.get("deleted_at"),
            })
        }).fetch_all(&self.db).await.unwrap();

        Ok(pets)
    }

    async fn find_by_id(&self, id: i32) -> Result<Pet, Error> {
        let pet_query =
            sqlx::query("SELECT id, name, k.id, k.name,
            o.first_name, o.last_name, created_at, updated_at, deleted_at
                        FROM pets p
                        JOIN kind k ON p.kind_id = k.id
                        JOIN owners o ON p.owner_id = o.id
                        WHERE id = $1")
            .bind(id);

        let pet: Pet = pet_query.map(|row: sqlx::postgres::PgRow| Pet {
            id: row.get("id"),
            name: row.get("name"),
            birth_date: Default::default(),
            owner: Owner {
                id: row.get("o.id"),
                first_name: row.get("o.first_name"),
                last_name: row.get("o.last_name"),
                address: Default::default(),
                city: Default::default(),
                telephone: Default::default(),
                log_date: Default::default(),
            },
            kind: Kind {
                id: row.get("k.id"),
                name: row.get("k.name"),
                log_date: Default::default(),
            },
            log_date: Some(LogDate {
                created_at: row.get("created_at"),
                updated_at: row.get("updated_at"),
                deleted_at: row.get("deleted_at"),
            }),
        }).fetch_one(&self.db).await.unwrap();

        Ok(pet)
    }
}