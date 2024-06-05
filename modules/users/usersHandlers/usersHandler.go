package usersHandlers

import (
	"fmt"

	"github.com/MarkTBSS/067_Refresh_Token/config"
	"github.com/MarkTBSS/067_Refresh_Token/modules/entities"
	"github.com/MarkTBSS/067_Refresh_Token/modules/users"
	"github.com/MarkTBSS/067_Refresh_Token/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	err := c.BodyParser(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			"users-001",
			err.Error(),
		).Res()
	}
	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			"users-001",
			"email pattern is invalid",
		).Res()
	}
	// Insert
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				"users-001",
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				"users-001",
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				"users-001",
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	err := c.BodyParser(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			"users-002",
			err.Error(),
		).Res()
	}
	fmt.Println("usersHandler.SignIn() : "+req.Email, req.Password)
	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			"users-002",
			err.Error(),
		).Res()
	}
	//passportJSON, _ := json.MarshalIndent(passport, "", "  ")
	//fmt.Printf("usersHandler.SignIn() : %s\n", passportJSON)

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	err := c.BodyParser(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			"users-003",
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			"users-003",
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}
