import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { TodosModule } from './todos/todos.module';
import { HttpModule } from '@nestjs/axios';

@Module({
  imports: [TodosModule, HttpModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
