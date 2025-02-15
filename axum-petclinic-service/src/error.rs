use axum::http::StatusCode;
use axum::response::IntoResponse;

#[derive(Debug, Clone)]
struct ApiError {
    status: StatusCode,
    message: String,
    timestamp: String,
}

// impl IntoResponse for ApiError {
//     fn into_response(self) -> axum::http::Response<Self::Body> {
//         axum::http::Response::builder()
//             .status(self.status)
//             .header("Content-Type", "application/json")
//             .unwrap()
//     }
// }