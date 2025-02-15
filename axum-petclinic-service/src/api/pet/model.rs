use chrono::NaiveDate;
use serde::{Deserialize, Serialize};
use crate::repository::schema;

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct Response {
    pub id : i32,
    pub name: String,
    pub birth_date: NaiveDate,
    pub kind: String,
}

impl Response {
    fn from(pet: schema::Pet) -> Self {
        Self {
            id: pet.id,
            name: pet.name,
            birth_date: pet.birth_date,
            kind: pet.kind.name,
        }
    }
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct Request {
    name: String,
    birthday: String,
    kind_id: u64,
    kind_name: String,
}
