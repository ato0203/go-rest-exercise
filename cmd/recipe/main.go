package main

import (
	"bukeuw/recipe/api/recipe"
	"bukeuw/recipe/pkg/recipes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupHandler() *recipe.RecipeHandler {
  store := recipes.NewMemStore()
  handler := recipe.NewRecipeHandler(store)
  return handler
}

func setupRouter(handler *recipe.RecipeHandler) *gin.Engine {
  r := gin.Default()
  api := r.Group("/api")
  {
    api.GET("health", func(ctx *gin.Context) {
      ctx.JSON(http.StatusOK, gin.H{
        "status": "OK",
      })
    })
    api.GET("recipes", handler.ListRecipe)
    api.POST("recipes", handler.CreateRecipe)
    api.GET("recipes/:id", handler.GetRecipe)
    api.PATCH("recipes/:id", handler.UpdateRecipe)
    api.DELETE("recipes/:id", handler.DeleteRecipe)
  }
  return r
}

func main() {
  handler := setupHandler()
  r := setupRouter(handler)
  r.Run()
}
