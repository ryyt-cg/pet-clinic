use crate::repository;
use crate::repository::schema::Owner;

#[mockall::automock]
pub trait Servicer {
    async fn get_all_owners(&self) -> Vec<Owner>;
    async fn get_owner_by_id(&self, id: i32) -> Owner;
}

#[derive(Clone, Debug)]
pub struct Service<R> {
    repo: R,
}

impl<R> Service<R> {
    pub fn new(repo: R) -> Self {
        Service {
            repo,
        }
    }
}

impl<R> Servicer for Service<R> where R: repository::BaseQuery<Owner> {
    async fn get_all_owners(&self) -> Vec<Owner> {
        self.repo.find_all().await.unwrap()
    }

    async fn get_owner_by_id(&self, id: i32) -> Owner {
        self.repo.find_by_id(id).await.unwrap()
    }
}

#[cfg(test)]
mod tests {

}