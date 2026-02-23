# Fullstack Todo List Project

โปรเจกต์นี้คือแอปพลิเคชัน Todo List แบบ Fullstack ที่พัฒนาขึ้นเพื่อแสดงการทำงานร่วมกันระหว่าง Backend ที่สร้างด้วย Golang (Gin, GORM) และ Frontend ที่สร้างด้วย NestJS โดยมีการเชื่อมต่อผ่าน RESTful API และใช้ Docker สำหรับการจัดการสภาพแวดล้อม

## ภาพรวมระบบ

-   **Backend**: Golang, Gin Framework, GORM ORM, PostgreSQL Database
-   **Frontend**: NestJS
-   **การเชื่อมต่อ**: RESTful API
-   **การจัดการสภาพแวดล้อม**: Docker, Docker Compose

## 1. Backend Requirements (Golang + Gin + GORM)

### 1.1 Project Structure

โครงสร้างโฟลเดอร์ของ Backend ได้รับการออกแบบมาให้เป็นแบบ production-ready เพื่อความสามารถในการขยายและบำรุงรักษา:

```
backend/
├── cmd/
│   └── main.go             # Entry point of the application
├── config/
│   └── database.go         # Database connection and initialization
├── controllers/
│   └── todo_controller.go  # Handles HTTP requests and responses
├── middleware/
│   ├── cors.go             # CORS middleware
│   └── logging.go          # Request logging middleware
├── models/
│   └── todo.go             # Defines the Todo data model
├── repository/
│   └── todo_repository.go  # Handles database operations
├── routes/
│   └── todo_routes.go      # Defines API routes
├── services/
│   └── todo_service.go     # Contains business logic
└── utils/
    └── response.go         # Standardized JSON response format
```

### 1.2 Database

ใช้ PostgreSQL เป็นฐานข้อมูล โดยมีตาราง `todos` ดังนี้:

-   `id` (uuid, primary key)
-   `title` (string, required)
-   `description` (string)
-   `completed` (boolean, default false)
-   `created_at` (timestamp)
-   `updated_at` (timestamp)

**GORM Auto Migrate**: GORM จะทำการ migrate โครงสร้างตาราง `todos` โดยอัตโนมัติเมื่อแอปพลิเคชันเริ่มต้นทำงาน

### 1.3 API Endpoints (RESTful)

Backend มี API endpoints สำหรับจัดการ Todo items ดังนี้:

| Method   | Endpoint           | Description          | Request Body (JSON)                               | Response Body (JSON)                               |
| :------- | :----------------- | :------------------- | :------------------------------------------------ | :------------------------------------------------- |
| `GET`    | `/api/todos`       | ดึงรายการ Todo ทั้งหมด | -                                                 | `{"success": true, "data": [...], "message": ""}` |
| `GET`    | `/api/todos/:id`   | ดึง Todo ตาม ID       | -                                                 | `{"success": true, "data": {...}, "message": ""}`  |
| `POST`   | `/api/todos`       | สร้าง Todo ใหม่       | `{"title": "string", "description": "string"}` | `{"success": true, "data": {...}, "message": ""}`  |
| `PUT`    | `/api/todos/:id`   | อัปเดต Todo ตาม ID    | `{"title": "string", "description": "string", "completed": boolean}` | `{"success": true, "data": {...}, "message": ""}`  |
| `DELETE` | `/api/todos/:id`   | ลบ Todo ตาม ID        | -                                                 | `{"success": true, "data": null, "message": ""}`   |

### 1.4 Architecture

Backend ใช้ **Layered Architecture** (Controller -> Service -> Repository) เพื่อแยกความรับผิดชอบของแต่ละส่วนอย่างชัดเจน:

-   **Controller**: รับผิดชอบการจัดการ HTTP request และ response
-   **Service**: รับผิดชอบ business logic
-   **Repository**: รับผิดชอบการสื่อสารกับฐานข้อมูล

**คุณสมบัติเพิ่มเติม:**

-   **Environment Variables**: ใช้ไฟล์ `.env` สำหรับการตั้งค่าฐานข้อมูลและพอร์ต
-   **CORS Middleware**: เปิดใช้งาน CORS เพื่อให้ Frontend สามารถเรียก API ได้
-   **JSON Response Format**: ใช้รูปแบบ response มาตรฐาน `{"success": true, "data": {}, "message": ""}`
-   **Error Handling**: มีการจัดการข้อผิดพลาดแบบ global
-   **Validation Request**: ใช้ `github.com/go-playground/validator/v10` สำหรับการตรวจสอบข้อมูลที่ส่งเข้ามาใน request
-   **Logging Middleware**: บันทึกข้อมูล request และ response เพื่อการ debug

