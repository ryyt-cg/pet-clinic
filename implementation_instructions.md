# Endpoint Implementation
Two approaches: top-down approach or botton-up approach.  My choice is top-down approach because I can write and test the code easily and quickly.

## Top-down approach
First is to implement the router.  For example, implement the `product` endpoint.
1. create product package under api directory.
2. create product_router.go under `product` package
3. define `Router` struct with service as Servicer interface.  This service acts as dependency service
4. define func 'NewRouter' as a way to instantiate the product Router object
5. define func(receiver) Register() to register endpoints

`product_router.go`
```go
// create router type with  service as dependency injection
type Router struct {
	service Servicer
	// add more services below
}

// NewRouter - Product Router constructor
func NewRouter(service Servicer) *Router {
	return &Router{service}
}

// Register registers the router to the gin engine
func (r *Router) Register(routerGroup *gin.RouterGroup) {
	routerGroup.GET(":id", r.getProductByID)
	routerGroup.POST("", r.addNewProduct)
	routerGroup.PUT(":id", r.updateProduct)
}

func(r *Router) getProductByID(c *gin.Context) {
	// implementation here
}

func(r *Router) addNewProduct(c *gin.Context) {
    // implementation here
}

func(r *Router) updateProduct(c *gin.Context) {
    // implementation here
}
```

Second is to implement `product` service.
1. create product_service.go under product package
2. define Servicer interface with some implementations.  This service will be injected into Product Router
3. define Service struct with some dependencies like repository or gateway
4. define func NewService() as a constructor function to instantiate product service
5. implement the service functions


```go
// Servicer - owner service interface
type Servicer interface {
	retrieveProductByID(id uint) (*response, error)
	create(productRequest *addRequest) (*response, error)
	update(id uint, updateProduct *updateRequest) (*updateResponse, error)
}

type Service struct {
	repository repository.ProductRepositorier
	// or gateway gateway.ProductGatewayer
}

func NewService(repository repository.ProductRepositorier) *Service {
	return &Service{repository: repository}
}

// Define the functions
func (service *Service) retrieveProductByID(id uint) (*response, error) {}
func (service *Service) create(productRequest *addRequest) (*response, error) {}
func (service *Service) update(id uint, request *updateRequest) (*updateResponse, error) {}
```

Third is to implement `product` repository (database) or gateway (downstream API).  The repository and gateway functions are usually shared logic service.  I create repository package under internal directory.
1. create product_repository.go under repository package
2. define Product interface as ProductRepositorier
3. define ProductRepository struct
4. define func NewProductRepository
5. implement repository implementations

```go
type ProductRepositorier interface {
	FindByID(id uint) (*Product, error)
	Insert(product *Product) (*Product, error)
	Update(product *Product) (*Product, error)
}

// ProductRepository searches products from the database
type ProductRepository struct {
	gdb *gorm.DB
}

// NewProductRepository - ProductRepository factory
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		gdb: db,
	}
}

// implement repository functions here
func (repository *ProductRepository) FindById(id uint) (*Product, error) {}
func (repository *ProductRepository) Insert(product *Product) (*Product, error) {}
func (repository *ProductRepository) Update(product *Product) (*Product, error) {}
```

