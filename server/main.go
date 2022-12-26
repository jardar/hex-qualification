package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"postapi/entity"
	"postapi/usecase"
	"strconv"
	"strings"
	"time"
)

const (
	JWTSECRET  = "QBZEypZVl0zfQrEkn4bOgA=="
	SSLPrivKey = "/etc/letsencrypt/live/%s/privkey.pem"
	SSLCertKey = "/etc/letsencrypt/live/%s/fullchain.pem"
)

var (
	listenPort = flag.String("port", "3000", "listen port")
	domain     = flag.String("d", "dash.playgaming.cc", "domain name for ssl")
	enableSSL  = flag.Bool("ssl", false, "enable TLS")
)

type MyClaims struct {
	UserEmail string `json:"email"`
	jwt.RegisteredClaims
}
type JWTOutput struct {
	UserEmail string  `json:"email"`
	Token     string  `json:"token"`
	Expires   float64 `json:"expires"` //second
}

func newJWTOutput(email string, expired time.Duration) (*JWTOutput, error) {
	expirationTime := time.Now().Add(expired)
	claims := &MyClaims{
		UserEmail: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	tokenString, err := token.SignedString([]byte(JWTSECRET))
	if err != nil {
		return nil, err
	}

	return &JWTOutput{
		UserEmail: email,
		Token:     tokenString,
		Expires:   expired.Seconds(),
	}, nil
}

func main() {
	flag.Parse()

	if len(*domain) == 0 {
		log.Fatal("must specify domain")
	}
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET,POST,PUT,DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			//return origin == "http://127.0.0.1:10600"
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/version", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"ok":      true,
			"msg":     "六角Vue直播班審核 API",
			"version": "1.0",
		})
	})

	//- 用戶註冊
	router.POST("/signup", func(c *gin.Context) {
		var bind entity.Member
		if err := c.ShouldBindJSON(&bind); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		if !bind.Valid() {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": "not valid",
			})
			return
		}

		if err := usecase.SignUp(bind.Email, bind.Pass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":  true,
			"msg": "sign up ok",
		})
	})

	//- 用戶登入
	router.POST("/login", func(c *gin.Context) {
		var bind entity.Member
		if err := c.ShouldBindJSON(&bind); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		if !bind.Valid() {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": "not valid",
			})
			return
		}

		if err := usecase.Login(bind.Email, bind.Pass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		dataOut, err := newJWTOutput(bind.Email, 10*time.Minute)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"msg":  "login ok",
			"data": dataOut,
		})
	})
	//- 景點列表
	router.GET("/posts", func(c *gin.Context) {
		posts, err := usecase.ReadPosts()
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": "error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"msg":  "success",
			"data": posts,
		})

	})
	//- 單一景點詳細資訊
	router.GET("/posts/:id", func(c *gin.Context) {
		postId := c.Param("id")
		if len(postId) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": "invalid (0)",
			})
			return
		}
		wantPostId, err := strconv.Atoi(postId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": "invalid",
			})
			return
		}
		post, err := usecase.ReadPostById(uint32(wantPostId))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"msg":  "success",
			"data": post,
		})
	})
	//// 需登入
	authorizedJWT := router.Group("/v1")
	authorizedJWT.Use(AuthMiddlewareJWT())

	//- 用戶登出
	// client 端把 jwt token 刪掉就可
	authorizedJWT.GET("/logout", func(c *gin.Context) {
		email, exist := c.Get("email")
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": "no user",
			})
		}
		got := email.(string)
		c.JSON(http.StatusOK, gin.H{
			"ok":    true,
			"msg":   "logout ok",
			"email": got,
		})
	})

	//- 單一景點詳細資訊 (已登入，帶入收藏狀態)
	authorizedJWT.GET("/posts/:id", func(c *gin.Context) {
		postId := c.Param("id")
		if len(postId) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": "invalid (0)",
			})
			return
		}
		wantPostId, err := strconv.Atoi(postId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": "invalid",
			})
			return
		}
		post, err := usecase.ReadPostById(uint32(wantPostId))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		// if login
		email, exist := c.Get("email")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{
				"ok":  false,
				"msg": "401 Unauthorized",
			})
			return
		}

		isBookmark := usecase.GetBookmarkState(email.(string), post.Id)

		c.JSON(http.StatusOK, gin.H{
			"ok":  true,
			"msg": "success",
			"data": &usecase.PostBookmark{
				Post:     post,
				Bookmark: isBookmark,
			},
		})
	})

	//- 單一景點收藏/取消收藏 （需登入）
	authorizedJWT.GET("/bookmarks/:id", func(c *gin.Context) {
		postId := c.Param("id")
		if len(postId) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": "invalid (0)",
			})
			return
		}
		wantPostId, err := strconv.Atoi(postId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": "invalid",
			})
			return
		}

		post, err := usecase.ReadPostById(uint32(wantPostId))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		// if login
		email, exist := c.Get("email")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{
				"ok":  false,
				"msg": "401 Unauthorized",
			})
			return
		}

		isBookmark := usecase.ToggleBookmark(email.(string), post.Id)

		c.JSON(http.StatusOK, gin.H{
			"ok":  true,
			"msg": "success",
			"data": &usecase.PostBookmark{
				Post:     post,
				Bookmark: isBookmark,
			},
		})
	})
	//- 收藏的景點列表（需登入）
	authorizedJWT.GET("/bookmarks", func(c *gin.Context) {

		// if login
		email, exist := c.Get("email")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{
				"ok":  false,
				"msg": "401 Unauthorized",
			})
			return
		}

		posts, err := usecase.ReadPostsBookmarked(email.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"msg":  "success",
			"data": posts,
		})
	})

	//- 新增景點（需登入）
	authorizedJWT.POST("/posts", func(c *gin.Context) {
		var bind entity.Post
		if err := c.ShouldBindJSON(&bind); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}
		log.Println("/posts bind=", bind)
		post, err := usecase.CreatePost(&bind)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"msg":  "success",
			"data": post,
		})
	})
	//- 編輯景點（需登入）
	authorizedJWT.PUT("/posts", func(c *gin.Context) {
		var bind entity.Post
		if err := c.ShouldBindJSON(&bind); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		post, err := usecase.UpdatePost(&bind)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"msg":  "success",
			"data": post,
		})
	})
	//- 刪除景點（需登入）
	authorizedJWT.DELETE("/posts/:id", func(c *gin.Context) {
		postId := c.Param("id")
		if len(postId) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": "invalid (0)",
			})
			return
		}
		wantPostId, err := strconv.Atoi(postId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ok":  false,
				"msg": "invalid",
			})
			return
		}

		post, err := usecase.DeletePost(uint32(wantPostId))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":  false,
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"msg":  "success",
			"data": post,
		})
	})

	//////////
	if *enableSSL {
		privkey := fmt.Sprintf(SSLPrivKey, *domain)
		certkey := fmt.Sprintf(SSLCertKey, *domain)

		err := router.RunTLS(":"+*listenPort, certkey, privkey)
		//err := router.Run(":" + *listenPort) // listen and serve on 0.0.0.0:3000
		if err != nil {
			log.Fatalln("Route can not be run", err)
		}
	} else {
		err := router.Run(":" + *listenPort) // listen and serve on 0.0.0.0:3000
		if err != nil {
			log.Fatalln("Route can not be run", err)
		}
	}
}

func AuthMiddlewareJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":  false,
				"msg": "401 Unauthorized",
			})
			return
		}

		if strings.HasPrefix(token, "Bearer") {
			token = token[7:]
		}
		log.Println("token=", token)
		myClaims := &MyClaims{}

		tkn, err := jwt.ParseWithClaims(token, myClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSECRET), nil
		})
		if err != nil {
			log.Println("JWT err=", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if claims, ok := tkn.Claims.(*MyClaims); ok && tkn.Valid {
			fmt.Printf("%v %v", claims.UserEmail, claims.RegisteredClaims.Issuer)
			c.Set("email", claims.UserEmail)

			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		}

		//// 在閞發時沒認證
		//c.Set("user", "JarDar")
		//c.Set("roles", []string{"admin", "user", "op"})
		//c.Next()

	}
}
