use serde::Serialize;
use crate::repository::schema;

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct Response {
    pub id : i32,
    pub first_name : String,
    pub last_name : String,
    pub address : String,
    pub city: String,
    pub telephone : String,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct Responses {
    pub owners: Vec<Response>,
}

impl Response {
    pub fn from_owner(owner: schema::Owner) -> Self {
        Self {
            id: owner.id,
            first_name: owner.first_name,
            last_name: owner.last_name,
            address: owner.address,
            city: owner.city,
            telephone: owner.telephone,
        }
    }

    pub fn from_owners(owners: Vec<schema::Owner>) -> Vec<Response> {
        owners.into_iter().map(|owner| Response::from_owner(owner)).collect()
    }
}