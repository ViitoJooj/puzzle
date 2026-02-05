package generator

var modelTpl = `
package auth

type User struct {
	ID uint ` + "`gorm:\"primaryKey\"`" + `

{{range .Fields}}
	{{title .}} string ` + "`json:\"{{.}}\"`" + `
{{end}}

	Password string
	RefreshToken string
}
`

var dtoTpl = `
package auth

type LoginDTO struct {
	Email string
	Password string
}

type RegisterDTO struct {
{{range .Fields}}
	{{title .}} string
{{end}}
	Password string
}

type RefreshDTO struct {
	Refresh string
}
`

var routesTpl = `
package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	h := NewHandler(db)

	g := r.Group("/auth")

	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.POST("/refresh", h.Refresh)
}
`

var handlersTpl = `
package auth

import "github.com/gin-gonic/gin"

type Handler struct {
	service *Service
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		service: NewService(db),
	}
}

func (h *Handler) Register(c *gin.Context) {

	var dto RegisterDTO
	c.ShouldBindJSON(&dto)

	user, err := h.service.Register(dto)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (h *Handler) Login(c *gin.Context) {

	var dto LoginDTO
	c.ShouldBindJSON(&dto)

	res, err := h.service.Login(dto)

	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)
}

func (h *Handler) Refresh(c *gin.Context) {

	var dto RefreshDTO
	c.ShouldBindJSON(&dto)

	token := h.service.Refresh(dto.Refresh)

	c.JSON(200, gin.H{"access": token})
}
`

var serviceTpl = `
package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) Register(dto RegisterDTO) (*User, error) {

	hash, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), 10)

	user := User{
		Email: dto.Email,
		Password: string(hash),
	}

	s.db.Create(&user)

	return &user, nil
}

func (s *Service) Login(dto LoginDTO) (any, error) {

	user := User{}

	if err := s.db.Where("email = ?", dto.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	access := GenerateJWT(user.ID)
	refresh := GenerateRefreshJWT(user.ID)

	user.RefreshToken = refresh
	s.db.Save(&user)

	return map[string]string{
		"access": access,
		"refresh": refresh,
	}, nil
}

func (s *Service) Refresh(token string) string {

	user := User{}
	s.db.Where("refresh_token = ?", token).First(&user)

	return GenerateJWT(user.ID)
}
`

var jwtTpl = `
package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(id uint) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": id,
			"exp": time.Now().Add(time.Minute * 15).Unix(),
		})

	signed, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return signed
}

func GenerateRefreshJWT(id uint) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": id,
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		})

	signed, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return signed
}
`

var middlewareTpl = `
package auth

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatus(401)
			return
		}

		c.Next()
	}
}
`
