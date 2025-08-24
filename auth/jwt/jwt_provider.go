package jwt

import (
	"fmt"
	"github.com/Lumicrate/gompose/auth"
	"github.com/Lumicrate/gompose/db"
	"github.com/Lumicrate/gompose/http"
	"github.com/Lumicrate/gompose/utils"
	"reflect"
	"time"
)

type JWTAuthProvider struct {
	SecretKey string
	UserModel any // optional: developer can override
	DB        db.DBAdapter
	TokenTTL  time.Duration
}

func NewJWTAuthProvider(secretKey string, dbAdapter db.DBAdapter) *JWTAuthProvider {
	return &JWTAuthProvider{
		SecretKey: secretKey,
		DB:        dbAdapter,
		UserModel: auth.UserModel{},
		TokenTTL:  time.Hour * 72, // default is 3 days
	}
}

func (j *JWTAuthProvider) Init() error {
	if j.SecretKey == "" {
		return fmt.Errorf("jwt: SecretKey must be provided")
	}

	if j.UserModel == nil {
		return fmt.Errorf("jwt: UserModel must be provided via SetUserModel")
	}

	if err := j.DB.Migrate([]any{j.UserModel}); err != nil {
		return fmt.Errorf("jwt: failed to migrate user model: %w", err)
	}

	return nil
}

func (j *JWTAuthProvider) RegisterRoutes(engine http.HTTPEngine) {
	engine.RegisterRoute("POST", "/auth/register", j.registerHandler, j.UserModel, false)
	engine.RegisterRoute("POST", "/auth/login", j.loginHandler, j.UserModel, false)
}

func (j *JWTAuthProvider) SetUserModel(model any) *JWTAuthProvider {
	if _, ok := model.(auth.AuthUser); !ok {
		panic("SetUserModel: model must implement AuthUser interface")
	}
	j.UserModel = model
	return j
}

func (j *JWTAuthProvider) SetTokenTTL(t time.Duration) *JWTAuthProvider {
	j.TokenTTL = t
	return j
}

func (j *JWTAuthProvider) registerHandler(ctx http.Context) {
	t := reflect.TypeOf(j.UserModel)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	newUser := reflect.New(t).Interface()

	if err := ctx.Bind(newUser); err != nil {
		ctx.JSON(400, map[string]string{"error": "invalid input: " + err.Error()})
		return
	}

	authUser, ok := newUser.(auth.AuthUser)
	if !ok {
		ctx.JSON(500, map[string]string{"error": "user model must implement AuthUser"})
		return
	}

	password := authUser.GetHashedPassword()
	hashed, err := utils.GenerateFromPassword(password)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "invalid input: " + err.Error()})
	}

	reflect.ValueOf(newUser).Elem().FieldByName("Password").SetString(hashed)
	idField := reflect.ValueOf(newUser).Elem().FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		switch idField.Kind() {
		case reflect.String:
			idField.SetString(utils.GenerateUUID())
		}
	}

	if err := j.DB.Create(newUser); err != nil {
		ctx.JSON(500, map[string]string{"error": "failed to create user: " + err.Error()})
		return
	}

	ctx.JSON(201, map[string]string{"message": "user registered successfully"})
}

func (j *JWTAuthProvider) loginHandler(ctx http.Context) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(400, map[string]string{"error": "invalid input: " + err.Error()})
		return
	}

	t := reflect.TypeOf(j.UserModel)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	user := reflect.New(t).Interface()

	foundUsers, err := j.DB.FindAll(user, map[string]any{
		"email": payload.Email,
	}, db.Pagination{Limit: 1}, nil)

	if err != nil {
		ctx.JSON(500, map[string]string{"error": "failed to query user"})
		return
	}

	usersVal := reflect.ValueOf(foundUsers)
	if usersVal.Len() == 0 {
		ctx.JSON(401, map[string]string{"error": "invalid username or password"})
		return
	}

	userVal := usersVal.Index(0)
	var authUser auth.AuthUser
	var ok bool

	if userVal.Kind() == reflect.Ptr {
		authUser, ok = userVal.Interface().(auth.AuthUser)
	} else {
		authUser, ok = userVal.Addr().Interface().(auth.AuthUser)
	}

	if !ok {
		ctx.JSON(500, map[string]string{"error": "user model must implement AuthUser"})
		return
	}

	if err := utils.CompareHashAndPassword(authUser.GetHashedPassword(), payload.Password); err != nil {
		ctx.JSON(401, map[string]string{"error": "invalid username or password"})
		return
	}

	token, err := utils.GenerateJWT(authUser.GetID(), j.SecretKey, j.TokenTTL)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "failed to generate token: " + err.Error()})
	}

	ctx.JSON(200, map[string]string{"token": token})
}

func (j *JWTAuthProvider) Middleware() http.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) {
			tokenStr, err := utils.ExtractBearerToken(ctx.Header("Authorization"))
			if err != nil {
				ctx.JSON(401, map[string]string{"error": err.Error()})
				ctx.Abort()
				return
			}

			claims, err := utils.ValidateJWT(tokenStr, j.SecretKey)
			if err != nil {
				ctx.JSON(401, map[string]string{"error": err.Error()})
				ctx.Abort()
				return
			}

			ctx.Set("user_id", claims["sub"])
			next(ctx)
		}
	}
}
