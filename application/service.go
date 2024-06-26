package application

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

const (
	// FormTokenName is used to define HTML input field parameter name.
	FormTokenName = "csrftoken"
	// HeaderTokenName is used to define the header field name which can contain the token.
	HeaderTokenName = "Csrf-Token" //nolint:gosec // false positive
)

type (
	// Service is interface to define usage of service responsible for creating and validation csrf token.
	Service interface {
		Generate(session *web.Session) string
		IsValid(request *web.Request) bool
		IsValidPost(request *web.Request) bool
		IsValidHeader(request *web.Request) bool
	}

	// ServiceImpl is actual implementation of Service interface
	ServiceImpl struct {
		secret []byte
		ttl    int

		logger flamingo.Logger
	}

	csrfToken struct {
		ID   string    `json:"id"`
		Date time.Time `json:"date"`
	}
)

// Inject dependencies
func (s *ServiceImpl) Inject(l flamingo.Logger, cfg *struct {
	Secret string  `inject:"config:csrf.secret"`
	TTL    float64 `inject:"config:csrf.ttl"`
}) {
	hash := sha256.Sum256([]byte(cfg.Secret))
	s.secret = hash[:]
	s.ttl = int(cfg.TTL)
	s.logger = l
}

// Generate creates csrf token depending on user session ID and time.
// It uses AES standard for encrypting data.
func (s *ServiceImpl) Generate(session *web.Session) string {
	token := csrfToken{
		ID:   session.ID(),
		Date: time.Now(),
	}

	body, err := json.Marshal(token)
	if err != nil {
		s.logger.WithField("csrf", "jsonMarshal").Error(err.Error())
		return ""
	}

	gcm, err := s.getGcm()
	if err != nil {
		s.logger.WithField("csrf", "newGCM").Error(err.Error())
		return ""
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		s.logger.WithField("csrf", "nonceGenerate").Error(err.Error())
		return ""
	}

	cipherText := gcm.Seal(nil, nonce, body, nil)
	cipherText = append(nonce, cipherText...) //nolint:makezero // nonce size is a placeholder in the beginning

	return hex.EncodeToString(cipherText)
}

// IsValidPost validates csrf token from POST request.
// It uses AES standard for decrypting data.
// Session ID from csrf token must be the one in the request and token life time must be valid.
func (s *ServiceImpl) IsValidPost(request *web.Request) bool {
	if request.Request().Method != http.MethodPost {
		return true
	}

	formToken, err := request.Form1(FormTokenName)
	if err != nil {
		return false
	}

	return s.isValidToken(formToken, request)
}

// IsValidHeader validates csrf token in request header field.
// It uses AES standard for decrypting data.
// Session ID from csrf token must be the one in the request and token life time must be valid.
func (s *ServiceImpl) IsValidHeader(request *web.Request) bool {
	headerCsrfToken := request.Request().Header.Get(HeaderTokenName)
	if headerCsrfToken == "" {
		return false
	}

	return s.isValidToken(headerCsrfToken, request)
}

// IsValid validates csrf token from POST request. Deprecated - use IsVaildPost instead.
// It uses AES standard for decrypting data.
// Session ID from csrf token must be the one in the request and token life time must be valid.
func (s *ServiceImpl) IsValid(request *web.Request) bool {
	return s.IsValidPost(request)
}

func (s *ServiceImpl) isValidToken(inputToken string, request *web.Request) bool {
	data, err := hex.DecodeString(inputToken)
	if err != nil {
		return false
	}

	gcm, err := s.getGcm()
	if err != nil {
		return false
	}

	nonceSize := gcm.NonceSize()
	if len(data) <= nonceSize {
		return false
	}

	nonce := data[:nonceSize]
	cipherText := data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return false
	}

	var token csrfToken

	err = json.Unmarshal(plainText, &token)
	if err != nil {
		return false
	}

	if request.Session().ID() != token.ID {
		return false
	}

	if time.Now().Add(time.Duration(-s.ttl) * time.Second).After(token.Date) {
		return false
	}

	return true
}

func (s *ServiceImpl) getGcm() (cipher.AEAD, error) {
	block, err := aes.NewCipher(s.secret)
	if err != nil {
		return nil, fmt.Errorf("creaton of cipher block from configured secret csrf.secret failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creation of gcm failed: %w", err)
	}

	return gcm, nil
}
