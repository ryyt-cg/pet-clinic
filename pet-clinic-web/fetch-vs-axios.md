# Tanstack Query Mutation - await fetch vs axios
For a TanStack Query mutation, you can use either the native fetch API or the Axios library, as TanStack Query is HTTP-client agnostic. The core difference lies not in the mutation itself, but in the specific features and syntax of the underlying data-fetching library.

The choice depends on your project's needs. fetch is a lightweight, built-in option, while Axios provides a more feature-rich developer experience, particularly for complex applications.

## Core difference: Error handling
The main distinction when integrating with TanStack Query is how each library handles HTTP errors (e.g., status codes like 404 or 500).

* Axios: Rejects the Promise by default for non-successful status codes, which is exactly what TanStack Query needs to register a mutation as an error.

* fetch: Does not reject the Promise on HTTP errors. The Promise is only rejected on network failures. You must manually check response.ok and throw an error yourself to trigger TanStack Query's error state.

## Comparison: await fetch vs. await axios
Feature 	fetch (Native)	Axios (Library)
Setup	No installation required, built into modern browsers.	Requires installation (npm install axios).
Error handling	Manual: You must check response.ok and throw an error yourself.	Automatic: Rejects the Promise automatically for HTTP error codes (e.g., 4xx, 5xx).
JSON parsing	Manual: Requires an extra step to call await response.json().	Automatic: The response data is automatically parsed and available at response.data.
Interceptors	Not built-in: Requires manual, boilerplate-heavy implementation.	Built-in: Allows you to globally handle requests (e.g., add auth tokens) or responses (e.g., handle errors).
Configuration	Manual: All options (headers, body) are configured on each call.	Centralized: Create an instance with a base URL and default headers, reducing repetition.
Request body	Manual: You must manually call JSON.stringify() for a JSON payload.	Automatic: Handles JSON stringification for you.
Bundle size	Smaller: As a native API, it adds no overhead to your bundle.	Slightly larger: Adds a small amount of overhead, though often worth the additional features.


## Example TanStack Query mutations
Using await fetch with useMutation
```typescript
import { useMutation, useQueryClient } from '@tanstack/react-query';

interface NewTodo {
  title: string;
}

interface Todo {
  id: number;
  title: string;
}

const addTodo = async (newTodo: NewTodo): Promise<Todo> => {
  const response = await fetch('/api/todos', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(newTodo),
  });

  if (!response.ok) {
    throw new Error('Failed to add todo');
  }

  return response.json();
};

function useAddTodoMutation() {
  const queryClient = useQueryClient();
  return useMutation<Todo, Error, NewTodo>({
    mutationFn: addTodo,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['todos'] });
    },
  });
}
```

Using await axios with useMutation

```typescript
import { useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';

interface NewTodo {
  title: string;
}

interface Todo {
  id: number;
  title: string;
}

const addTodoAxios = async (newTodo: NewTodo): Promise<Todo> => {
  const response = await axios.post<Todo>('/api/todos', newTodo);
  return response.data;
};

function useAddTodoMutationAxios() {
  const queryClient = useQueryClient();
  return useMutation<Todo, Error, NewTodo>({
    mutationFn: addTodoAxios,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['todos'] });
    },
  });
}
```

## Recommendation
For most applications, Axios offers a better developer experience with TanStack Query.
* Less boilerplate code due to automatic JSON handling and simplified error management.
* Powerful features like interceptors for centralized authentication or error logging are invaluable for real-world applications.

Use fetch if you are building a lightweight application and want to minimize dependencies, or if you prefer having granular, manual control over every aspect of your network requests.

