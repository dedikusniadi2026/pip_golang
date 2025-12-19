package mock_auth_service

type MockAuthService struct {
	RegisterFn     func(username, password, role string) error
	LoginFn        func(username, password string) (string, string, error)
	RefreshTokenFn func(refreshToken string) (string, error)
}

func (m *MockAuthService) Register(username, password, role string) error {
	if m.RegisterFn != nil {
		return m.RegisterFn(username, password, role)
	}
	return nil
}

func (m *MockAuthService) Login(username, password string) (string, string, error) {
	if m.LoginFn != nil {
		return m.LoginFn(username, password)
	}
	return "", "", nil
}

func (m *MockAuthService) RefreshToken(refreshToken string) (string, error) {
	if m.RefreshTokenFn != nil {
		return m.RefreshTokenFn(refreshToken)
	}
	return "", nil
}
