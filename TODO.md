# TODO

## Fix Login Error "invalid credentials password"

### Problem
- Login fails with "invalid credentials password" because passwords in DB are not hashed with bcrypt.

### Solution
- Added Register endpoint to create users with properly hashed passwords.
- Fixed AuthService creation to use NewAuthService with tokenRepo.

### Steps Completed
- [x] Add Save to UserRepository interface
- [x] Add Register to AuthServiceInterface
- [x] Add Register method to AuthService
- [x] Add RegisterRequest to dto.go
- [x] Add Register to AuthHandler
- [x] Add /register route in server.go
- [x] Fix authService creation in server.go

### Next Steps
- [ ] Register new users using POST /register with username, password, role
- [ ] For existing users, update passwords in DB to bcrypt hashes or re-register them
- [ ] Test login with properly hashed passwords

### Testing
- Run the server: go run main.go
- Register a user: curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"username":"admin","password":"password123","role":"ADMIN"}'
- Login: curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"admin","password":"password123"}'
