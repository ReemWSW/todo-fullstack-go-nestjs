import TodoList from '@/components/TodoList';

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-zinc-50 to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">
      <div className="max-w-2xl mx-auto px-4 py-12">
        <header className="text-center mb-8">
          <h1 className="text-4xl font-bold text-zinc-900 dark:text-zinc-100 mb-2">
            Todo List
          </h1>
          <p className="text-zinc-500 dark:text-zinc-400">
            Organize your tasks efficiently
          </p>
        </header>
        <TodoList />
      </div>
    </div>
  );
}
