package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	name2ID = map[string]string{
		"Tim":    "935f871a-660f-4f19-801e-916c04bb0324",
		"Alex":   "a89b7b78-b9c1-4129-8cff-380bf53f3a49",
		"Arthur": "a98cd0f5-d6b2-4899-a1fb-ddf308d6f5c8",
		"Ray":    "a679ac51-08e8-45c7-80d7-019bf9dad64b",
		"HD":     "55b36756-6089-4756-bbd2-b0f66e50ee07",
		"peko":   "5a1e760e-76ea-4709-98ba-e1a701a4d340",
		"miko":   "201bef83-cc46-4acb-9c25-2eef60a59a9a",
		"rushia": "1c3e7209-fb42-4643-bfa6-c6a3fb42bf92",
		"gura":   "084e135f-78c7-406e-a347-94e38fa55b60",
		"Ame":    "8a180d2b-0965-4095-ba17-a880d196f04d",
	}
)

const ()

// simulate auth middleware to get user's information
func GetUserAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		accountID, ok := name2ID[token]
		if !ok {
			logrus.WithField("token", token).Warn("token not found")
			c.JSON(http.StatusUnauthorized, map[string]string{
				"errMessage": "token not found",
			})
			c.Abort()
		}

		c.Set("accountID", accountID)
		c.Next()
	}
}

func SetHandleContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ctx", context.Background())
		c.Next()
	}
}
