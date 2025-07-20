# Pet Clinic API

There are 4 projects in this repository

| Project Name                                         | Language | Status      |
|------------------------------------------------------|----------|-------------|
| spring-petclinic                                     | Java     | cloned      |
| kotlin-petclinic                                     | Kotlin   | cloned      |
| [axum-petclinic-service](./axum-petclinic-service)   | Rust     | in-progress |
| [fiber-petclinic-service](./fiber-petclinic-service) | Go       | in-progress |
| [gin-petclinic-service](./gin-petclinic-service)     | Go       | in-progress |


## Goal
The goal of this repository is to implement a pet clinic API in different languages and frameworks. These are endpoints:

| Method | Endpoint                     | Description                |
|--------|------------------------------|----------------------------|
| GET    | /owners?last-name={lastName} | Search owners by last name |
| POST   | /owners                      | Add new owner              |
| PUT    | /owners/{id}                 | Update owner               |
|        |                              |                            |
|        |                              |                            |
|        |                              |                            |
|        |                              |                            |
|        |                              |                            |
|        |                              |                            |
|        |                              |                            |
|        |                              |                            |
|        |                              |                            |


- GET /pets
- GET /pets/{id}
- POST /pets
- PUT /pets/{id}


## Spring Pet Clinic Wireframe

| Page              | Description                                                                                                                                 | Endpoints                        | Page Screenshot                                                              |   
|-------------------|---------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------|------------------------------------------------------------------------------|
| Home Page         | This is Pet Clinic Home Page to display the pet photos.  This page will not invoke any API. Click on pet image to display Search Owner Page |                                  | ![01-home-page.png](pet-clinic-page-screnshots/01-home-page.png)             |
| Search Owner Page | Search Owner Page provides a Last Name text box to search owner(s) by last name.                                                            | GET /owners?last-name={lastname} | ![02-find-owners.png](pet-clinic-page-screnshots/02-find-owners.png)         |
|                   |                                                                                                                                             |                                  | ![02-find-owners.png](pet-clinic-page-screnshots/02-find-owners.png)         |
|                   |                                                                                                                                             |                                  | ![03-owner-page.png](pet-clinic-page-screnshots/03-owner-page.png)           |
|                   |                                                                                                                                             | POST /owners                     | ![05-add-owner-page.png](pet-clinic-page-screnshots/05-add-owner-page.png)   |
|                   |                                                                                                                                             |                                  | ![05-update-pet-page.png](pet-clinic-page-screnshots/05-update-pet-page.png) |
|                   |                                                                                                                                             |                                  | ![06-add-visit-page.png](pet-clinic-page-screnshots/06-add-visit-page.png)   |
|                   |                                                                                                                                             |                                  | ![07-update-pet-page.png](pet-clinic-page-screnshots/07-update-pet-page.png) |
|                   |                                                                                                                                             |                                  | ![08-error-page.png](pet-clinic-page-screnshots/08-error-page.png)           |
|                   |                                                                                                                                             |                                  |                                                                              |








