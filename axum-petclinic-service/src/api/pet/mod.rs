mod model;
mod service;

use axum::extract::{Path};
use axum::{Json, Router};
use axum::routing::get;
use chrono::NaiveDate;
use tracing::info;
use crate::api::owner::service::{Service, Servicer};
use crate::repository::owner::Repository;
use crate::repository::schema;

pub fn routes() -> Router {
    Router::new()
        // .route("/pet/all", get(crate::api::owner::get_all_pets))
        .route("/pet/id/:id", get(get_pet))
}

async fn get_pet(
    Path(id): Path<i32>) -> Json<model::Response> {

    info!("Getting pet by id: {}", id);
    let response = model::Response {
        id: 1,
        name: "Fido".to_string(),
        birth_date: NaiveDate::from_ymd_opt(2010, 5, 1)?,
        kind: "Dog".to_string(),
    };
    // let owner = service.get_owner_by_id(id).await;
    // let response = model::Response::from_owner(owner);
    Json(response)
}