### 1.5 ตัวอย่างไฟล์สำคัญ (Backend)

#### `backend/cmd/main.go`

```go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"todo-fullstack/backend/config"
	"todo-fullstack/backend/controllers"
	"todo-fullstack/backend/middleware"
	"todo-fullstack/backend/models"
	"todo-fullstack/backend/repository"
	"todo-fullstack/backend/services"
	"todo-fullstack/backend/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize database
	db := config.InitDB()

	// Auto migrate models
	err := db.AutoMigrate(&models.Todo{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// Initialize repository, service, and controller
	todoRepository := repository.NewTodoRepository(db)
	todoService := services.NewTodoService(todoRepository)
	todoController := controllers.NewTodoController(todoService)

	// Setup Gin router
	r := gin.Default()

	// CORS Middleware
	r.Use(middleware.CORSMiddleware())

	// Logging Middleware
	r.Use(middleware.LoggingMiddleware())

	// Setup routes
	routes.SetupTodoRoutes(r, todoController)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
```

#### `backend/models/todo.go`

```go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Todo represents a todo item
type Todo struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Title       string    `gorm:"not null" json:"title" binding:"required"`
	Description string    `json:"description"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
}

// BeforeCreate will set a UUID for the Todo ID
func (todo *Todo) BeforeCreate(tx *gorm.DB) (err error) {
	todo.ID = uuid.New()
	return
}
```

#### `backend/controllers/todo_controller.go`

```go
package controllers

import (
	"net/http"

	"todo-fullstack/backend/models"
	"todo-fullstack/backend/services"
	"todo-fullstack/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TodoController interface {
	GetAllTodos(c *gin.Context)
	GetTodoByID(c *gin.Context)
	CreateTodo(c *gin.Context)
	UpdateTodo(c *gin.Context)
	DeleteTodo(c *gin.Context)
}

type todoController struct {
	todoService services.TodoService
	validate    *validator.Validate
}

func NewTodoController(service services.TodoService) TodoController {
	return &todoController{
		todoService: service,
		validate:    validator.New(),
	}
}

func (ctrl *todoController) GetAllTodos(c *gin.Context) {
	todos, err := ctrl.todoService.GetAllTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(todos, "Todos retrieved successfully"))
}

func (ctrl *todoController) GetTodoByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Todo ID"))
		return
	}

	todo, err := ctrl.todoService.GetTodoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Todo not found"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(todo, "Todo retrieved successfully"))
}

func (ctrl *todoController) CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	// Validate the struct
	if err := ctrl.validate.Struct(todo); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	createdTodo, err := ctrl.todoService.CreateTodo(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse(createdTodo, "Todo created successfully"))
}

func (ctrl *todoController) UpdateTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Todo ID"))
		return
	}

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	// Validate the struct
	if err := ctrl.validate.Struct(todo); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	updatedTodo, err := ctrl.todoService.UpdateTodo(id, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(updatedTodo, "Todo updated successfully"))
}

func (ctrl *todoController) DeleteTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Todo ID"))
		return
	}

	if err := ctrl.todoService.DeleteTodo(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Todo deleted successfully"))
}
```

### 1.6 วิธีรัน Backend

1.  **สร้างไฟล์ `.env`**: ในโฟลเดอร์ `backend` สร้างไฟล์ชื่อ `.env` และเพิ่มการตั้งค่าดังนี้:

    ```dotenv
    PORT=8080
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=todo_db
    DB_PORT=5432
    DB_SSLMODE=disable
    DB_TIMEZONE=Asia/Bangkok
    ```

    *หมายเหตุ: หากรันด้วย Docker Compose, `DB_HOST` จะเป็น `db` (ชื่อ service ของ PostgreSQL)*

2.  **ติดตั้ง Dependencies**: หากยังไม่ได้ติดตั้ง ให้รันคำสั่งในโฟลเดอร์ `backend`:

    ```bash
    go mod tidy
    ```

3.  **รันแอปพลิเคชัน**: ในโฟลเดอร์ `backend` รันคำสั่ง:

    ```bash
    go run cmd/main.go
    ```

    Backend จะทำงานบนพอร์ต 8080 (หรือตามที่ระบุใน `PORT` ใน `.env`)

### 1.7 ตัวอย่าง Curl Test API (Backend)

#### Create Todo (POST)

```bash
curl -X POST \\
  http://localhost:8080/api/todos \\
  -H 'Content-Type: application/json' \\
  -d '{"title": "Buy groceries", "description": "Milk, eggs, bread"}'
