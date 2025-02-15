use std::env;
use config::{Config, ConfigError, Environment, File};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Configuration {
    pub database: DatabaseConfiguration,
    pub application_port: u16,
}

#[derive(Debug, Deserialize)]
pub struct AppInfo {
    pub version: String,
    pub name: String,
    pub description: String,
    pub author: String,
    pub license: String,
}

#[derive(Debug, Deserialize)]
pub struct DatabaseConfiguration {
    username: String,
    password: String,
    port: u16,
    host: String,
    database_name: String,
}

impl Configuration {
    pub fn new() -> Result<Self, ConfigError> {
        let run_mode = env::var("RUN_MODE").unwrap_or_else(|_| "dev".into());

        let s = Config::builder()
            // Start off by merging in the "default" configuration file
            .add_source(File::with_name("config/app"))
            // Add in the current environment file
            // Default to 'development' env
            // Note that this file is _optional_
            .add_source(
                File::with_name(&format!("config/{}", run_mode))
                    .required(false),
            )
            // Add in a local configuration file
            // This file shouldn't be checked in to git
            .add_source(File::with_name("config/local").required(false))
            // Add in settings from the environment (with a prefix of APP)
            // Eg.. `APP_DEBUG=1 ./target/app` would set the `debug` key
            .add_source(Environment::with_prefix("app"))
            .build()?;

        // You can deserialize (and thus freeze) the entire configuration as
        s.try_deserialize()
    }
}
