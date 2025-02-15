use sqlx::Error;

pub mod owner;
pub mod pet;
pub mod schema;


pub trait BaseQuery<T> {
    async fn find_all(&self) -> Result<Vec<T>, Error>;
    async fn find_by_id(&self, id: i32) -> Result<T, Error>;
}