```

#### Get All Todos (GET)

```bash
curl -X GET \\
  http://localhost:8080/api/todos
```

#### Get Todo by ID (GET)

```bash
# แทนที่ <TODO_ID> ด้วย ID ของ Todo ที่ต้องการ
curl -X GET \\
  http://localhost:8080/api/todos/<TODO_ID>
```

#### Update Todo (PUT)

```bash
# แทนที่ <TODO_ID> ด้วย ID ของ Todo ที่ต้องการ
curl -X PUT \\
  http://localhost:8080/api/todos/<TODO_ID> \\
  -H 'Content-Type: application/json' \\
  -d '{"title": "Buy groceries", "description": "Milk, eggs, bread, cheese", "completed": true}'
```

#### Delete Todo (DELETE)

```bash
# แทนที่ <TODO_ID> ด้วย ID ของ Todo ที่ต้องการ
curl -X DELETE \\
  http://localhost:8080/api/todos/<TODO_ID>
```

## 2. Frontend Requirements (NestJS)

Frontend สร้างด้วย NestJS เพื่อทำหน้าที่เป็น API consumer และ expose endpoint ของตัวเอง

### 2.1 Module Structure

-   `todos` module: โมดูลหลักสำหรับจัดการ Todo items
    -   `todos.service.ts`: Service สำหรับเรียก Backend API
    -   `todos.controller.ts`: Controller สำหรับ expose frontend endpoint
    -   `dto/`: โฟลเดอร์สำหรับ Data Transfer Objects (DTOs) เช่น `create-todo.dto.ts`, `update-todo.dto.ts`
    -   `interfaces/`: โฟลเดอร์สำหรับ TypeScript interfaces เช่น `todo.interface.ts`

### 2.2 ใช้ HTTP Module เรียก Backend

Frontend ใช้ `@nestjs/axios` (HTTP Module) ใน `TodosService` เพื่อเชื่อมต่อไปยัง Backend API ที่ `http://localhost:8080/api/todos` (หรือ `http://backend:8080/api/todos` เมื่อรันด้วย Docker Compose)

### 2.3 Endpoint ฝั่ง Frontend

Frontend จะมี endpoints ที่ mirror ของ Backend เพื่อให้ client อื่นๆ สามารถเรียกใช้งานได้:

| Method   | Endpoint     | Description          |
| :------- | :----------- | :------------------- |
| `GET`    | `/todos`     | ดึงรายการ Todo ทั้งหมด |
| `GET`    | `/todos/:id` | ดึง Todo ตาม ID       |
| `POST`   | `/todos`     | สร้าง Todo ใหม่       |
| `PUT`    | `/todos/:id` | อัปเดต Todo ตาม ID    |
| `DELETE` | `/todos/:id` | ลบ Todo ตาม ID        |

### 2.4 คุณสมบัติเพิ่มเติม (Frontend)

-   **DTO Validation**: ใช้ `class-validator` และ `class-transformer` สำหรับการตรวจสอบข้อมูลใน DTOs
-   **Environment Config**: ใช้ `.env` สำหรับการตั้งค่า URL ของ Backend
-   **Error Handling**: มีการจัดการข้อผิดพลาดเมื่อเรียก Backend API
-   **Response Mapping**: มีการ map response จาก Backend ให้เป็นรูปแบบที่เหมาะสมสำหรับ Frontend

### 2.5 ตัวอย่างไฟล์สำคัญ (Frontend)

#### `frontend/src/todos/interfaces/todo.interface.ts`

```typescript
export interface Todo {
  id: string;
  title: string;
  description?: string;
  completed: boolean;
  createdAt: string;
  updatedAt: string;
}
```

#### `frontend/src/todos/dto/create-todo.dto.ts`

```typescript
import { IsBoolean, IsNotEmpty, IsOptional, IsString } from 'class-validator';

export class CreateTodoDto {
  @IsNotEmpty()
  @IsString()
  title: string;

  @IsOptional()
  @IsString()
  description?: string;

  @IsOptional()
  @IsBoolean()
  completed?: boolean;
}
```

#### `frontend/src/todos/todos.service.ts`

