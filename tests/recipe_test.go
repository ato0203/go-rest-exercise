package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"

	"bukeuw/recipe/api/recipe"
	"bukeuw/recipe/pkg/recipes"
)

func TestListRecipe(t *testing.T) {
  router := gin.Default()
  api := router.Group("/api")
  store := recipes.NewMemStore()
  handler := recipe.NewRecipeHandler(store)
  api.GET("recipes", handler.ListRecipe)

  w := httptest.NewRecorder()

  ingredients := [3]recipes.Ingredient{
    { Name: "Tea"},
    { Name: "Water"},
    { Name: "Ice Cube"},
  }
  newRecipe := recipes.Recipe{
    Name: "Iced Tea",
    Ingredients: ingredients[:],
  }

  id := slug.Make(newRecipe.Name)
  store.Add(id, newRecipe)
  recipeJson, _ := json.Marshal(newRecipe)

  req, _ := http.NewRequest("GET", "/api/recipes", nil)
  router.ServeHTTP(w, req)

  assert.Equal(t, http.StatusOK, w.Code)
  assert.Equal(t, fmt.Sprintf("{\"recipes\":{\"%s\":%s}}", id, recipeJson), w.Body.String())
}

func TestCreateRecipe(t *testing.T) {
  router := gin.Default()
  api := router.Group("/api")
  store := recipes.NewMemStore()
  handler := recipe.NewRecipeHandler(store)
  api.POST("recipes", handler.CreateRecipe)

  w := httptest.NewRecorder()

  ingredients := [3]recipes.Ingredient{
    { Name: "Tea"},
    { Name: "Water"},
    { Name: "Ice Cube"},
  }
  newRecipe := recipes.Recipe{
    Name: "Iced Tea",
    Ingredients: ingredients[:],
  }
  recipeJson, _ := json.Marshal(newRecipe)
  req, _ := http.NewRequest("POST", "/api/recipes", strings.NewReader(string(recipeJson)))
  router.ServeHTTP(w, req)

  assert.Equal(t, http.StatusCreated, w.Code)
  assert.Equal(t, fmt.Sprintf("{\"recipe\":%s}", recipeJson), w.Body.String())
}

func TestGetRecipe(t *testing.T) {
  router := gin.Default()
  api := router.Group("/api")
  store := recipes.NewMemStore()
  handler := recipe.NewRecipeHandler(store)
  api.GET("recipes/:id", handler.GetRecipe)

  w := httptest.NewRecorder()

  ingredients := [3]recipes.Ingredient{
    { Name: "Tea"},
    { Name: "Water"},
    { Name: "Ice Cube"},
  }
  newRecipe := recipes.Recipe{
    Name: "Iced Tea",
    Ingredients: ingredients[:],
  }

  id := slug.Make(newRecipe.Name)
  store.Add(id, newRecipe)
  recipeJson, _ := json.Marshal(newRecipe)

  req, _ := http.NewRequest("GET", "/api/recipes/" + id, nil)
  router.ServeHTTP(w, req)

  assert.Equal(t, http.StatusOK, w.Code)
  assert.Equal(t, fmt.Sprintf("{\"recipe\":%s}", recipeJson), w.Body.String())
}

func TestUpdateRecipe(t *testing.T) {
  router := gin.Default()
  api := router.Group("/api")
  store := recipes.NewMemStore()
  handler := recipe.NewRecipeHandler(store)
  api.PATCH("recipes/:id", handler.UpdateRecipe)

  w := httptest.NewRecorder()

  ingredients := [2]recipes.Ingredient{
    { Name: "Tea"},
    { Name: "Water"},
  }
  teaRecipe := recipes.Recipe{
    Name: "Tea",
    Ingredients: ingredients[:],
  }

  id := slug.Make(teaRecipe.Name)
  store.Add(id, teaRecipe)

  newIngredients := [3]recipes.Ingredient{
    { Name: "Tea"},
    { Name: "Water"},
    { Name: "Ice Cube"},
  }
  icedTeaRecipe := recipes.Recipe{
    Name: "Iced Tea",
    Ingredients: newIngredients[:],
  }
  recipeJson, _ := json.Marshal(icedTeaRecipe)

  req, _ := http.NewRequest("PATCH", "/api/recipes/" + id, strings.NewReader(string(recipeJson)))
  router.ServeHTTP(w, req)

  assert.Equal(t, http.StatusOK, w.Code)
  assert.Equal(t, fmt.Sprintf("{\"recipe\":%s}", recipeJson), w.Body.String())

  updatedRecipe, _ := store.Get(id)
  assert.Equal(t, updatedRecipe.Name, icedTeaRecipe.Name)
  assert.Equal(t, updatedRecipe.Ingredients, icedTeaRecipe.Ingredients)
}

func TestDeleteRecipe(t *testing.T) {
  router := gin.Default()
  api := router.Group("/api")
  store := recipes.NewMemStore()
  handler := recipe.NewRecipeHandler(store)
  api.DELETE("recipes/:id", handler.DeleteRecipe)

  w := httptest.NewRecorder()

  ingredients := [3]recipes.Ingredient{
    { Name: "Tea"},
    { Name: "Water"},
    { Name: "Ice Cube"},
  }
  newRecipe := recipes.Recipe{
    Name: "Iced Tea",
    Ingredients: ingredients[:],
  }

  id := slug.Make(newRecipe.Name)
  store.Add(id, newRecipe)

  req, _ := http.NewRequest("DELETE", "/api/recipes/" + id, nil)
  router.ServeHTTP(w, req)

  assert.Equal(t, http.StatusOK, w.Code)

  _, err := store.Get(id)
  assert.Equal(t, err, recipes.NotFoundErr)
}
