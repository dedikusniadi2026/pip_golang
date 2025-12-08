# TODO: Update Coverage to 100% for auth-service/service

## Tasks
- [x] Create service/auth_service_test.go with MockUserRepository
- [x] Add test for Login method - success case
- [x] Add test for Login method - invalid password
- [x] Add test for Login method - user not found error
- [x] Add test for HashPassword function
- [x] Add test for CheckPasswordHash function - valid
- [x] Add test for CheckPasswordHash function - invalid
- [x] Add test for GenerateJWT function
- [x] Add test for GenerateRefreshToken function
- [x] Run go test ./service/... to execute tests
- [x] Run go test -cover ./service/... to check coverage
- [x] Verify coverage.out shows 100% for auth_service.go
