mod api;
mod configuration;
mod db;
pub mod error;
mod middleware;
mod repository;

use std::time::Duration;
use axum::{Router};
use tokio::net::TcpListener;
use tokio::signal;
use tower_http::timeout::TimeoutLayer;
use tower_http::trace::TraceLayer;
use tracing::log::debug;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};
use tracing_subscriber::fmt::layer;
use api::{health, info, owner, pet};


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
    let pool = db::create_connection().await;
    let repo = repository::owner::Repository::new(pool);
    let service = owner::service::Service::new(repo);

    // Create a new router
    let app_routes =
        Router::new()
            .merge(health::routes())
            .merge(info::routes())
            .merge(owner::routes(service))
            .merge(pet::routes())
            .layer((
                TraceLayer::new_for_http(),
                // Graceful shutdown will wait for outstanding requests to complete.
                // Add a timeout so requests don't respond forever.
                TimeoutLayer::new(Duration::from_secs(10))),);

    let listener = TcpListener::bind("127.0.0.1:3003")
        .await
        .unwrap();
    debug!("Listening on {}", listener.local_addr().unwrap());

    // Start the server
    axum::serve(listener, app_routes)
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}

/// Wait for a shutdown signal to gracefully shut down the server.
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
