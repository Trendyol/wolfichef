package controllers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"trendyol.com/security/appsec/devsecops/wolfichef/api"
	"trendyol.com/security/appsec/devsecops/wolfichef/api/models"
	"trendyol.com/security/appsec/devsecops/wolfichef/api/utils"
)

var verifiers map[string]Verifier

type Verifier struct {
	Verifier string
	Hashed   string
}

func init() {
	verifiers = make(map[string]Verifier)
}

func GitlabOAuthUrl(c *fiber.Ctx) error {
	identifier := fmt.Sprintf("%v", c.Locals("identifier"))

	_, ok := verifiers[identifier]
	if !ok {
		rnd, err := utils.GenerateRandomString(100)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("")
		}
		hash := sha256.Sum256([]byte(rnd))
		hashed := b64.RawURLEncoding.EncodeToString(hash[:])
		verifiers[identifier] = Verifier{
			Verifier: rnd,
			Hashed:   hashed,
		}
	}
	rnd, err := utils.GenerateRandomString(30)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("")
	}

	return c.JSON(struct {
		Url   string `json:"url"`
		State string `json:"state"`
	}{
		Url: fmt.Sprintf(
			"https://%s/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s&response_type=code&scope=api&code_challenge_method=S256",
			api.Http.Gitlab.Domain,
			api.Http.Gitlab.AppId,
			api.Http.Gitlab.RedirectUri,
			rnd,
			verifiers[identifier].Hashed),
		State: rnd,
	})
}

func GitlabFetchToken(c *fiber.Ctx) error {
	var body models.Token
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	_, errors := govalidator.ValidateStruct(body)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	reqBody := url.Values{}

	verifier := verifiers[fmt.Sprintf("%v", c.Locals("identifier"))]
	reqBody.Set("client_id", api.Http.Gitlab.AppId)
	reqBody.Set("client_secret", api.Http.Gitlab.SecretKey)
	reqBody.Set("code", body.Code)
	reqBody.Set("grant_type", "authorization_code")
	reqBody.Set("redirect_uri", api.Http.Gitlab.RedirectUri)
	reqBody.Set("code_verifier", verifier.Verifier)

	resp, err := http.Post(fmt.Sprintf("https://%s/oauth/token", api.Http.Gitlab.Domain), "application/x-www-form-urlencoded", strings.NewReader(reqBody.Encode()))
	if err != nil {
		return c.JSON(err)
	}
	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return c.JSON(err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func GitlabRefreshToken(c *fiber.Ctx) error {
	var body models.Token
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	_, errors := govalidator.ValidateStruct(body)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	reqBody := url.Values{}

	reqBody.Set("client_id", api.Http.Gitlab.AppId)
	reqBody.Set("client_secret", api.Http.Gitlab.SecretKey)
	reqBody.Set("refresh_token", body.Code)
	reqBody.Set("grant_type", "refresh_token")
	reqBody.Set("redirect_uri", api.Http.Gitlab.RedirectUri)

	resp, err := http.Post(fmt.Sprintf("https://%s/oauth/token", api.Http.Gitlab.Domain), "application/x-www-form-urlencoded", strings.NewReader(reqBody.Encode()))
	if err != nil {
		return c.JSON(err)
	}
	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return c.JSON(err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
