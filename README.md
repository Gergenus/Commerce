# e-commerce web platform for trading.
## Structure of the project
- user-service
- product-service
- cart-service
- order-serivce
- notification-service
## Each of them accounts for their specific task using up-to-date tech stack. Here is a break down
- user-service handles user registration, authentication, and profile management.
- product-service manages product listings, categories, and inventory.
- cart-service manages users’ shopping carts, including adding/removing items and updating quantities.
- order-service processes orders, places orders etc.
- notification-service sends emails for profile confirmation.
## Techiques used in the project
- clean architecture
- saga pattern
- rate limiting
## Tech stack used in the project
- Golang + echo framework
- Postgresql + pgx package
- Redis
- JWT
- Kafka
- Faktory
- Elasticsearch
- gRPC
- Mockery
## Testing
Mockery was used to create mocks on service and repository layers. product-service has essential unit tests.
