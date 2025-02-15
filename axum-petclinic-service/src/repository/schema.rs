use chrono::{DateTime, Local, NaiveDate};
use sqlx::FromRow;

#[derive(Debug, Clone, FromRow)]
pub struct Owner {
    pub id: i32,
    pub first_name: String,
    pub last_name: String,
    pub address: String,
    pub city: String,
    pub telephone: String,
    pub log_date: Option<LogDate>,
    // pub pets: Vec<Pet>,
}

#[derive(Debug, Clone, FromRow)]
pub struct Pet {
    pub id: i32,
    pub name: String,
    pub birth_date: NaiveDate,
    pub owner: Owner,
    pub kind: Kind,
    pub log_date: Option<LogDate>,
}

#[derive(Debug, Clone, FromRow)]
pub struct Vet {
    pub id: i32,
    pub first_name: String,
    pub last_name: String,
    pub specialties: Vec<Specialty>,
    pub log_date: Option<LogDate>,
}

#[derive(Debug, Clone, FromRow)]
pub struct Visit {
    pub id: i32,
    pub visit_date: String,
    pub description: String,
    pub pet: Pet,
    pub owner: Owner,
    pub log_date: Option<LogDate>,
}

#[derive(Debug, Clone, FromRow)]
pub struct Specialty {
    pub id: i32,
    pub name: String,
    pub log_date: Option<LogDate>,
}

#[derive(Debug, Clone, FromRow)]
pub struct Kind {
    pub id: i32,
    pub name: String,
    pub log_date: Option<LogDate>,
}

#[derive(Debug, Clone, FromRow)]
pub struct LogDate {
    pub created_at: DateTime<Local>,
    pub updated_at : DateTime<Local>,
    pub deleted_at: Option<DateTime<Local>>,
}
