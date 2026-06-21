package components

import (
	"strings"
)

// UIBlock — интерфейс для любого блока страницы
type UIBlock interface {
	HTML() string
}

// Page — главный контейнер для сборки страницы
type Page struct {
	Title      string
	Components []UIBlock
}

// Render собирает финальный HTML-документ
func (p *Page) Render() string {
	var htmlBuilder strings.Builder

	// Рендерим компоненты
	for _, comp := range p.Components {
		htmlBuilder.WriteString(comp.HTML())
	}

	// Получаем глобальные стили и JS
	globalCSS := GetGlobalCSS()
	globalJS := GetGlobalJS()

	// Сборка итогового HTML
	return `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
  <title>` + p.Title + `</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
  <style>` + globalCSS + `</style>
</head>
<body>
  <div id="app">` + htmlBuilder.String() + `</div>
  <script>` + globalJS + `</script>
</body>
</html>`
}
