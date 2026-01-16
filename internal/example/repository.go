package example

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type ExampleRepository struct {
	// fields
	Dbpool *pgxpool.Pool
}

func NewExampleRepository(dbpool *pgxpool.Pool) *ExampleRepository {
	return &ExampleRepository{
		Dbpool: dbpool,
	}
}

func (r *ExampleRepository) addExample(entity *ExampleEntity) error {
	// В этот момент поле entity.Id равно 0 (Zero Value), так как оно не пришло в JSON.
	// Go автоматически инициализирует пустые числовые поля нулями.

	// implementation
	// Мы намеренно НЕ указываем 'id' в части INSERT, чтобы PostgreSQL использовал свой
	// встроенный счетчик (SERIAL) и сам сгенерировал уникальный номер.
	query := `
				INSERT INTO example (example_column)
				VALUES (@example_column)
				RETURNING id, example_column
			`

	args := pgx.NamedArgs{
		"example_column": entity.ExampleColumn,
		// Важно: мы НЕ передаем сюда "id": entity.Id.
		// Если бы мы передали 0, база попыталась бы создать запись с ID=0, что вызвало бы ошибку дубликата.
	}

	// RETURNING id возвращает нам сгенерированное число (например, 55).
	// .Scan(&entity.Id) берет это число и перезаписывает наш "ноль" в структуре.
	// Так как entity передан по указателю (*ExampleEntity), это изменение увидит и Handler.
	err := r.Dbpool.QueryRow(context.Background(), query, args).Scan(
		&entity.Id,
		&entity.ExampleColumn,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to add example")
		return err
	}
	return nil
}
