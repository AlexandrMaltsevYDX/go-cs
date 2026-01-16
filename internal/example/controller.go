package example

import "github.com/gofiber/fiber/v2"

type ExampleController struct {
	// fields

	router     fiber.Router
	Repository *ExampleRepository
}

func NewController(repository *ExampleRepository) *ExampleController {
	return &ExampleController{
		Repository: repository,
	}
}

func (ctrl *ExampleController) Create(c *fiber.Ctx) error {
	// implementation

	// new(ExampleEntity) выделяет память в куче и возвращает указатель (*ExampleEntity).
	// При этом Go гарантирует, что все поля будут иметь "Zero Values" (нулевые значения):
	// Id = 0, ExampleColumn = ""
	entity := new(ExampleEntity)

	// BodyParser пытается сопоставить ключи из JSON с полями структуры.
	// 1. Если поле есть в JSON ("example_column": "val"), оно записывается в entity.
	// 2. Если поля НЕТ в JSON ("id"), парсер его пропускает.
	//    В итоге entity.Id остается равным 0 (Zero Value).
	// Парсер НЕ ругается на отсутствие полей, если нет явной валидации.
	if err := c.BodyParser(entity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Мы передаем структуру с Id=0 в репозиторий.
	// Репозиторий должен проигнорировать этот 0 и позволить базе сгенерировать настоящий ID.
	if err := ctrl.Repository.addExample(entity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Здесь entity.Id уже обновлен репозиторием (например, стал 55),
	// поэтому мы возвращаем клиенту полный объект.
	return c.Status(fiber.StatusCreated).JSON(entity)
}
