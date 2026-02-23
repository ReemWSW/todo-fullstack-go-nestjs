import { Injectable, InternalServerErrorException, NotFoundException } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { catchError, firstValueFrom, map } from 'rxjs';
import { CreateTodoDto } from './dto/create-todo.dto';
import { UpdateTodoDto } from './dto/update-todo.dto';
import { Todo } from './interfaces/todo.interface';

@Injectable()
export class TodosService {
  private readonly backendUrl = 'http://localhost:8080/api/todos';

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
