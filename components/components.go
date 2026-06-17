package components

import (
	"strings"
)

// UIBlock представляет рендерируемый блок пользовательского интерфейса
type UIBlock interface {
	HTML() string
}

// Page представляет главный контейнер страницы
type Page struct {
	Title      string
	Components []UIBlock
}

func (p *Page) Render() string {
	var htmlBuilder strings.Builder

	// Рендерим все дочерние компоненты
	for _, comp := range p.Components {
		htmlBuilder.WriteString(comp.HTML())
	}

	// Получаем глобальные стили и скрипты
	globalCSS := GetGlobalCSS()
	globalJS := GetGlobalJS()

	// Сборка итогового HTML-кода
	return `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
  <title>` + p.Title + `</title>
  <!-- Google Fonts -->
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Outfit:wght@500;600;700;800&display=swap" rel="stylesheet">
  <!-- MediaPipe Pose -->
  <script src="https://cdn.jsdelivr.net/npm/@mediapipe/camera_utils/camera_utils.js" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/@mediapipe/control_utils/control_utils.js" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/@mediapipe/drawing_utils/drawing_utils.js" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/@mediapipe/pose/pose.js" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
  <style>
` + globalCSS + `
  </style>
</head>
<body>
  <div class="app-container">
    ` + htmlBuilder.String() + `
  </div>
  <script type="module">
` + globalJS + `
  </script>
</body>
</html>`
}
