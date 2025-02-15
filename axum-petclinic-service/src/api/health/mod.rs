use axum::{Json, Router};
use axum::routing::get;
use serde::Serialize;
use tracing::{info};


#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
struct HealthCheck {
    status: String,
}

pub fn routes() -> Router {
    // let app_state = AppState { mc };
    Router::new().route("/health", get(check))
}

async fn check() -> Json<HealthCheck> {
    info!("Checking health");
    let health_check = HealthCheck {
        status: "UP".to_string(),
    };

    Json(health_check)
}