```typescript
import { Injectable, InternalServerErrorException, NotFoundException } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { catchError, firstValueFrom, map } from 'rxjs';
import { CreateTodoDto } from './dto/create-todo.dto';
import { UpdateTodoDto } from './dto/update-todo.dto';
import { Todo } from './interfaces/todo.interface';

@Injectable()
export class TodosService {
  private readonly backendUrl = process.env.BACKEND_URL || 'http://localhost:8080/api/todos';

  constructor(private readonly httpService: HttpService) {}

  async findAll(): Promise<Todo[]> {
    try {
      const { data } = await firstValueFrom(
        this.httpService.get<any>(this.backendUrl).pipe(
          catchError((error) => {
            console.error('Error fetching todos:', error.response?.data || error.message);
            throw new InternalServerErrorException('Failed to fetch todos from backend');
          }),
          map((response) => response.data),
        ),
      );
      return data;
    } catch (error) {
      throw error;
    }
  }

  async findOne(id: string): Promise<Todo> {
    try {
      const { data } = await firstValueFrom(
        this.httpService.get<any>(`${this.backendUrl}/${id}`).pipe(
          catchError((error) => {
            console.error(`Error fetching todo with ID ${id}:`, error.response?.data || error.message);
            if (error.response && error.response.status === 404) {
              throw new NotFoundException(`Todo with ID ${id} not found`);
            }
            throw new InternalServerErrorException('Failed to fetch todo from backend');
          }),
          map((response) => response.data),
        ),
      );
      return data;
    } catch (error) {
      throw error;
    }
  }

  async create(createTodoDto: CreateTodoDto): Promise<Todo> {
    try {
      const { data } = await firstValueFrom(
        this.httpService.post<any>(this.backendUrl, createTodoDto).pipe(
          catchError((error) => {
            console.error('Error creating todo:', error.response?.data || error.message);
            throw new InternalServerErrorException('Failed to create todo in backend');
          }),
          map((response) => response.data),
        ),
      );
      return data;
    } catch (error) {
      throw error;
    }
  }

  async update(id: string, updateTodoDto: UpdateTodoDto): Promise<Todo> {
    try {
      const { data } = await firstValueFrom(
        this.httpService.put<any>(`${this.backendUrl}/${id}`, updateTodoDto).pipe(
          catchError((error) => {
            console.error(`Error updating todo with ID ${id}:`, error.response?.data || error.message);
            if (error.response && error.response.status === 404) {
              throw new NotFoundException(`Todo with ID ${id} not found`);
            }
            throw new InternalServerErrorException('Failed to update todo in backend');
          }),
          map((response) => response.data),
        ),
      );
      return data;
    } catch (error) {
      throw error;
    }
  }

  async remove(id: string): Promise<void> {
    try {
      await firstValueFrom(
        this.httpService.delete<any>(`${this.backendUrl}/${id}`).pipe(
          catchError((error) => {
            console.error(`Error deleting todo with ID ${id}:`, error.response?.data || error.message);
            if (error.response && error.response.status === 404) {
              throw new NotFoundException(`Todo with ID ${id} not found`);
            }
            throw new InternalServerErrorException('Failed to delete todo from backend');
          }),
        ),
      );
    } catch (error) {
      throw error;
    }
  }
}
```

#### `frontend/src/todos/todos.controller.ts`

```typescript
import { Controller, Get, Post, Body, Put, Param, Delete, HttpCode, HttpStatus, UsePipes, ValidationPipe } from '@nestjs/common';
import { TodosService } from './todos.service';
import { CreateTodoDto } from './dto/create-todo.dto';
import { UpdateTodoDto } from './dto/update-todo.dto';
import { Todo } from './interfaces/todo.interface';

@Controller('todos')
export class TodosController {
  constructor(private readonly todosService: TodosService) {}

  @Get()
  async findAll(): Promise<Todo[]> {
    return this.todosService.findAll();
  }

  @Get(':id')
  async findOne(@Param('id') id: string): Promise<Todo> {
    return this.todosService.findOne(id);
  }

  @Post()
  @HttpCode(HttpStatus.CREATED)
  @UsePipes(new ValidationPipe({ transform: true }))
  async create(@Body() createTodoDto: CreateTodoDto): Promise<Todo> {
    return this.todosService.create(createTodoDto);
  }

  @Put(':id')
  @UsePipes(new ValidationPipe({ transform: true }))
  async update(@Param('id') id: string, @Body() updateTodoDto: UpdateTodoDto): Promise<Todo> {
    return this.todosService.update(id, updateTodoDto);
  }

  @Delete(':id')
  @HttpCode(HttpStatus.NO_CONTENT)
  async remove(@Param('id') id: string): Promise<void> {
    return this.todosService.remove(id);
  }
}
```

