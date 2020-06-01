package main

import (
	fmt "fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/tidwall/gjson"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
)

type PlacesController struct{}

func (*PlacesController) ListPlaces(c echo.Context) error {
	request := echo.Map{}
	response := echo.Map{}
	if err := c.Bind(&request); err != nil {
		return c.NoContent(500)
	}
	fmt.Println(request)
	query := `SELECT "Places".id,"Places".name, "Images".src as landing_image
	FROM "Places" LEFT JOIN "Images" ON "Images".id = "Places".landing_image 
	WHERE "Places".name  LIKE ?
	GROUP BY "Places".id, "Images".id`
	type Result struct {
		LandingImage string `json:"landing_image"`
		Name         string `json:"name"`
		Id           uint   `json:"id"`
	}
	var places []Result
	db.Raw(query, request["keyword"].(string)+"%").Scan(&places)
	response["places"] = places

	return c.JSON(200, response)
}

func (*PlacesController) GetPlaceByID(c echo.Context) error {
	request := echo.Map{}
	response := echo.Map{}
	if err := c.Bind(&request); err != nil {
		return c.NoContent(200)
	}
	query := `
	SELECT "Places".id,
	"Places".description,
	(SELECT COALESCE(json_agg(item),'[]'::json) FROM (SELECT * FROM "Images" WHERE "Images".id = ANY("Places".images)) item ) as "images",
	"Places".name,	
	(SELECT json_agg(item) FROM (SELECT * FROM "Images" WHERE "Images".id = "Places".landing_image) item ) as "background_image",
	(SELECT avg("Post".rating) FROM "Post" Where "Post".id = ANY("Places".post_ids)) as rating
	FROM "Places",
	"Images"
	WHERE "Places".id = ?
	GROUP BY "Places".id
	`

	type Result struct {
		ID              string         `json:"id"`
		Description     string         `json:"description"`
		Rating          float64        `json:"rating"`
		Images          postgres.Jsonb `json:"images"`
		Name            string         `json:"name"`
		BackgroundImage postgres.Jsonb `json:"background_image,omitempty"`
		Posts           []struct {
			ID        uint           `json:"id"`
			Content   string         `json:"content"`
			Avatar    string         `json:"avatar"`
			FullName  string         `json:"full_name"`
			AccountID string         `json:"account_id"`
			UserLiked postgres.Jsonb `json:"user_liked"`
			Comments  postgres.Jsonb `json:"comments"`
			Images    postgres.Jsonb `json:"images"`
		} `json:"posts"`
	}
	var result Result
	db.Raw(query, request["id"]).Scan(&result)
	query = `
	SELECT "Post".*,"Profile".avatar,"Profile".full_name ,
	(
		SELECT COALESCE(json_agg(DISTINCT item),'[]'::json)
		FROM (
			SELECT "Images".* FROM "Images" WHERE "Images".id = ANY("Post".image_ids)
		) item
	) as images,
	(
		SELECT COALESCE(json_agg(DISTINCT item),'[]'::json)
		FROM ( 
			SELECT "Profile".* FROM "Profile" WHERE "Profile".account_id = ANY("Post".account_liked)) item
	) as user_liked ,
	(
		SELECT COALESCE(json_agg(DISTINCT item),'[]'::json)
		FROM (SELECT "Comment".*,"Profile".* as user FROM "Comment","Profile"
		WHERE "Comment".account_id = "Profile".account_id AND "Comment".post_id = "Post".id		
		) as item
	) as comments
	FROM "Post","Profile","Places"
	WHERE 
	"Places".id = ?
	AND "Post".account_id = "Profile".account_id
	AND "Post".id = Any("Places".post_ids)
	GROUP BY "Post".id,"Profile".full_name,"Profile".avatar
	ORDER BY "Post".id DESC
	`
	db.Raw(query, request["id"]).Scan(&result.Posts)
	response["content"] = result
	return c.JSON(200, response)
}

func (*PlacesController) AddPlace(c echo.Context) error {
	fmt.Println(c.FormValue("address"))
	fmt.Println(c.FormValue("name"))
	fmt.Println(c.FormValue("description"))
	mf, _ := c.MultipartForm()
	imagesID := []int64{}
	if len(mf.File["file"]) > -1 {
		for _, fh := range mf.File["file"] {
			file, _ := fh.Open()
			fileData, _ := ioutil.ReadAll(file)
			fileName, _ := uuid.NewUUID()
			pathString := fileName.String() + filepath.Ext(fh.Filename)
			err := ioutil.WriteFile("./assets/"+pathString, fileData, 0644)
			if err != nil {
				panic(err)
			}
			images := Images{
				Size: uint(fh.Size),
				Type: filepath.Ext(fh.Filename)[1:],
				Src:  localServer + "/public/" + pathString,
			}
			if err := db.Create(&images).Error; err != nil {
				panic(err)
			}
			fmt.Println(images.ID)
			imagesID = append(imagesID, int64(images.ID))
		}
	}
	address := url.PathEscape(c.FormValue("address"))
	url := "https://google-maps-geocoding.p.rapidapi.com/geocode/json?language=en&address=" + address
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-host", "google-maps-geocoding.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "f93e17dc2emsh5d6001a1203bfe7p1c9d21jsnd5e65b818b81")
	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		return c.NoContent(500)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	lat := gjson.Get(string(body), "results.0.geometry.location.lat")
	lng := gjson.Get(string(body), "results.0.geometry.location.lng")

	place := Places{
		Coordinate:   `{"lat":` + lat.String() + `,"lng":` + lng.String() + `}`,
		AccountLiked: []int64{},
		Images:       imagesID,
		Address:      c.FormValue("address"),
		LandingImage: uint(imagesID[0]),
		Description:  c.FormValue("description"),
		Name:         c.FormValue("name"),
	}
	if err := db.Create(&place).Error; err != nil {
		panic(err)
	}
	return c.NoContent(200)
}
func (*PlacesController) UserReview(c echo.Context) error {
	// request := echo.Map{}
	fh, _ := c.FormFile("file")
	f, _ := fh.Open()
	fileData, _ := ioutil.ReadAll(f)
	fileName, _ := uuid.NewUUID()
	err := ioutil.WriteFile("./assets/"+fileName.String()+filepath.Ext(fh.Filename), fileData, 0644)
	if err != nil {
		return c.NoContent(500)
	}
	image := Images{Size: uint(fh.Size), Src: localServer + "/public/" + fileName.String() + filepath.Ext(fh.Filename)}
	err = db.Create(&image).Error
	if err != nil {
		panic(err)
	}
	accountID, _ := DecodeToken(c, "id")
	rating, _ := strconv.ParseFloat(c.FormValue("rating"), 64)

	post := Post{
		Content:      c.FormValue("content"),
		AccountID:    uint(int(accountID.(float64))),
		Like:         0,
		AccountLiked: []int64{},
		ImageIds:     []int64{int64(image.ID)},
		Rating:       rating,
	}
	err = db.Create(&post).Error
	if err != nil {
		return c.NoContent(500)
	}
	db.Table("Places").Where("id = ? ", c.FormValue("post_id")).Update("post_ids", gorm.Expr("array_append(post_ids,?)", post.ID)).Update("images", gorm.Expr("array_append(images,?)", image.ID))
	return c.NoContent(200)
}
