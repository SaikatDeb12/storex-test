# Asset Management System

A production-ready backend API service for managing organizational assets, users, authentication, and role-based access control.

## Features

### 1. User Authentication & Authorization
- Secure user registration and login functionality
- Passwords hashed using industry-standard algorithms
- JWT-based authentication with server-side session validation
- Role-Based Access Control (RBAC) for granular permissions
- Secure logout functionality with session invalidation

### 2. User Management
- Fetch all registered users (role-restricted)
- Retrieve user details by ID
- Role-based middleware enforcement for all operations
- Atomic database transactions for user and session creation

### 3. Asset Management
- Create new assets with structured validation
- Fetch assets through protected routes
- Update existing asset details securely
- Assign assets to users
- Mark assets as "sent to service" when required
- Comprehensive asset lifecycle management
- All critical operations protected by role-based authorization

### 4. Server & Security
- Protected routes with JWT authentication
- Multi-layer middleware validation:
  - Token signature verification
  - Session validity checking
  - Role permission enforcement
- Centralized error handling for consistent API responses
- Health check endpoint for server monitoring
- Modular and scalable architecture

## Tech Stack

- **Backend**: Go (Golang)
- **Router**: Chi
- **Database**: PostgreSQL
- **Query Layer**: sqlx
- **Authentication**: JWT
- **Architecture**: Modular layered architecture

## Project Structure

├── cmd/

│ └── main.go

├── internal/

│ ├── handler/ # Request handlers

│ ├── middleware/ # Auth & role middleware

│ ├── dbhelper/ # Database helper functions

│ ├── database/ # Database connection & setup

│ ├── models/ # Data models

│ └── utils/ # Utility functions


## API Routes

Base URL: `/v1`

### Authentication Routes (`/v1/auth`)

#### Register User

POST /v1/auth/login

**Request Body:**
```json
{
  "name": "string",
  "email": "string",
  "phoneNumber": "string",
  "role": "string",
  "employment": "string",
  "password": "string"
}
```

#### Login User

POST /v1/auth/login

**Request Body:**
```json
{
  "email": "string",
  "password": "string"
}
```

#### Logout User

POST /v1/auth/logout

### Protected Routes (Requires Authorization: Bearer <jwt_token>)

User Routes (/v1/users)

**Routes:**
```
GET /v1/users           # Fetch all users
GET /v1/users/{id}      # Fetch user by ID
```

Asset Rroutes (/v1/asset)

**Routes:**
```
GET    /v1/asset                    # Fetch all assets
POST   /v1/asset                    # Create new asset
PUT    /v1/asset/update/{id}         # Update asset details
PUT    /v1/asset/assign              # Assign asset to user
PUT    /v1/asset/service/{id}        # Mark asset for service
```

Authentication Flow
----------------------

1.  **User Registration/Login**: Server generates:
    
    -   User ID
        
    -   Session ID
        
    -   JWT Token (contains user\_id, session\_id, and role)
        
2.  **Token Storage**: Client securely stores the JWT token
    
3.  textAuthorization: Bearer
    
4.  **Middleware Validation**:
    
    -   Verifies token signature
        
    -   Validates session exists and is active
        
    -   Checks role permissions for the requested operation
        

Installation
---------------

### Prerequisites

-   Go 1.16+

-   PostgreSQL
    
-   Git

#### Steps
Clone the repository

```
git clone https://github.com/SaikatDeb12/storeX.git
cd asset-management-system
```

#### Configure environment variables
Create a .env file in the root directory:

```
DATABASE_URL=postgresql://username:password@localhost:5432/dbname?sslmode=disable
JWT_SECRET=your_jwt_secret_key
```

#### Start the server

```
go run cmd/main.go
```