### 2.6 วิธีรัน Frontend

1.  **สร้างไฟล์ `.env`**: ในโฟลเดอร์ `frontend` สร้างไฟล์ชื่อ `.env` และเพิ่มการตั้งค่าดังนี้:

    ```dotenv
    PORT=3000
    BACKEND_URL=http://localhost:8080/api/todos
    ```

    *หมายเหตุ: หากรันด้วย Docker Compose, `BACKEND_URL` จะเป็น `http://backend:8080/api/todos` (ชื่อ service ของ Backend)*

2.  **ติดตั้ง Dependencies**: ในโฟลเดอร์ `frontend` รันคำสั่ง:

    ```bash
    npm install
    ```

3.  **รันแอปพลิเคชัน**: ในโฟลเดอร์ `frontend` รันคำสั่ง:

    ```bash
    npm run start:dev
    ```

    Frontend จะทำงานบนพอร์ต 3000 (หรือตามที่ระบุใน `PORT` ใน `.env`)

## 3. Docker และ Docker Compose

โปรเจกต์นี้มาพร้อมกับ Dockerfile สำหรับทั้ง Backend และ Frontend รวมถึงไฟล์ `docker-compose.yml` เพื่อให้ง่ายต่อการตั้งค่าและรันทั้งระบบ

### 3.1 `docker-compose.yml`

```yaml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    container_name: todo_db
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: todo_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  backend:
    build: ./backend
    container_name: todo_backend
    restart: always
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=todo_db
      - DB_PORT=5432
      - DB_SSLMODE=disable
      - DB_TIMEZONE=Asia/Bangkok
    depends_on:
      - db

  frontend:
    build: ./frontend
    container_name: todo_frontend
    restart: always
    ports:
      - "3000:3000"
    environment:
      - BACKEND_URL=http://backend:8080/api/todos
    depends_on:
      - backend

volumes:
  postgres_data:
```

### 3.2 `backend/Dockerfile`

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

CMD ["./main"]
```

### 3.3 `frontend/Dockerfile`

```dockerfile
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM node:20-alpine

WORKDIR /app

COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
COPY .env .

EXPOSE 3000

CMD ["npm", "run", "start:prod"]
```

### 3.4 วิธีรันด้วย Docker Compose

1.  **สร้างไฟล์ `.env`**: ตรวจสอบให้แน่ใจว่ามีไฟล์ `.env` ในโฟลเดอร์ `backend` และ `frontend` ตามตัวอย่างด้านบน (โดยเฉพาะ `DB_HOST=db` และ `BACKEND_URL=http://backend:8080/api/todos`)

2.  **รัน Docker Compose**: ในโฟลเดอร์หลักของโปรเจกต์ (`todo-fullstack/`) รันคำสั่ง:

    ```bash
    docker-compose up --build
    ```

    คำสั่งนี้จะสร้างและรัน Docker containers สำหรับ PostgreSQL, Backend และ Frontend

    -   Backend จะสามารถเข้าถึงได้ที่ `http://localhost:8080`
    -   Frontend จะสามารถเข้าถึงได้ที่ `http://localhost:3000`

3.  **หยุด Docker Compose**: หากต้องการหยุดการทำงาน ให้กด `Ctrl+C` ใน Terminal แล้วรัน:

    ```bash
    docker-compose down
    ```

## 4. SQL Schema Example

ไฟล์ `schema.sql` แสดงโครงสร้างของตาราง `todos` สำหรับ PostgreSQL:

```sql
-- SQL Schema for Todo List
-- Note: GORM will auto-migrate this, but this is for reference

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS todos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed);
```

## 5. สรุปและพร้อมใช้งาน

โปรเจกต์นี้ได้รับการออกแบบมาให้เป็นตัวอย่างที่ครบถ้วนและพร้อมใช้งานสำหรับการพัฒนา Fullstack ToDo List ด้วย Golang และ NestJS โดยเน้นที่ Best Practices ในการจัดโครงสร้างโค้ด, การจัดการ API, การเชื่อมต่อฐานข้อมูล และการใช้ Docker เพื่อความสะดวกในการ deploy และ scale

หากมีข้อสงสัยหรือต้องการปรับปรุงเพิ่มเติม สามารถแจ้งได้เลยครับ
