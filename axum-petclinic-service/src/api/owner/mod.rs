mod model;
pub(crate) mod service;

use axum::{Json, Router};
use axum::extract::{Path, State};
use axum::routing::get;
use tracing::info;
use crate::api::owner::model::Responses;
use crate::api::owner::service::{Service, Servicer};
use crate::repository::owner::{Repository};

pub fn routes(service: Service<Repository>) -> Router {
    Router::new()
        .route("/owner/all", get(get_all_owners))
        .route("/owner/id/:id", get(get_owner))
        .with_state(service)
}

async fn get_all_owners(
    State(service): State<Service<Repository>>) -> Json<Responses> {

    info!("Getting all owners");
    let owners = service.get_all_owners().await;
    let responses = model::Response::from_owners(owners);
    Json(Responses { owners: responses })
}

async fn get_owner(
    State(service): State<Service<Repository>>,
    Path(id): Path<i32>) -> Json<model::Response> {

    info!("Getting owner by id: {}", id);
    let owner = service.get_owner_by_id(id).await;
    let response = model::Response::from_owner(owner);
    Json(response)
}
