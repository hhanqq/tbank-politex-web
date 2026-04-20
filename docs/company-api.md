# Выдача данных компаниям по API

## 1. Принцип

Компания не получает доступ к данным конкретного клиента. Она получает только обезличенные агрегированные дата-продукты через REST API или через асинхронную выгрузку.

Общий flow:

1. Компания проходит B2B-аутентификацию по OAuth2 client credentials.
2. Смотрит доступные дата-продукты.
3. Создает запрос с нужными фильтрами.
4. После одобрения получает данные либо по API, либо в виде файла.

## 2. Линейка B2B-продуктов

Ниже тот же набор данных, что пользователь разрешает передавать, но упакованный в понятные для компании продукты.

### Product 1. `audience-core`

Что входит:

- категории расходов;
- демографические срезы;
- агрегаты по периодам и регионам.

Когда нужен:

- сегментация аудитории;
- оценка спроса по категориям;
- построение маркетинговых гипотез.

Основные dimensions:

- `period`
- `region`
- `ageBand`
- `gender`
- `spendCategory`

Основные metrics:

- `users`
- `txCount`
- `avgCheck`
- `totalVolume`
- `categoryShare`

### Product 2. `behavior-geo`

Что входит:

- все из `audience-core`;
- поведенческие паттерны;
- агрегированная география;
- частота покупок и сезонность.

Когда нужен:

- поиск новых точек роста;
- оценка поведения сегментов;
- планирование локального маркетинга.

Основные dimensions:

- `period`
- `regionCluster`
- `segment`
- `spendCategory`
- `weekdayGroup`

Основные metrics:

- `users`
- `avgCheck`
- `purchaseFrequency`
- `seasonalityIndex`
- `onlineShare`

### Product 3. `financial-signals`

Что входит:

- все из `behavior-geo`;
- инвестиционные сегменты;
- агрегированные финансовые индикаторы;
- оперативные сигналы спроса.

Когда нужен:

- продвинутый market intelligence;
- оценка финансового поведения сегментов;
- быстрый мониторинг изменений спроса.

Основные dimensions:

- `period`
- `region`
- `segment`
- `incomeBand`
- `savingsBand`

Основные metrics:

- `users`
- `investmentActivityIndex`
- `savingsPropensityIndex`
- `financialStabilityIndex`
- `demandPulse`

Важно:

- не выдаем `user_id`, событие по клиенту или сырой транзакционный поток;
- `demandPulse` - это агрегированный сигнал по сегменту и периоду, а не realtime stream конкретных операций.

## 3. Endpoint'ы для компаний

### Каталог продуктов

`GET /v1/company/products`

Возвращает список доступных продуктов:

```json
{
  "items": [
    {
      "id": "audience-core",
      "name": "Audience Core",
      "description": "Категории расходов и демография в агрегированном виде",
      "accessLevel": "standard",
      "deliveryModes": ["api", "export"],
      "refreshRate": "daily"
    },
    {
      "id": "behavior-geo",
      "name": "Behavior Geo",
      "description": "Поведенческие сегменты и агрегированная география",
      "accessLevel": "advanced",
      "deliveryModes": ["api", "export"],
      "refreshRate": "daily"
    },
    {
      "id": "financial-signals",
      "name": "Financial Signals",
      "description": "Агрегированные финансовые индикаторы и сигналы спроса",
      "accessLevel": "premium",
      "deliveryModes": ["api", "export"],
      "refreshRate": "hourly"
    }
  ]
}
```

### Создание запроса

`POST /v1/company/requests`

Пример:

```json
{
  "productId": "behavior-geo",
  "filters": {
    "period": "2026-03",
    "region": "moscow",
    "segment": "travellers"
  },
  "deliveryMode": "api"
}
```

### Проверка статуса

`GET /v1/company/requests/{requestId}`

Ответ:

```json
{
  "requestId": "req_01J4YQ3",
  "productId": "behavior-geo",
  "status": "ready",
  "price": 0
}
```

## 4. Как именно отдаем данные

Есть 2 режима выдачи.

### Вариант A. Синхронный API

Подходит для небольших и средних агрегатов.

`GET /v1/company/datasets/{datasetId}?period=2026-03&region=moscow`

#### Пример для `audience-core`

```json
{
  "datasetId": "audience-core",
  "productId": "audience-core",
  "rows": [
    {
      "period": "2026-03",
      "region": "moscow",
      "ageBand": "25-34",
      "gender": "female",
      "spendCategory": "groceries",
      "users": 18420,
      "txCount": 53210,
      "avgCheck": 786.4,
      "totalVolume": 41852144.0,
      "categoryShare": 0.31
    }
  ],
  "privacy": {
    "kMin": 30,
    "noiseApplied": true,
    "suppressedRows": 2
  }
}
```

#### Пример для `behavior-geo`

```json
{
  "datasetId": "behavior-geo",
  "productId": "behavior-geo",
  "rows": [
    {
      "period": "2026-03",
      "regionCluster": "moscow_north",
      "segment": "travellers",
      "spendCategory": "travel",
      "weekdayGroup": "weekend",
      "users": 4210,
      "avgCheck": 4210.7,
      "purchaseFrequency": 3.8,
      "seasonalityIndex": 1.24,
      "onlineShare": 0.67
    }
  ],
  "privacy": {
    "kMin": 30,
    "noiseApplied": true,
    "suppressedRows": 1
  }
}
```

#### Пример для `financial-signals`

```json
{
  "datasetId": "financial-signals",
  "productId": "financial-signals",
  "rows": [
    {
      "period": "2026-03-20T10:00:00Z",
      "region": "moscow",
      "segment": "premium",
      "incomeBand": "high",
      "savingsBand": "medium",
      "users": 2180,
      "investmentActivityIndex": 0.74,
      "savingsPropensityIndex": 0.61,
      "financialStabilityIndex": 0.83,
      "demandPulse": 1.12
    }
  ],
  "privacy": {
    "kMin": 50,
    "noiseApplied": true,
    "suppressedRows": 3
  }
}
```

### Вариант B. Выгрузка

Подходит для больших наборов данных и BI-интеграций.

1. `POST /v1/company/datasets/{datasetId}/exports`
2. `GET /v1/company/exports/{exportId}`
3. Получить `downloadUrl`

Рекомендуемые форматы:

- `csv` для простых выгрузок;
- `parquet` для аналитических хранилищ;
- `json` для технических интеграций.

## 5. Ограничения и правила

- только whitelist компаний;
- только агрегированные срезы;
- минимальный размер группы проверяется privacy policy;
- запрещены фильтры, которые делают выборку слишком узкой;
- все запросы логируются в аудит;
- доступ может ограничиваться по продукту, отрасли и SLA.

## 6. Как это лучше продавать на защите

Правильная формулировка:

> Компания покупает не доступ к данным клиента, а доступ к обезличенному дата-продукту с заранее определенным набором агрегатов, частотой обновления и правилами приватности.
