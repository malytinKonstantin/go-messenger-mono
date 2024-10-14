   # Генерация gRPC кода

   Этот документ описывает шаги для генерации gRPC кода в проекте.

   ## Установка инструментов

   Сначала установите необходимые инструменты:

   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

   ## Настройка PATH

   Добавьте путь к бинарным файлам Go в переменную PATH:

   ```bash
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

   ## Подготовка скрипта генерации

   Сделайте скрипт генерации исполняемым:

   ```bash
   chmod +x generate.sh
   ```

   ## Запуск генерации

   Запустите скрипт для генерации кода:

   ```bash
   ./generate.sh
   ```

   После выполнения этих шагов, ваш gRPC код будет сгенерирован и готов к использованию.